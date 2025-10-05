# UP REPL

Interactive Read-Eval-Print Loop for UP (Unified Properties).

## Installation

```bash
go install github.com/uplang/tools/repl@latest
```

Or use via the UP CLI:
```bash
up repl
```

## Features

- ✅ Interactive UP document editing
- ✅ Real-time parsing and validation
- ✅ Multiline input support
- ✅ Block and list detection
- ✅ Syntax highlighting (coming soon)
- ✅ History support (readline)

## Usage

### Direct

```bash
up-repl
```

### Via UP CLI

```bash
up repl
```

### With Options

```bash
# Enable debug output
up repl --debug

# Or directly
up-repl --debug
```

## REPL Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `.help` | `.h` | Show help message |
| `.exit` | `.quit`, `.q` | Exit REPL |
| `.clear` | `.c` | Clear input buffer |

## Examples

### Simple Values

```
up> name John Doe
Parsed:
  name: "John Doe"

up> age!int 30
Parsed:
  age!int: "30"
```

### Blocks

```
up> config {
... debug!bool true
... port!int 8080
... }
Parsed:
  config: { 2 entries }
```

### Lists

```
up> tags [
... important
... urgent
... ]
Parsed:
  tags: [ 2 items ]
```

### Multiline Strings

```
up> description ```
... This is a long
... multiline description
... ```
Parsed:
  description: "This is a long\nmultiline description"
```

## Features in Detail

### Automatic Multiline Detection

The REPL automatically detects when you're entering a multiline block or list:

- Entering `{` switches to multiline mode (prompt changes to `...`)
- Entering `[` switches to multiline mode
- Closing `}` or `]` exits multiline mode and evaluates
- Entering ` ``` ` starts a multiline string

### Buffer Management

- Input accumulates in a buffer across multiple lines
- Use `.clear` to reset the buffer if you make a mistake
- The buffer automatically clears after successful evaluation

### Error Handling

Parse errors are displayed immediately:

```
up> invalid syntax here
Error: unexpected token at line 1
```

## Development

```bash
# Clone and build
git clone https://github.com/uplang/tools
cd tools/repl
go build -v .

# Run
./up-repl
```

## Integration

The REPL can be launched from the UP CLI:

```bash
up repl
```

This automatically finds and executes `up-repl` from your PATH or the same directory as the `up` binary.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../../spec/CONTRIBUTING.md) for guidelines.

## License

GNU General Public License v3.0 - see [LICENSE](../LICENSE) for details.

