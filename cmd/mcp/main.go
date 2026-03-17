// cmd/mcp/main.go — stdio-to-SSE proxy for Claude Desktop.
//
// Claude Desktop launches this binary as a child process. Each JSON-RPC line
// from stdin is POSTed to https://forge-cms.dev/mcp/message with
// Authorization: Bearer $MCP_TOKEN. The response is written to stdout.
//
// Build: go build -o forge-mcp-proxy.exe ./cmd/mcp/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const mcpURL = "https://forge-cms.dev/mcp/message"

// jsonRPCError wraps a non-JSON-RPC HTTP error into a valid JSON-RPC 2.0
// error response so Claude Desktop can display it without a schema crash.
func jsonRPCError(id json.RawMessage, code int, msg string) []byte {
	if id == nil {
		id = json.RawMessage("null")
	}
	type errObj struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	type resp struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Error   errObj          `json:"error"`
	}
	b, _ := json.Marshal(resp{
		JSONRPC: "2.0",
		ID:      id,
		Error:   errObj{Code: code, Message: msg},
	})
	return b
}

func main() {
	token := os.Getenv("MCP_TOKEN")
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1<<20)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		// Extract id from the request for error responses.
		// Notifications have no "id" — responses must be suppressed for them.
		var peek struct {
			ID json.RawMessage `json:"id"`
		}
		_ = json.Unmarshal(line, &peek)
		isRequest := len(peek.ID) > 0 && string(peek.ID) != "null"

		req, err := http.NewRequest("POST", mcpURL, bytes.NewReader(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "forge-mcp proxy: %v\n", err)
			if isRequest {
				fmt.Fprintf(os.Stdout, "%s\n", jsonRPCError(peek.ID, -32603, err.Error()))
			}
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "forge-mcp proxy: %v\n", err)
			if isRequest {
				fmt.Fprintf(os.Stdout, "%s\n", jsonRPCError(peek.ID, -32603, err.Error()))
			}
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()

		if !isRequest {
			// Notification: discard any server response, write nothing to stdout.
			continue
		}

		trimmed := bytes.TrimSpace(body)
		// If the response is not JSON (e.g. "Unauthorized" plain text), wrap it.
		if len(trimmed) == 0 || trimmed[0] != '{' {
			msg := strings.TrimSpace(string(trimmed))
			if msg == "" {
				msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
			}
			fmt.Fprintf(os.Stderr, "forge-mcp proxy: HTTP %d: %s\n", resp.StatusCode, msg)
			fmt.Fprintf(os.Stdout, "%s\n", jsonRPCError(peek.ID, -32603, msg))
			continue
		}
		fmt.Fprintf(os.Stdout, "%s\n", trimmed)
	}
}
