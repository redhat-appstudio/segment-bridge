package querygen

import (
	"fmt"
	"strings"
)

// The Filter interface must be implemented by each filter.
type Filter interface {
	// Commands provides a sequence of Splunk commands for narrowing down search results.
	Commands() []string
	// FieldSet returns a map of all possible fields that can be included in the output
	// from this filter.
	FieldSet() FieldSet
}

// Optional fields for customizing StatusConditionFilter.
// All options will be combined using an "AND" boolean expression.
type StatusConditionOpts struct {
	// If defined, only match conditions with one of the specified reasons.
	reasons []string
	// If defined, only match conditions with one of the specified statuses.
	statuses []string
	// If defined, only match conditions with this message.
	// Evaluated using the Splunk 'like' function to allow for the use of wildcards.
	message string
}

// StatusConditionFilter matches audit records for any k8s resource based on its
// status conditions.
type StatusConditionFilter struct {
	// The status condition's 'type' fields value
	cType string
	// The name of the field used to track the position of the desired status condition.
	indexField string
	// Optional params
	opts StatusConditionOpts
}

// NewStatusConditionFilter creates a default StatusConditionFilter
func NewStatusConditionFilter(cType string) *StatusConditionFilter {
	return &StatusConditionFilter{
		cType:      cType,
		indexField: "status_condition_index",
	}
}

func (f *StatusConditionFilter) FieldSet() FieldSet {
	return FieldSet{
		"status_message": {
			subObj:  "properties",
			srcExpr: fmt.Sprintf(`mvindex('responseObject.status.conditions{}.message', %s)`, f.indexField),
		},
		"status_reason": {
			subObj:  "properties",
			srcExpr: fmt.Sprintf(`mvindex('responseObject.status.conditions{}.reason', %s)`, f.indexField),
		},
	}
}

func (f *StatusConditionFilter) Commands() []string {
	evalCmd := fmt.Sprintf(
		`eval %s=mvfind('responseObject.status.conditions{}.type', "%s")`,
		f.indexField, f.cType,
	)
	whereCmd := fmt.Sprintf(`where isnotnull(%s)`, f.indexField)

	if len(f.opts.reasons) > 0 {
		whereCmd += fmt.Sprintf(
			` AND mvindex('responseObject.status.conditions{}.reason', %s) IN (%s)`,
			f.indexField,
			`"`+strings.Join(f.opts.reasons, `", "`)+`"`,
		)
	}

	if len(f.opts.statuses) > 0 {
		whereCmd += fmt.Sprintf(
			` AND mvindex('responseObject.status.conditions{}.status', %s) IN (%s)`,
			f.indexField,
			`"`+strings.Join(f.opts.statuses, `", "`)+`"`,
		)
	}

	if f.opts.message != "" {
		whereCmd += fmt.Sprintf(
			` AND like(mvindex('responseObject.status.conditions{}.message', %s), "%s")`,
			f.indexField, f.opts.message,
		)
	}

	return []string{evalCmd, whereCmd}
}

// TektonTaskResultFilter will match audit records for Tekton TaskRun resources
// based on conditions for its task results.
type TektonTaskResultFilter struct {
	// The name of the result key
	name string
	// The name of the field used to track the position of the matching task result.
	indexField string
}

// NewTektonTaskResultFilter creates a default TektonTaskResultFilter
func NewTektonTaskResultFilter(name string) *TektonTaskResultFilter {
	return &TektonTaskResultFilter{
		name:       name,
		indexField: "tekton_task_result_index",
	}
}

func (f *TektonTaskResultFilter) FieldSet() FieldSet {
	return FieldSet{
		"tekton_task_result": {
			subObj:  "properties",
			srcExpr: fmt.Sprintf(`mvindex('responseObject.status.taskResults{}.value', %s)`, f.indexField),
		},
	}
}

func (f *TektonTaskResultFilter) Commands() []string {
	return []string{
		fmt.Sprintf(
			`eval %s=mvfind('responseObject.status.taskResults{}.name', "%s")`,
			f.indexField, f.name,
		),
		fmt.Sprintf(`where isnotnull(%s)`, f.indexField),
	}
}
