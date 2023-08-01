package querygen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldSet_QueryGen(t *testing.T) {
	fldSet := FieldSet{
		"fallback":     {srcFields: []string{"orig_a1", "orig_a2"}},
		"fallback2":    {srcFields: []string{"orig_b1", "orig_b2", "orig_b3"}},
		"fixed_expr":   {srcExpr: `"foo"`},
		"fixed_expr2":  {srcExpr: `'e1'+'e2'`, srcFields: []string{"e1", "e2"}},
		"plain_field":  {},
		"plain_field2": {},
		"renamed_fld":  {srcFields: []string{"orig.field"}},
		"so_fld1":      {subObj: "sub_obj"},
		"so_fld2":      {subObj: "sub_obj2"},
		"so_fld3":      {subObj: "sub_obj", srcFields: []string{"so_fld3_orig"}},
	}
	includeFieldsCmd := `fields fallback,fallback2,fixed_expr,fixed_expr2,` +
		`plain_field,plain_field2,renamed_fld,sub_obj,sub_obj2`
	type args struct {
		searchExpr string
		fields     []string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		want_error bool
	}{
		{
			name: "Single plain field",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"plain_field"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"plain_field"}) +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "Non-existent field",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"plain_field", "no_such_field"},
			},
			want_error: true,
		},
		{
			name: "Multiple plain fields",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"plain_field", "plain_field2"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"plain_field", "plain_field2"}) +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "With renamed field",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"plain_field", "renamed_fld"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"plain_field", "orig.field"}) +
				`|eval renamed_fld='orig.field'` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "With fallback fields",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"fallback", "fallback2"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"orig_a1", "orig_a2", "orig_b1", "orig_b2", "orig_b3"}) +
				`|eval ` +
				`fallback=if(isnull('orig_a1'),'orig_a2','orig_a1'),` +
				`fallback2=if(isnull('orig_b1'),if(isnull('orig_b2'),'orig_b3','orig_b2'),'orig_b1')` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "With expression fields",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"fixed_expr", "fixed_expr2"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"e1", "e2"}) +
				`|eval fixed_expr="foo",fixed_expr2='e1'+'e2'` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
		{
			name: "With sub objects",
			args: args{
				searchExpr: `search index="foo"`,
				fields:     []string{"plain_field", "so_fld1", "plain_field2", "so_fld2", "so_fld3"},
			},
			want: `search index="foo"` +
				`|` + GenDedupEval([]string{"plain_field", "so_fld1", "plain_field2", "so_fld2", "so_fld3_orig"}) +
				`|eval ` +
				`sub_obj=json_object("so_fld1",'so_fld1',"so_fld3",'so_fld3_orig'),` +
				`sub_obj2=json_object("so_fld2",'so_fld2')` +
				`|` + includeFieldsCmd +
				`|` + excludeFieldsCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fldSet.QueryGen(tt.args.searchExpr, tt.args.fields)
			if tt.want_error {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
