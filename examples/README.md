# UP Examples Test Runner

Tool for running and testing UP example files across namespaces and projects.

## Installation

### From Go

```bash
go install github.com/uplang/tools/examples@latest
```

### From Source

```bash
git clone https://github.com/uplang/tools
cd tools/examples
go build -o examples .
```

## Usage

### Basic Usage

Run all examples in the current directory:

```bash
examples
```

Run examples in a specific directory:

```bash
examples /path/to/namespace
```

With verbose output:

```bash
examples -v
```

### Options

- `PATH` - Directory to search for examples (default: current directory)
- `-v, --verbose` - Enable verbose output

## How It Works

The tool:

1. **Searches** for all `.up` files in `examples/` subdirectories
2. **Detects** the namespace from the directory structure
3. **Parses** each file using the UP CLI (`up parse`)
4. **Reports** success/failure for each example
5. **Summarizes** results at the end

## Example Output

```
=========================================
UP Namespace Examples Runner
=========================================

Found 12 example files

[1/12] Testing: string - uppercase.up
File: ns/string/examples/uppercase.up

Output:
{
  "key": "result",
  "type": "string",
  "value": "HELLO WORLD"
}

---

[2/12] Testing: math - addition.up
File: ns/math/examples/addition.up

Output:
{
  "key": "sum",
  "type": "int",
  "value": 42
}

---

=========================================
All examples processed!
=========================================

Note: These examples show UP syntax parsing.
To evaluate with actual namespace functions:
1. Build namespaces: cd <namespace> && go build
2. Use UP CLI with full namespace support
```

## Directory Structure

The tool expects this directory structure:

```
namespace/
├── examples/
│   ├── example1.up
│   ├── example2.up
│   └── ...
├── namespace.up-schema
└── main.go (or binary)
```

Examples can be in any subdirectory named `examples/`.

## Requirements

- **UP CLI** must be installed and in PATH
  ```bash
  go install github.com/uplang/tools/up@latest
  ```

- Examples must be valid UP syntax
- For namespace evaluation, the namespace binary must be built

## Exit Codes

- `0` - All examples processed successfully
- `1` - Error scanning directories or no examples found
- `2` - UP CLI not found

## Environment Variables

- `UP_CLI` - Path to UP CLI binary (default: searches PATH for `up`)

## Testing

```bash
# Test from namespace directory
cd ns/string
examples

# Test all namespaces
cd ns
examples .

# Test specific namespace
examples ns/math
```

## Integration with CI

Example GitHub Actions workflow:

```yaml
name: Test Examples

on: [push, pull_request]

jobs:
  test-examples:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install UP CLI
        run: go install github.com/uplang/tools/up@latest

      - name: Install examples runner
        run: go install github.com/uplang/tools/examples@latest

      - name: Run examples
        run: examples
```

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../../spec/CONTRIBUTING.md) for guidelines.

## License

GNU General Public License v3.0 - see [LICENSE](../LICENSE) for details.

