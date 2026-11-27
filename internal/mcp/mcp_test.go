package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func call(t *testing.T, body string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	rec := httptest.NewRecorder()
	NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body)))
	if rec.Body.Len() == 0 {
		return rec, nil
	}
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v: %s", err, rec.Body.String())
	}
	return rec, out
}

func TestInitialize(t *testing.T) {
	_, out := call(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	result, ok := out["result"].(map[string]any)
	if !ok {
		t.Fatalf("no result: %v", out)
	}
	if result["protocolVersion"] != ProtocolVersion {
		t.Errorf("protocolVersion = %v", result["protocolVersion"])
	}
	info := result["serverInfo"].(map[string]any)
	if info["name"] != "explain-loop" {
		t.Errorf("serverInfo.name = %v", info["name"])
	}
}

func TestToolsListEmpty(t *testing.T) {
	_, out := call(t, `{"jsonrpc":"2.0","id":"a","method":"tools/list"}`)
	tools := out["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 0 {
		t.Errorf("tools = %v, want empty", tools)
	}
	if out["id"] != "a" {
		t.Errorf("id not echoed: %v", out["id"])
	}
}

func TestNotificationAccepted(t *testing.T) {
	rec, out := call(t, `{"jsonrpc":"2.0","method":"notifications/initialized"}`)
	if rec.Code != http.StatusAccepted || out != nil {
		t.Errorf("code=%d body=%v, want 202 with empty body", rec.Code, out)
	}
}

func TestErrors(t *testing.T) {
	cases := map[string]struct {
		body string
		code float64
	}{
		"parse":          {`{not json`, codeParseError},
		"batch":          {`[{"jsonrpc":"2.0","id":1,"method":"ping"}]`, codeInvalidRequest},
		"version":        {`{"jsonrpc":"1.0","id":1,"method":"ping"}`, codeInvalidRequest},
		"unknown method": {`{"jsonrpc":"2.0","id":1,"method":"nope"}`, codeMethodNotFound},
		"unknown tool":   {`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"x"}}`, codeInvalidParams},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, out := call(t, tc.body)
			e, ok := out["error"].(map[string]any)
			if !ok || e["code"] != tc.code {
				t.Errorf("got %v, want error code %v", out, tc.code)
			}
		})
	}
}
