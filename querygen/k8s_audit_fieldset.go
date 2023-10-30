package querygen

import (
	"fmt"
	"strings"
)

// K8sApiId defines a K8s API by including details about the API group and
// resource names
type K8sApiId struct {
	apiGroup string
	resource string
}

// Maps K8s API identifiers to field sets. This map should at least contain
// a value for the K8sApiId zero value. That value is used as the FieldSet for
// querying. When querying for an API object for which a value exists in the map
// the fields in the value are added to the zero-value fields to make up the
// final FieldSet used for the query.
// This allows to have different field settings for different K8s APIs.
type K8sAuditFieldSet map[K8sApiId]FieldSet

func (kfs K8sAuditFieldSet) QueryGen(
	index string, api K8sApiId, searchExpr string, fields []string, extra ...FieldSet,
) (string, error) {
	searchCmd := strings.TrimSpace(fmt.Sprintf(
		`search index="%s" log_type=audit `+
			`"objectRef.apiGroup"="%s" `+
			`"objectRef.resource"="%s" `+
			`%s`,
		index,
		api.apiGroup,
		api.resource,
		searchExpr,
	))
	allFieldSets := []FieldSet{kfs[K8sApiId{}], kfs[api]}
	allFieldSets = append(allFieldSets, extra...)
	fieldSet := FieldSet{}
	for _, fieldSetToAdd := range allFieldSets {
		for fld, spec := range fieldSetToAdd {
			fieldSet[fld] = spec
		}
	}

	return fieldSet.QueryGen(searchCmd, fields)
}
