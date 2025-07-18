#!/bin/bash
 
# interactive_test.sh
SERVER_CMD="go run main.go"
 
echo "Starting MCP STDIO server test..."
 
# Function to send JSON-RPC request
send_request() {
    local request="$1"
    echo "Sending: $request"
    echo "$request" | $SERVER_CMD
    echo "---"
}
 
# Initialize
send_request '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test","version":"1.0.0"}}}'
 
# List tools
send_request '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
 
# List resources
send_request '{"jsonrpc":"2.0","id":3,"method":"resources/list","params":{}}'
 
# Call tool
send_request '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"create_project","arguments":{"ProjectName":"TestProject1","Tasks":["Task 1","Task 2"]}}}'
 
echo "Test completed."
