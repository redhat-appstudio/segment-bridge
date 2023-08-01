// Package querygen is used to generate Splunk queries for fetching user journey
// events from the RHTAP K8s event log.
package querygen

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// UserIdExpr is a Splunk eval expression for obtaining the human username
	// behind a given event even if the event action was done on their behalf by
	// a service account
	UserIdExpr = `if(isnull('impersonatedUser.username'),` +
		`'user.username','impersonatedUser.username')`

	// IncludeFieldsCmd is a Splunk query "fields" command for specifying which
	// fields to include in the final query results. Since the structure of a
	// Segment record is quite fixed, most queries should return the same
	// fields. Some queries may return a subset of these fields by not
	// calculating values for all of them.
	IncludeFieldsCmd = `fields ` +
		`event_subject,event_verb,messageId,namespace,properties,timestamp,type,userId`

	// ExcludeFieldsCmd is a Splunk query "fields" command for removing fields
	// from query results that splunk includes by default and we don't need.
	ExcludeFieldsCmd = `fields - _*`
)

// GenPropertiesJSONExpr generates a Splunk "eval" expression for generating a
// JSON object containing the fields and values specified in properties_map. The
// field names are sorted in the output to make it predictable.
func GenPropertiesJSONExpr(properties_map map[string]string) string {
	var prop_args []string = make([]string, 0, len(properties_map))
	for json_field, splunk_field := range properties_map {
		prop_args = append(
			prop_args,
			fmt.Sprintf(`"%s",'%s'`, json_field, splunk_field),
		)
	}
	sort.Strings(prop_args)
	return fmt.Sprintf("json_object(%s)", strings.Join(prop_args, ","))
}

// TrackFieldSpec specifies which optional fields of a Segment track object do
// we want to include in the query output. Queried objects should have the
// values needed to generate these fields for this to work.
type TrackFieldSpec struct {
	with_userid, with_ev_verb, with_ev_subject, with_namespace bool
}

// Generate Splunk query commands for the output fields needed to build a
// Segment track object. spec is used to toggle some optional fields that
// different objects may/may not include while properties_map specified which
// fields to include in the properties sub-object.
func GenTrackFields(spec TrackFieldSpec, properties_map map[string]string) string {
	var builder strings.Builder
	builder.WriteString("eval ")
	if spec.with_ev_subject {
		builder.WriteString(`event_subject='objectRef.resource',`)
	}
	if spec.with_ev_verb {
		builder.WriteString(`event_verb='verb',`)
	}
	builder.WriteString(`messageId='auditID',`)
	if spec.with_namespace {
		builder.WriteString(`namespace='objectRef.namespace',`)
	}
	builder.WriteString(`timestamp='requestReceivedTimestamp',`)
	builder.WriteString(`type="track",`)
	if spec.with_userid {
		builder.WriteString(`userId=` + UserIdExpr + `,`)
	}
	builder.WriteString(`properties=` + GenPropertiesJSONExpr(properties_map))
	builder.WriteString(`|` + IncludeFieldsCmd)
	builder.WriteString(`|` + ExcludeFieldsCmd)
	return builder.String()
}

// GenApplicationQuery returns a Splunk query for generating Segment events
// representing AppStudio Application object events.
func GenApplicationQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=create ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="appstudio.redhat.com" ` +
		`"objectRef.resource"="applications" ` +
		`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) ` +
		`(verb!=create OR "responseObject.metadata.resourceVersion"="*")`,
		[]string{"name", "userId"},
	)
	return q
}

// GenPipelineRunQuery returns a Splunk query for generating Segment events
// representing creation of AppStudio build PipelineRuns.
func GenPipelineRunQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=create ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="tekton.dev" ` +
		`"objectRef.resource"="pipelineruns" ` +
		`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build` +
		`"responseObject.metadata.resourceVersion"="*"`,
		[]string{"namespace","application","component"},
	)
	return q
}
