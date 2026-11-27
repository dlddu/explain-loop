package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(App())
	defer srv.Close()

	cases := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/", "", http.StatusOK},
		{http.MethodGet, "/healthz", "", http.StatusOK},
		{http.MethodGet, "/readyz", "", http.StatusOK},
		{http.MethodPost, "/mcp", `{"jsonrpc":"2.0","id":1,"method":"ping"}`, http.StatusOK},
		{http.MethodGet, "/mcp", "", http.StatusMethodNotAllowed},
		{http.MethodGet, "/nope", "", http.StatusNotFound},
	}
	for _, tc := range cases {
		req, _ := http.NewRequest(tc.method, srv.URL+tc.path, strings.NewReader(tc.body))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", tc.method, tc.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != tc.want {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, resp.StatusCode, tc.want)
		}
	}
}
