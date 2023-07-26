// Package webfixture includes a simple HTTP-server test fixture that collects
// the requests made to it
package webfixture

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceRequestsFrom(t *testing.T) {
	tests := []struct {
		name         string
		test_func func(url string, c *http.Client)
		want []RequestTrace
	}{
		{
			name: "Single GET",
			test_func: func(url string, c *http.Client) {
				c.Get(url)
			},
			want: []RequestTrace{
				{"GET", "/", ""},
			},
		},
		{
			name: "Multiple GETs",
			test_func: func(url string, c *http.Client) {
				c.Get(url + "/foo")
				c.Get(url + "/bar")
			},
			want: []RequestTrace{
				{"GET", "/foo", ""},
				{"GET", "/bar", ""},
			},
		},
		{
			name: "POST request with a body",
			test_func: func(url string, c *http.Client) {
				c.Post(url, "text/json", strings.NewReader(
					`{"foo": "bar}`,
				))
			},
			want: []RequestTrace{
				{"POST", "/", `{"foo": "bar}`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequests := TraceRequestsFrom(tt.test_func)
			assert.Equal(t, tt.want, gotRequests)
		})
	}
}
