package querygen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestK8sAuditFieldSet_QueryGen(t *testing.T) {
	kfs := K8sAuditFieldSet{
		K8sApiId{}: {
			"plain_field":  {},
			"override_field": {},
		},
		K8sApiId{"api1.com", "SomeObj"}: {
			"override_field": {srcFields: []string{"other_field"}},
		},
	}
	fields := []string{"plain_field", "override_field"}
	includeFieldsCmd := `fields override_field,plain_field`
	type args struct {
		index      string
		api        K8sApiId
		searchExpr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple query",
			args: args{
				index: "some_idx",
				api: K8sApiId{"other.api.com", "SomeOtherObj"},
				searchExpr: "foo bar baz",
			},
			want: `search index="some_idx" log_type=audit ` +
				`"objectRef.apiGroup"="other.api.com" `+
				`"objectRef.resource"="SomeOtherObj" ` +
				`foo bar baz` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "Query on customized object",
			args: args{
				index: "some_idx",
				api: K8sApiId{"api1.com", "SomeObj"},
				searchExpr: "foo bar baz",
			},
			want: `search index="some_idx" log_type=audit ` +
				`"objectRef.apiGroup"="api1.com" `+
				`"objectRef.resource"="SomeObj" ` +
				`foo bar baz` +
				`|eval override_field='other_field'` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kfs.QueryGen(tt.args.index, tt.args.api, tt.args.searchExpr, fields)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
