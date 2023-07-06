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
	// fields.
	IncludeFieldsCmd = `fields ` +
		`messageId, timestamp, type, userId, event_verb, event_subject, properties`

	// ExcludeFieldsCmd is a Splunk query "fields" command for removing fields
	// from query results that splunk includes by default and we don't need.
	ExcludeFieldsCmd = `fields - _*`
)

// GenDedupEval generates a Splunk "eval" command for converting multivalue
// fields into single-value fields. See RHTAPWATCH-293 for why we need this.
func GenDedupEval(fields []string) string {
	var builder strings.Builder
	builder.WriteString("eval ")
	for i, field := range fields {
		if i > 0 {
			builder.WriteString(",")
		}
		fmt.Fprintf(&builder, `%s=mvindex('%s', 0)`, field, field)
	}
	return builder.String()
}

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
	with_userid, with_ev_verb, with_ev_subject bool
}

// Generate Splunk query commands for the output fields needed to build a
// Segment track object. spec is used to toggle some optional fields that
// different objects may/may not include while properties_map specified which
// fields to include in the properties sub-object.
func GenTrackFields(spec TrackFieldSpec, properties_map map[string]string) string {
	var builder strings.Builder
	builder.WriteString("eval ")
	builder.WriteString(`messageId=auditID,`)
	builder.WriteString(`timestamp=requestReceivedTimestamp,`)
	builder.WriteString(`type="track",`)
	if spec.with_userid {
		builder.WriteString(`userId=` + UserIdExpr + `,`)
	}
	if spec.with_ev_verb {
		builder.WriteString(`event_verb=verb,`)
	}
	if spec.with_ev_subject {
		builder.WriteString(`event_subject='objectRef.resource',`)
	}
	builder.WriteString(`properties=` + GenPropertiesJSONExpr(properties_map))
	builder.WriteString(`|` + IncludeFieldsCmd)
	builder.WriteString(`|` + ExcludeFieldsCmd)
	return builder.String()
}

// GenApplicationQuery returns a Splunk query for generating Segment events
// representing AppStudio Application object events.
func GenApplicationQuery(index string) string {
	fields := []string{
		"auditID",
		"impersonatedUser.username",
		"user.username",
		"objectRef.resource",
		"objectRef.namespace",
		"objectRef.apiGroup",
		"objectRef.apiVersion",
		"objectRef.name",
		"verb",
		"requestReceivedTimestamp",
	}
	json_properties := map[string]string{
		"apiGroup":   "objectRef.apiGroup",
		"apiVersion": "objectRef.apiVersion",
		"kind":       "objectRef.resource",
		"name":       "objectRef.name",
	}
	return `search ` +
		`index="` + index + `" ` +
		`log_type=audit ` +
		`NOT verb IN (get, watch, list, deletecollection) ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="appstudio.redhat.com" ` +
		`"objectRef.resource"="applications" ` +
		`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) ` +
		`(verb!=create OR "responseObject.metadata.resourceVersion"="*")` +
		`|` + GenDedupEval(fields) +
		`|` + GenTrackFields(TrackFieldSpec{true, true, true}, json_properties)
}
