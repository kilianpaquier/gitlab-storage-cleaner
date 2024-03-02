package mocks

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// MockPages simplifies mocking of objects paging for input url.
//
// It mocks two pages, the first one with input s slice and the second one with an empty slice.
func MockPages[S ~[]E, E any](t testing.TB, url string, s S) {
	t.Cleanup(httpmock.Reset)

	// mock first page
	httpmock.RegisterMatcherResponder(http.MethodGet, url,
		httpmock.NewMatcher("with_page_one", func(req *http.Request) bool {
			return req.URL.Query().Get("page") == "1"
		}),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, s))

	// mock second page
	httpmock.RegisterMatcherResponder(http.MethodGet, url,
		httpmock.NewMatcher("with_page_two", func(req *http.Request) bool {
			return req.URL.Query().Get("page") == "2"
		}),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, []E{}))
}
