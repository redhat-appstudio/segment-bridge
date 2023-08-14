package querygen

import (
	"fmt"
	"sort"
)

var UserJourneyFieldSet = FieldSet{
	"messageId":     {srcFields: []string{"auditID"}},
	"timestamp":     {srcFields: []string{"requestReceivedTimestamp"}},
	"type":          {srcExpr: `"track"`},
	"userAgent":     {subObj: "context"},
	"userId":        {srcFields: []string{"impersonatedUser.username", "user.username"}},
	"namespace":     {srcFields: []string{"objectRef.namespace"}},
	"event":         {srcFields: []string{"event"}},
	"event_verb":    {srcFields: []string{"verb"}},
	"event_subject": {srcFields: []string{"objectRef.resource"}},
	"apiGroup":      {subObj: "properties", srcFields: []string{"objectRef.apiGroup"}},
	"apiVersion":    {subObj: "properties", srcFields: []string{"objectRef.apiVersion"}},
	"kind":          {subObj: "properties", srcFields: []string{"objectRef.resource"}},
	"name":          {subObj: "properties", srcFields: []string{"objectRef.name"}},
	"application": {
		subObj:    "properties",
		srcFields: []string{"responseObject.metadata.labels.appstudio.openshift.io/application"},
	},
	"component": {
		subObj:    "properties",
		srcFields: []string{"responseObject.metadata.labels.appstudio.openshift.io/component"},
	},
}

var UserJourneyCommonFields = [...]string{
	"apiGroup",
	"apiVersion",
	"event_subject",
	"event_verb",
	"kind",
	"messageId",
	"timestamp",
	"type",
	"userAgent",
}

func UserJourneyQueryGen(index, predicate string, fields []string) (string, error) {
	var realFields []string = UserJourneyCommonFields[:]
	realFields = append(realFields, fields...)
	sort.Strings(realFields) // To make test results predictable
	searchExpr := fmt.Sprintf(`search index="%s" log_type=audit %s`, index, predicate)
	return UserJourneyFieldSet.QueryGen(searchExpr, realFields)
}
