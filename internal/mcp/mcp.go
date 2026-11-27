// Package mcp serves the Model Context Protocol over Streamable HTTP (JSON responses only).
//
// 골격 단계라 도구는 아직 없다. PRD-1~4 의 기능은 registry 에 도구를 등록하는 방식으로 붙인다.
package mcp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dlddu/explain-loop/internal/version"
)

// ProtocolVersion is the MCP revision this server speaks (same as homelab-k3s-mcp).
const ProtocolVersion = "2025-06-18"

// maxBodyBytes caps a single JSON-RPC request.
const maxBodyBytes = 1 << 20

// JSON-RPC error codes.
const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
)

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool is an MCP tool definition as returned by tools/list.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

// Handler answers POST /mcp.
type Handler struct {
	tools []Tool
}

// NewHandler returns a handler with no tools registered.
func NewHandler() *Handler {
	return &Handler{tools: []Tool{}}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		writeJSON(w, errorResponse(nullID(), codeParseError, "parse error"))
		return
	}
	if len(body) > maxBodyBytes {
		writeJSON(w, errorResponse(nullID(), codeInvalidRequest, "request too large"))
		return
	}
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		writeJSON(w, errorResponse(nullID(), codeInvalidRequest, "batch requests are not supported"))
		return
	}

	var req request
	if err := json.Unmarshal(trimmed, &req); err != nil {
		writeJSON(w, errorResponse(nullID(), codeParseError, "parse error"))
		return
	}

	// 알림(id 없음)과 클라이언트가 보낸 응답에는 본문 없이 202 를 돌려준다(Streamable HTTP).
	if len(bytes.TrimSpace(req.ID)) == 0 {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if req.JSONRPC != "2.0" {
		writeJSON(w, errorResponse(req.ID, codeInvalidRequest, "invalid jsonrpc version"))
		return
	}

	result, rerr := h.dispatch(req.Method)
	if rerr != nil {
		writeJSON(w, response{JSONRPC: "2.0", ID: req.ID, Error: rerr})
		return
	}
	writeJSON(w, response{JSONRPC: "2.0", ID: req.ID, Result: result})
}

func (h *Handler) dispatch(method string) (any, *rpcError) {
	switch method {
	case "initialize":
		return map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{"listChanged": false}},
			"serverInfo":      map[string]any{"name": version.Name, "version": version.Version},
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": h.tools}, nil
	case "tools/call":
		return nil, &rpcError{Code: codeInvalidParams, Message: "unknown tool"}
	default:
		return nil, &rpcError{Code: codeMethodNotFound, Message: "method not found: " + method}
	}
}

func errorResponse(id json.RawMessage, code int, msg string) response {
	return response{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}}
}

func nullID() json.RawMessage { return json.RawMessage("null") }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
