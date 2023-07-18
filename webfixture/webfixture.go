// Package webfixture includes a simple HTTP-server test fixture that collects
// the requests made to it
package webfixture

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// RequestTrace represents simple data about a web request
type RequestTrace struct {
	Method,	Path, Body string
}

// TraceRequestsFrom runs a small web server, and then invokes a provided test
// function while passing it the server URL and an HTTP client. The requests
// made to the web server while the test function is running are then logged
// and returned.
func TraceRequestsFrom(test_func func(url string, c *http.Client)) (requests []RequestTrace) {
	requests_chan := make(chan RequestTrace)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		requests_chan <- RequestTrace{r.Method, r.URL.Path, string(body)}
	}))
	defer svr.Close()
	defer close(requests_chan)
	go func() {
		for request := range requests_chan {
			requests = append(requests, request)
		}
	}()
	test_func(svr.URL, svr.Client())
	return
}