package querygen

import (
	"testing"
)

func TestGenDedupEval(t *testing.T) {
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
	expected := `eval ` +
		`auditID=mvindex('auditID', 0),` +
		`impersonatedUser.username=mvindex('impersonatedUser.username', 0),` +
		`user.username=mvindex('user.username', 0),` +
		`objectRef.resource=mvindex('objectRef.resource', 0),` +
		`objectRef.namespace=mvindex('objectRef.namespace', 0),` +
		`objectRef.apiGroup=mvindex('objectRef.apiGroup', 0),` +
		`objectRef.apiVersion=mvindex('objectRef.apiVersion', 0),` +
		`objectRef.name=mvindex('objectRef.name', 0),` +
		`verb=mvindex('verb', 0),` +
		`requestReceivedTimestamp=mvindex('requestReceivedTimestamp', 0)`
	out := GenDedupEval(fields)
	if out != expected {
		t.Errorf("GenDedupEval() returns:\n  %q, expected\n  %q", out, expected)
	}
}

func TestGenApplicationQuery(t *testing.T) {
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
	expected := `search ` +
		`index="some_index" ` +
		`log_type=audit ` +
		`NOT verb IN (get, watch, list, deletecollection) ` +
		`"responseStatus.code" IN (200, 201) ` +
		`"objectRef.apiGroup"="appstudio.redhat.com" ` +
		`"objectRef.resource"="applications" ` +
		`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) ` +
		`(verb!=create OR "responseObject.metadata.resourceVersion"="*")` +
		`|` + GenDedupEval(fields) +
		`|` + GenTrackFields(TrackFieldSpec{true, true, true}, json_properties)
	out := GenApplicationQuery("some_index")
	if out != expected {
		t.Errorf("GenApplicationQuery() = %v, want %v", out, expected)
	}
}

func TestGenPropertiesJSONExpr(t *testing.T) {
	type args struct {
		properties_map map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Application query properties",
			args: args{properties_map: map[string]string{
				"apiGroup":   "objectRef.apiGroup",
				"apiVersion": "objectRef.apiVersion",
				"kind":       "objectRef.resource",
				"name":       "objectRef.name",
			}},
			want: `json_object(` +
				`"apiGroup",'objectRef.apiGroup',` +
				`"apiVersion",'objectRef.apiVersion',` +
				`"kind",'objectRef.resource',` +
				`"name",'objectRef.name'` +
				`)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenPropertiesJSONExpr(tt.args.properties_map); got != tt.want {
				t.Errorf("GenPropertiesJSONExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenTrackFields(t *testing.T) {
	json_properties := map[string]string{
		"apiGroup":   "objectRef.apiGroup",
		"apiVersion": "objectRef.apiVersion",
		"kind":       "objectRef.resource",
		"name":       "objectRef.name",
	}
	type args struct {
		spec           TrackFieldSpec
		properties_map map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with userid, verb, subject & app props",
			args: args{
				spec: TrackFieldSpec{
					with_userid:     true,
					with_ev_verb:    true,
					with_ev_subject: true,
				},
				properties_map: json_properties,
			},
			want: `eval ` +
				`messageId=auditID,` +
				`timestamp=requestReceivedTimestamp,` +
				`type="track",` +
				`userId=` + UserIdExpr + `,` +
				`event_verb=verb,` +
				`event_subject='objectRef.resource',` +
				`properties=` + GenPropertiesJSONExpr(json_properties) +
				`|` + IncludeFieldsCmd + `|` + ExcludeFieldsCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenTrackFields(tt.args.spec, tt.args.properties_map); got != tt.want {
				t.Errorf("GenTrackFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
