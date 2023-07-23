// Package queryprint contains utilities for printing one or more Splunk queries
package queryprint

import (
	"strings"
	"testing"

	. "github.com/lithammer/dedent"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

func TestPrettyPrintQueries(t *testing.T) {
	tests := []struct {
		name    string
		queries []QueryDesc
		want    string
	}{
		{
			name:    "With no queries",
			queries: []QueryDesc{},
			want:    "",
		},
		{
			name: "With a few short queries",
			queries: []QueryDesc{
				{"foo", "search index=foo"},
				{"foo count", "search index=foo | stats count by bar"},
				{"foo baz", "search index=foo bar=baz | fields bar, bal"},
			},
			want: strings.TrimSpace(Dedent(`
				foo
				---
				    search index=foo

				foo count
				---------
				    search index=foo | stats count by bar

				foo baz
				-------
				    search index=foo bar=baz | fields bar, bal`,
			)),
		},
		{
			name: "With a long query",
			queries: []QueryDesc{{
				"Some long query",
				`search index=some_long_index_name log_type=awesome match=value` +
					`|eval custom_field=some_expression,` +
					`other_field=other_expression` +
					`|fields fields,shown,in,results`,
			}},
			want: strings.TrimSpace(Dedent(`
				Some long query
				---------------
				    search index=some_long_index_name log_type=awesome match=value
				    |eval custom_field=some_expression,other_field=other_expression
				    |fields fields,shown,in,results
			`)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrettyPrintQueries(tt.queries); got != tt.want {
				diff := dmp.DiffPrettyText(dmp.DiffMain(tt.want, got, true))
				t.Errorf(
					"PrettyPrintQueries() incorrect output. "+
						"Diff:\n%v\nGot:\n%v\nWant:\n%v",
					diff, got, tt.want,
				)
			}
		})
	}
}

func TestMachinePrintQueries(t *testing.T) {
	tests := []struct {
		name string
		queries []QueryDesc
		want string
	}{
		{
			name:    "With no queries",
			queries: []QueryDesc{},
			want:    "",
		},
		{
			name: "With a few short queries",
			queries: []QueryDesc{
				{"foo", "search index=foo"},
				{"foo count", "search index=foo | stats count by bar"},
				{"foo baz", "search index=foo bar=baz | fields bar, bal"},
			},
			want: "search index=foo\x00" +
				"search index=foo | stats count by bar\x00" +
				"search index=foo bar=baz | fields bar, bal",
		},
		{
			name: "With a long query",
			queries: []QueryDesc{{
				"Some long query",
				`search index=some_long_index_name log_type=awesome match=value` +
					`|eval custom_field=some_expression,` +
					`other_field=other_expression` +
					`|fields fields,shown,in,results`,
			}},
			want: `search index=some_long_index_name log_type=awesome match=value` +
			`|eval custom_field=some_expression,` +
			`other_field=other_expression` +
			`|fields fields,shown,in,results`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MachinePrintQueries(tt.queries); got != tt.want {
				t.Errorf("MachinePrintQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}
