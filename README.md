# UP Tools

Collection of tools and utilities for working with UP (Unified Properties).

Follows the Go command pattern: a single `up` command with subcommands and tool dispatch.

## Installation

```bash
# Install main UP CLI (includes built-in commands)
go install github.com/uplang/tools/up@latest

# Install additional tools (optional)
go install github.com/uplang/tools/lsp@latest
go install github.com/uplang/tools/repl@latest
go install github.com/uplang/tools/examples@latest
```

## Main Command: `up`

The `up` command works like the `go` command - a single entry point with subcommands:

```
up <command> [arguments]
```

### Built-in Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `parse` | `p` | Parse UP documents and output as JSON |
| `format` | `fmt`, `f` | Format UP documents with consistent style |
| `validate` | `vet`, `v` | Validate UP documents against schemas |
| `eval` | `e` | Evaluate dynamic namespaces |
| `convert` | `c` | Convert between UP and other formats |
| `lsp` | `language-server` | Start the UP language server |
| `repl` | - | Start interactive REPL |
| `tool` | - | Run specified UP tool |
| `version` | - | Print UP version |

### Usage Examples

```bash
# Parse and pretty-print
up parse -i config.up --pretty

# Format document (shorthand)
up fmt -i config.up -o formatted.up

# Validate syntax (shorthand)
up vet -i config.up

# Start REPL
up repl

# Start language server
up lsp --debug --log /tmp/up-lsp.log

# Run external tool
up tool my-analyzer config.up
```

See [up/README.md](up/README.md) for complete CLI documentation.

## Additional Tools

### Language Server (`language-server`)

Language Server Protocol implementation for UP.

**Standalone:**
```bash
go install github.com/uplang/tools/lsp@latest
up-language-server
```

**Via UP CLI:**
```bash
up lsp
```

See [lsp/README.md](lsp/README.md) for full documentation.

### Interactive REPL (`repl`)

Read-Eval-Print Loop for interactive UP editing.

**Standalone:**
```bash
go install github.com/uplang/tools/repl@latest
up-repl
```

**Via UP CLI:**
```bash
up repl
```

See [repl/README.md](repl/README.md) for full documentation.

### Examples Runner (`examples`)

Tool for running and testing UP example files across namespaces.

**Standalone:**
```bash
go install github.com/uplang/tools/examples@latest
examples /path/to/namespace
```

**Via UP CLI:**
```bash
up tool examples /path/to/namespace
```

See [examples/README.md](examples/README.md) for full documentation.

## Development

Each tool is a standalone Go module with its own `go.mod` file:

```
tools/
├── up/
│   ├── go.mod          # Standalone module
│   ├── main.go
│   └── README.md
├── examples/
│   ├── go.mod          # Standalone module
│   ├── main.go
│   └── README.md
└── README.md
```

### Building Tools

```bash
# Build all tools
cd tools
for tool in */; do
  if [ -f "$tool/go.mod" ]; then
    echo "Building $tool..."
    cd "$tool"
    go build -v .
    cd ..
  fi
done
```

### Testing Tools

```bash
# Test all tools
cd tools
for tool in */; do
  if [ -f "$tool/go.mod" ]; then
    echo "Testing $tool..."
    cd "$tool"
    go test -v ./...
    cd ..
  fi
done
```

## Releases

Each tool uses [GoReleaser](https://goreleaser.com/) for automated releases.

### Creating a Release

1. Tag a new version:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions automatically:
   - Builds binaries for Linux, macOS, Windows (amd64 + arm64)
   - Creates a GitHub release
   - Uploads all artifacts with checksums
   - Updates Homebrew tap (if configured)

### Testing Locally

Test a release build without publishing:

```bash
# Test UP CLI
cd up
goreleaser release --snapshot --clean

# Test examples runner
cd ../examples
goreleaser release --snapshot --clean
```

Artifacts will be in the `dist/` directory.

### Installation Methods

After a release, users can install via:

```bash
# Via go install
go install github.com/uplang/tools/up@latest
go install github.com/uplang/tools/examples@latest

# Via Homebrew (once tap is configured)
brew install uplang/tap/up
brew install uplang/tap/examples

# Direct download from GitHub releases
# https://github.com/uplang/tools/releases
```

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../spec/CONTRIBUTING.md) for guidelines.

## License

GNU General Public License v3.0 - see [LICENSE](LICENSE) for details.

## Links

- [UP Language Specification](https://github.com/uplang/spec)
- [UP Go Parser](https://github.com/uplang/go)
- [UP Namespaces](https://github.com/uplang/ns)
