package querygen

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// ExcludeFieldsCmd is a Splunk query "fields" command for removing fields
	// from query results that splunk includes by default and we don't need.
	excludeFieldsCmd = `fields - _*`
)

// FieldSet defines a set of fields that can be included in the query output,
// how to generated the values for them from the input index and how to place
// then in the output
//
// The keys in the FieldSet are the output field names while the values include
// details about how to obtain the output values and how to present them.
//
// The zero value for the value struct defines that the output field is copied
// as-is from the input.
type FieldSet map[string]struct {
	// srcFields defines which input fields to use. Each field in the list is
	// used as a fallback for when all the fields that precede it are NULL.
	srcFields []string
	// srcExpr is an expression for generating the field value, if provided, it
	// overrides the fallback logic for srcFields. The expression may be a
	// literal value, in which case srcFields may be empty.
	srcExpr string
	// subObj defines a JSON sub-object for the field to reside in when its
	// included in the output
	subObj string
}

// QueryGen generates a Splunk query with searchExpr where its output includes
// the given fields. The values for the fields and how to present them in the
// output are determined from the FieldSet
func (fs FieldSet) QueryGen(searchExpr string, fields []string) (string, error) {
	var queryElements, evalElements []string
	var err error
	if evalElements, err = fs.collectEvalElements(fields); err != nil {
		return "", err
	}
	queryElements = append(
		queryElements,
		searchExpr,
	)
	if len(evalElements) > 0 {
		queryElements = append(queryElements, "eval "+commaSep(evalElements))
	}
	queryElements = append(
		queryElements,
		"fields "+commaSep(fs.collectIncludeFields()),
		excludeFieldsCmd,
	)
	return strings.Join(queryElements, "|"), nil
}

func commaSep(words []string) string {
	return strings.Join(words, ",")
}

// collectEvalElements generates a list of Splunk `eval` expressions for
// generating the values for the given output fields.
func (fs FieldSet) collectEvalElements(fields []string) ([]string, error) {
	var evalElements, subObjects []string
	subObjectFields := map[string][]string{}

	for _, field := range fields {
		if spec, ok := fs[field]; ok {
			expr := spec.srcExpr
			if expr == "" {
				expr = mkFieldSrcEvalExpr(spec.srcFields)
			}
			if spec.subObj == "" {
				if expr != "" {
					evalElements = append(evalElements, field+"="+expr)
				}
				continue
			}
			if expr == "" {
				expr = sQuot(field)
			}
			sof, ok := subObjectFields[spec.subObj]
			if !ok {
				subObjects = append(subObjects, spec.subObj)
			}
			subObjectFields[spec.subObj] = append(sof, dQuot(field), expr)
		} else {
			return []string{}, fmt.Errorf(`no field specification for: "%s"`, field)
		}
	}
	for _, subObject := range subObjects {
		evalElements = append(
			evalElements,
			fmt.Sprintf(
				"%s=json_object(%s)",
				subObject,
				commaSep(subObjectFields[subObject]),
			),
		)
	}
	return evalElements, nil
}

// collectIncludeFields generates a list of all the top-level fields that may be
// included in the query output
func (fs FieldSet) collectIncludeFields() (fields []string) {
	seenSubObj := map[string]bool{}
	for field, spec := range fs {
		if spec.subObj == "" {
			fields = append(fields, field)
		} else if !seenSubObj[spec.subObj] {
			seenSubObj[spec.subObj] = true
			fields = append(fields, spec.subObj)
		}
	}
	sort.Strings(fields)
	return
}

// mkFieldSrcEvalExpr generates a Splunk `eval` expression for getting the value
// from the given srcFields so that each field in the list is a fallback for the
// ones before it.
//
// Example: given srcFields = {"a", "b", "c"}
//
//	if(isnull('a'),if(isnull('b'),'c','b'),'a')
func mkFieldSrcEvalExpr(srcFields []string) string {
	if len(srcFields) <= 0 {
		return ""
	}
	expr := sQuot(srcFields[len(srcFields)-1])
	for i := len(srcFields) - 2; i >= 0; i-- {
		expr = fmt.Sprintf(
			"if(isnull('%s'),%s,'%s')",
			srcFields[i], expr, srcFields[i],
		)
	}
	return expr
}

func dQuot(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}

func sQuot(v string) string {
	return fmt.Sprintf(`'%s'`, v)
}
