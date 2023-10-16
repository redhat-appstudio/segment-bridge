package querygen

import (
	"fmt"
	"sort"
	"strings"
)

// UserJourneyQuery is a builder for Splunk queries.
type UserJourneyQuery struct {
	// The Splunk index to be searched
	index string

	// The initial search predicate used to narrow down results
	predicate string

	// Additional Splunk commands to execute in order immediately after the search
	// command.
	commands []string

	// Fields to return from the query
	fields []string

	// A FieldSet containing all possible output fields for the query. If a Filter
	// is used, its FieldSet is merged with the default set.
	fieldset FieldSet
}

// NewUserJourneyQuery constructs a default UserJourneyQuery
func NewUserJourneyQuery(index string) *UserJourneyQuery {
	return &UserJourneyQuery{
		index: index,
		fields: []string{
			"apiGroup",
			"apiVersion",
			"event_subject",
			"event_verb",
			"kind",
			"messageId",
			"namespace",
			"timestamp",
			"type",
			"userAgent",
		},
		fieldset: FieldSet{
			"messageId":     {srcFields: []string{"auditID"}},
			"timestamp":     {srcFields: []string{"requestReceivedTimestamp"}},
			"type":          {srcExpr: `"track"`},
			"userAgent":     {subObj: "context"},
			"userId":        {srcFields: []string{"impersonatedUser.username", "user.username"}},
			"namespace":     {srcFields: []string{"objectRef.namespace"}},
			"event_verb":    {srcFields: []string{"verb"}},
			"event_subject": {srcFields: []string{"objectRef.resource"}},
			"apiGroup":      {subObj: "properties", srcFields: []string{"objectRef.apiGroup"}},
			"apiVersion":    {subObj: "properties", srcFields: []string{"objectRef.apiVersion"}},
			"kind":          {subObj: "properties", srcFields: []string{"objectRef.resource"}},
			"name":          {subObj: "properties", srcFields: []string{"objectRef.name"}},
			"src_url":       {subObj: "properties", srcFields: []string{"responseObject.spec.source.git.url"}},
			"src_revision":  {subObj: "properties", srcFields: []string{"responseObject.spec.source.git.revision"}},
			"src_context":   {subObj: "properties", srcFields: []string{"responseObject.spec.source.git.context"}},
			"application": {
				subObj: "properties",
				srcFields: []string{
					"responseObject.spec.application",
					"responseObject.metadata.labels.appstudio.openshift.io/application",
				},
			},
			"commit_sha": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.build.appstudio.redhat.com/commit_sha"},
			},
			"component": {
				subObj: "properties",
				srcFields: []string{
					"responseObject.spec.componentName",
					"responseObject.metadata.labels.appstudio.openshift.io/component",
				},
			},
			"repo": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.build.appstudio.openshift.io/repo"},
				srcExpr:   `replace('responseObject.metadata.annotations.build.appstudio.openshift.io/repo',"^([^?]*)(.*)?","\1")`,
			},
			"target_branch": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.build.appstudio.redhat.com/target_branch"},
			},
			"git_trigger_event_type": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.pipelinesascode.tekton.dev/event-type"},
			},
			"git_trigger_provider": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.pipelinesascode.tekton.dev/git-provider"},
			},
			"pipeline_log_url": {
				subObj:    "properties",
				srcFields: []string{"responseObject.metadata.annotations.pipelinesascode.tekton.dev/log-url"},
			},
			"vulnerabilities_critical": {
				subObj:    "properties",
				srcFields: []string{"clair_scan_result.vulnerabilities.critical"},
			},
			"vulnerabilities_high": {
				subObj:    "properties",
				srcFields: []string{"clair_scan_result.vulnerabilities.high"},
			},
			"vulnerabilities_medium": {
				subObj:    "properties",
				srcFields: []string{"clair_scan_result.vulnerabilities.medium"},
			},
			"vulnerabilities_low": {
				subObj:    "properties",
				srcFields: []string{"clair_scan_result.vulnerabilities.low"},
			},
			"merge_url": {
				subObj:    "properties",
				srcFields: []string{"build_status.pac.merge-url"},
			},
		},
	}
}

// WithPredicate adds additional expressions to the leading search command.
func (q *UserJourneyQuery) WithPredicate(predicate string) *UserJourneyQuery {
	q.predicate = predicate
	return q
}

// WithCommands adds raw Splunk commands to the query.
// Each call appends to the existing set of commands so order of invocation is important.
func (q *UserJourneyQuery) WithCommands(commands ...string) *UserJourneyQuery {
	q.commands = append(q.commands, commands...)
	return q
}

// WithFilter adds all the Splunk commands for a Filter to the query.
// Each call appends to the existing set of commands so order of invocation is important.
func (q *UserJourneyQuery) WithFilter(filter Filter) *UserJourneyQuery {
	for k, v := range filter.FieldSet() {
		q.fieldset[k] = v
	}
	return q.WithCommands(filter.Commands()...)
}

// WithFields adds fields to the output of the query.
func (q *UserJourneyQuery) WithFields(fields ...string) *UserJourneyQuery {
	q.fields = append(q.fields, fields...)
	return q
}

// WithEventExpr adds a Splunk 'eval' expression specifically for the 'event' output field.
// This is handy when trying to override the default event naming logic.
func (q *UserJourneyQuery) WithEventExpr(expr string) *UserJourneyQuery {
	q.fieldset["event"] = &FieldSetSpec{srcExpr: expr}
	q.fields = append(q.fields, "event")
	return q
}

// String builds the Splunk query.
func (q *UserJourneyQuery) String() (string, error) {
	sort.Strings(q.fields) // To make test results predictable

	searchCmd := strings.TrimSpace(fmt.Sprintf(
		`search index="%s" log_type=audit %s`, q.index, q.predicate,
	))
	commands := append([]string{searchCmd}, q.commands...)
	query := strings.Join(commands, " | ")

	return q.fieldset.QueryGen(query, q.fields)
}
