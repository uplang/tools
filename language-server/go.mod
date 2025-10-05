module github.com/uplang/tools/language-server

go 1.25

toolchain go1.25.1

require (
	github.com/urfave/cli/v2 v2.27.1
	go.lsp.dev/protocol v0.12.0
)

tool (
	github.com/golangci/golangci-lint/cmd/golangci-lint
)

