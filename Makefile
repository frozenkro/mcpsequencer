BIN=build
LOG=log
MCP_BIN=mcpsequencer-mcp
TUI_BIN=mcpsequencer-tui

build: build.mcp build.tui

build.mcp:
	go build -o ${BIN}/${MCP_BIN} cmd/mcp/main.go

run.mcp: 
	go run cmd/mcp/main.go

build.tui:
	go build -o ${BIN}/${TUI_BIN} cmd/tui/main.go

run.tui:
	go run cmd/tui/main.go

test:
	go test ./...

debug.mcp:
	# # TODO figure out how to reliably kill this npx process when done with it.
	# mkdir -p ${LOG}
	# nohup npx @modelcontextprotocol/inspector > ${LOG}/inspector.log 2>&1 &
	dlv debug cmd/mcp/main.go -- --dev --http

debug.mcp.prod:
	# mkdir -p ${LOG}
	# nohup npx @modelcontextprotocol/inspector > ${LOG}/inspector.log 2>&1 &
	dlv debug cmd/mcp/main.go -- --http
