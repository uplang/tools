module github.com/uplang/tools/examples

go 1.25

toolchain go1.25.1

require github.com/urfave/cli/v2 v2.27.1

tool (
	github.com/golangci/golangci-lint/cmd/golangci-lint
)

