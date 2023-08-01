package querygen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserJourneyQueryGen(t *testing.T) {
	queryTemplate := `search index="%s" log_type=audit %s|%s|` +
		`eval ` +
		`event_subject='objectRef.resource',` +
		`event_verb='verb',` +
		`messageId='auditID',` +
		`timestamp='requestReceivedTimestamp',` +
		`type="track",` +
		`%s` +
		`properties=json_object(` +
		`"apiGroup",'objectRef.apiGroup',` +
		`"apiVersion",'objectRef.apiVersion',` +
		`"kind",'objectRef.resource'` +
		`)` +
		`|fields ` +
		`event_subject,event_verb,messageId,namespace,` +
		`properties,timestamp,type,userId` +
		`|` + excludeFieldsCmd
	type args struct {
		index     string
		predicate string
		fields    []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "With invalid fields",
			args: args{fields: []string{"no_such"}},
			wantErr: true,
		},
		{
			name: "With only built-in fields",
			args: args{index: "idx1"},
			want: fmt.Sprintf(
				queryTemplate,
				"idx1",
				"",
				GenDedupEval([]string{
					"objectRef.apiGroup",
					"objectRef.apiVersion",
					"objectRef.resource",
					"verb",
					"auditID",
					"requestReceivedTimestamp",
				}),
				"",
			),
		},
		{
			name: "With an extra field and predicate",
			args: args{
				index: "idx2",
				predicate: "verb=create",
				fields: []string{"userId"},
			},
			want: fmt.Sprintf(
				queryTemplate,
				"idx2",
				"verb=create",
				GenDedupEval([]string{
					"objectRef.apiGroup",
					"objectRef.apiVersion",
					"objectRef.resource",
					"verb",
					"auditID",
					"requestReceivedTimestamp",
					"impersonatedUser.username",
					"user.username",
				}),
				`userId=if(isnull('impersonatedUser.username'),` +
				`'user.username','impersonatedUser.username'),`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserJourneyQueryGen(tt.args.index, tt.args.predicate, tt.args.fields)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
