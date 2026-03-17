// cmd/mcp/main.go — stdio-to-SSE proxy for Claude Desktop.
//
// Claude Desktop launches this binary as a child process and communicates
// over stdin/stdout using newline-delimited JSON-RPC. Each line is forwarded
// to the remote Forge MCP endpoint (POST /mcp/message) with a Bearer token,
// and the JSON response is written back to stdout.
//
// Usage:
//
//	MCP_TOKEN=<bearer> forge-mcp-proxy
//
// Build:
//
//	go build -o forge-mcp-proxy.exe ./cmd/mcp/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const mcpURL = "https://forge-cms.dev/mcp/message"

func main() {
	token := os.Getenv("MCP_TOKEN")
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		req, err := http.NewRequest("POST", mcpURL, bytes.NewReader(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "forge-mcp proxy: %v\n", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "forge-mcp proxy: %v\n", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Fprintf(os.Stdout, "%s\n", bytes.TrimSpace(body))
	}
}
