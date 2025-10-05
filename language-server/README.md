# UP Language Server

Language Server Protocol (LSP) implementation for UP (Unified Properties).

## Features

- ✅ **Syntax Validation** - Real-time error detection
- ✅ **Auto-completion** - Type annotations and namespace functions
- ✅ **Hover Information** - Documentation on hover
- ✅ **Go to Definition** - Navigate to definitions
- ✅ **Document Symbols** - Outline view
- ✅ **Formatting** - Auto-format documents
- 🚧 **Diagnostics** - Advanced linting (planned)
- 🚧 **Rename** - Symbol renaming (planned)

## Installation

```bash
go install github.com/uplang/tools/language-server@latest
```

## Usage

### Command Line

```bash
# Start the language server
up-language-server

# With debug logging
up-language-server -debug -log /tmp/up-lsp.log
```

### Editor Integration

#### VS Code

Install the UP extension (coming soon) or configure manually:

```json
{
  "up.languageServer": {
    "command": "up-language-server",
    "args": ["-log", "/tmp/up-lsp.log"]
  }
}
```

#### Neovim

```lua
require'lspconfig'.configs.up = {
  default_config = {
    cmd = {'up-language-server'},
    filetypes = {'up'},
    root_dir = function(fname)
      return vim.fn.getcwd()
    end,
  },
}

require'lspconfig'.up.setup{}
```

#### Emacs (lsp-mode)

```elisp
(lsp-register-client
 (make-lsp-client
  :new-connection (lsp-stdio-connection '("up-language-server"))
  :major-modes '(up-mode)
  :server-id 'up-ls))
```

## Development

### Building

```bash
git clone https://github.com/uplang/tools
cd tools/language-server
go build -o up-language-server .
```

### Testing

```bash
go test ./...
```

### Architecture

```
language-server/
├── main.go           # Entry point
├── server/
│   ├── server.go     # LSP server implementation
│   ├── completion.go # Auto-completion
│   ├── hover.go      # Hover information
│   ├── diagnostics.go # Error detection
│   └── formatting.go  # Document formatting
└── protocol/
    └── types.go      # LSP protocol types
```

## Protocol Support

Implements LSP version 3.17:

| Feature | Status |
|---------|--------|
| textDocument/didOpen | ✅ |
| textDocument/didChange | ✅ |
| textDocument/didClose | ✅ |
| textDocument/completion | ✅ |
| textDocument/hover | ✅ |
| textDocument/definition | ✅ |
| textDocument/documentSymbol | ✅ |
| textDocument/formatting | ✅ |
| textDocument/rename | 🚧 |
| textDocument/codeAction | 🚧 |

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../spec/CONTRIBUTING.md) for guidelines.

## License

GNU General Public License v3.0 - see [LICENSE](LICENSE) for details.

## Links

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [UP Language Specification](https://github.com/uplang/spec)

