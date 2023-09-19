package querygen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFilter struct{}

func (f *TestFilter) FieldSet() FieldSet {
	return FieldSet{
		"tf1": {srcFields: []string{"tf1"}},
		"tf2": {srcFields: []string{"tf2"}},
	}
}

func (f *TestFilter) Commands() []string {
	return []string{`eval tf1="hello"`, `eval tf2="world"`}
}

func TestUserJourneyQuery(t *testing.T) {
	q, err := NewUserJourneyQuery("idx").
		WithPredicate("verb=created").
		WithEventExpr(`"Event Name"`).
		WithCommands(`eval foo="bar"`).
		WithFilter(&TestFilter{}).
		WithFields("tf1", "tf2", "userId").
		String()
	assert.Nil(t, err)
	assert.Equal(t,
		`search index="idx" log_type=audit verb=created `+
			`| eval foo="bar" `+
			`| eval tf1="hello" `+
			`| eval tf2="world"`+
			`|eval event="Event Name",`+
			`event_subject='objectRef.resource',`+
			`event_verb='verb',`+
			`messageId='auditID',`+
			`tf1='tf1',`+
			`tf2='tf2',`+
			`timestamp='requestReceivedTimestamp',`+
			`type="track",`+
			`userId=if(isnull('impersonatedUser.username'),`+
			`'user.username','impersonatedUser.username'),`+
			`properties=json_object(`+
			`"apiGroup",'objectRef.apiGroup',`+
			`"apiVersion",'objectRef.apiVersion',`+
			`"kind",'objectRef.resource'`+
			`),`+
			`context=json_object("userAgent",'userAgent')`+
			`|fields `+
			`context,event,event_subject,event_verb,messageId,namespace,`+
			`properties,tf1,tf2,timestamp,type,userId`+
			`|`+excludeFieldsCmd,
		q,
	)
}

func TestUserJourneyQueryUnknownField(t *testing.T) {
	_, err := NewUserJourneyQuery("idx").WithFields("not-found").String()
	assert.NotNil(t, err)
}
