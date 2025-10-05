# UP CLI

Command-line interface for parsing, formatting, and working with UP (Unified Properties) documents.

## Installation

### From Go

```bash
go install github.com/uplang/tools/up@latest
```

### From Source

```bash
git clone https://github.com/uplang/tools
cd tools/up
go build -o up .
```

## Usage

### Parse

Parse a UP document and output as JSON:

```bash
up parse -i config.up
```

Options:
- `-i, --input FILE` - Input UP file (default: stdin)
- `-o, --output FILE` - Output file (default: stdout)
- `--pretty` - Pretty-print JSON output
- `--validate` - Validate against schema if present

### Format

Format a UP document with consistent style:

```bash
up format -i config.up -o formatted.up
```

Options:
- `-i, --input FILE` - Input UP file (default: stdin)
- `-o, --output FILE` - Output file (default: stdout)
- `--indent N` - Indentation spaces (default: 2)
- `--sort-keys` - Sort keys alphabetically

### Validate

Validate a UP document against a schema:

```bash
up validate -i config.up -s schema.up-schema
```

Options:
- `-i, --input FILE` - Input UP file (required)
- `-s, --schema FILE` - Schema file (default: auto-detect)
- `--strict` - Strict validation mode

### Evaluate

Evaluate dynamic namespaces in a UP document:

```bash
up eval -i config.up
```

Options:
- `-i, --input FILE` - Input UP file (default: stdin)
- `-o, --output FILE` - Output file (default: stdout)
- `--ns-path DIR` - Namespace search path (default: ./up-namespaces)
- `--pretty` - Pretty-print output

### Convert

Convert between UP and other formats:

```bash
# UP to JSON
up convert -i config.up -o config.json --to json

# JSON to UP
up convert -i config.json -o config.up --from json

# UP to YAML
up convert -i config.up -o config.yaml --to yaml
```

Options:
- `-i, --input FILE` - Input file (required)
- `-o, --output FILE` - Output file (required)
- `--from FORMAT` - Input format (json, yaml, toml) - auto-detected if not specified
- `--to FORMAT` - Output format (json, yaml, toml, up) - required
- `--pretty` - Pretty-print output

## Examples

### Basic Parsing

```bash
# Parse and pretty-print
echo "name John Doe" | up parse --pretty

# Parse file
up parse -i examples/server.up --pretty
```

### Format Documents

```bash
# Format and write to new file
up format -i messy.up -o clean.up

# Format in-place
up format -i config.up -o config.up

# Format with sorted keys
up format -i config.up --sort-keys
```

### Validation

```bash
# Validate with auto-detected schema
up validate -i config.up

# Validate with explicit schema
up validate -i config.up -s config.up-schema

# Strict validation
up validate -i config.up --strict
```

### Dynamic Evaluation

```bash
# Evaluate namespaces
up eval -i dynamic.up --pretty

# Custom namespace path
up eval -i dynamic.up --ns-path ./custom-namespaces
```

### Format Conversion

```bash
# Convert UP to JSON
up convert -i config.up -o config.json --to json --pretty

# Convert JSON to UP
up convert -i config.json -o config.up --to up

# Convert UP to YAML
up convert -i config.up -o config.yaml --to yaml
```

## Exit Codes

- `0` - Success
- `1` - Parse error
- `2` - Validation error
- `3` - I/O error
- `4` - Invalid arguments

## Environment Variables

- `UP_NS_PATH` - Default namespace search path
- `UP_CONFIG` - Default config file location

## Testing

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../../spec/CONTRIBUTING.md) for guidelines.

## License

GNU General Public License v3.0 - see [LICENSE](../LICENSE) for details.

