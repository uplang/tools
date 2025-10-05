# Homebrew Installation Guide

## Installing UP CLI via Homebrew

The UP CLI can be easily installed on macOS and Linux using Homebrew.

### Quick Install

```bash
# Install directly (Homebrew will add the tap automatically on first install)
brew install uplang/brew/up
```

### Step-by-Step Install

```bash
# 1. Add the uplang brew tap to Homebrew
brew tap uplang/brew

# 2. Install the UP CLI
brew install up

# 3. Verify installation
up version
```

### What Gets Installed

When you install via Homebrew, you get:
- The `up` binary installed to `/opt/homebrew/bin/up` (Apple Silicon) or `/usr/local/bin/up` (Intel Mac)
- Automatic dependency management
- Easy upgrades via `brew upgrade`

### Upgrading

To upgrade to the latest version:

```bash
# Update Homebrew and upgrade up
brew update
brew upgrade uplang/brew/up
```

### Uninstalling

To remove the UP CLI:

```bash
# Uninstall the package
brew uninstall uplang/brew/up

# Optionally, remove the tap
brew untap uplang/brew
```

## Using the UP CLI

Once installed, you can use all `up` commands:

```bash
# Parse a UP document
up parse -i config.up --pretty

# Format a UP document
up format -i config.up -o formatted.up

# Validate a UP document
up validate -i config.up

# Get help
up --help
```

## For Linux Users

Homebrew also works on Linux! Install Homebrew first:

```bash
# Install Homebrew on Linux
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Then follow the same steps as macOS
brew install uplang/brew/up
```

## Troubleshooting

### Command Not Found

If you get `command not found: up` after installation:

1. Make sure Homebrew's bin directory is in your PATH:
   ```bash
   # For Apple Silicon
   echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc

   # For Intel Mac
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc

   # Reload shell
   source ~/.zshrc
   ```

2. Verify Homebrew installation:
   ```bash
   brew doctor
   ```

### Installation Fails

If installation fails:

1. Update Homebrew:
   ```bash
   brew update
   ```

2. Check for issues:
   ```bash
   brew doctor
   ```

3. Try installing with verbose output:
   ```bash
   brew install uplang/brew/up --verbose
   ```

### Wrong Version After Upgrade

If `up version` shows an old version after upgrading:

1. Clear cache:
   ```bash
   brew cleanup
   ```

2. Check which binary is being used:
   ```bash
   which up
   ```

3. Reinstall if needed:
   ```bash
   brew reinstall uplang/brew/up
   ```

## Alternative Installation Methods

If you prefer not to use Homebrew, see other installation options:

- **Go Install**: `go install github.com/uplang/tools/up@latest`
- **Pre-built Binaries**: Download from [GitHub Releases](https://github.com/uplang/tools/releases)
- **From Source**: Clone and build manually

See the [main README](README.md) for details on alternative installation methods.

## Homebrew Formula

The Homebrew formula is automatically generated and maintained by GoReleaser on each release.

You can view the formula at:
- https://github.com/uplang/brew/blob/main/Formula/up.rb

## Support

For issues related to:
- **Installation/Homebrew**: Open an issue at [uplang/brew](https://github.com/uplang/brew/issues)
- **UP CLI functionality**: Open an issue at [uplang/tools](https://github.com/uplang/tools/issues)

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](../spec/CONTRIBUTING.md) for guidelines.

## License

UP CLI is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for details.

