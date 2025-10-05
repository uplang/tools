# Setup Guide for Maintainers

This guide is for repository maintainers to configure the Homebrew tap integration.

## Prerequisites

### 1. Create the Homebrew Tap Repository

If it doesn't already exist, create a repository at `https://github.com/uplang/brew`:

```bash
# Create new repository
mkdir brew
cd brew

# Initialize
git init
mkdir Formula
echo "# Homebrew Formulas for UP Language" > README.md
echo "This repository provides Homebrew formulas for UP Language tools." >> README.md
echo "" >> README.md
echo "## Usage" >> README.md
echo "" >> README.md
echo "\`\`\`bash" >> README.md
echo "brew tap uplang/brew" >> README.md
echo "brew install up" >> README.md
echo "\`\`\`" >> README.md

# Add .gitignore
cat > .gitignore << 'EOF'
.DS_Store
*.swp
*~
EOF

# Commit and push
git add .
git commit -m "Initial commit"
git remote add origin git@github.com:uplang/brew.git
git push -u origin main
```

### 2. Create GitHub Personal Access Token

#### Option A: Fine-Grained Token (Recommended)

Fine-grained tokens are more secure as they limit access to specific repositories.

1. Go to GitHub Settings → Developer settings → Personal access tokens → Fine-grained tokens
2. Click "Generate new token"
3. Token name: `GoReleaser Homebrew Tap`
4. Expiration: Set as needed (or no expiration for automation)
5. **Repository access**: Select "Only select repositories"
   - Choose: `uplang/brew`
6. **Permissions** (Repository permissions):
   - ✅ `Contents`: Read and write
   - ✅ `Metadata`: Read-only (automatically selected)
7. Click "Generate token"
8. **IMPORTANT**: Copy the token immediately (you won't see it again)

#### Option B: Classic Token (Legacy)

If you prefer classic tokens:

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Name: `GoReleaser Homebrew Tap`
4. Scopes needed:
   - ✅ `repo` (Full control of private repositories)
5. Click "Generate token"
6. **IMPORTANT**: Copy the token immediately (you won't see it again)

### 3. Add Token to Repository Secrets

1. Go to `https://github.com/uplang/tools`
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add the following secret:
   - **Name**: `HOMEBREW_TAP_TOKEN`
   - **Value**: The personal access token you just created
5. Click "Add secret"

## What's Configured

The following has been configured in this repository:

### 1. GoReleaser Configuration
- **File**: `up/.goreleaser.yaml`
- **Features**:
  - Multi-platform builds (Linux, macOS, Windows)
  - Multi-architecture (amd64, arm64)
  - Automatic Homebrew formula generation
  - GitHub Release creation with binaries
  - Checksums and signatures

### 2. GitHub Actions Workflow
- **File**: `.github/workflows/release.yml`
- **Triggers**: On tag push matching `v*`
- **Actions**:
  - Builds binaries for all platforms
  - Creates GitHub Release
  - Uploads artifacts
  - Updates Homebrew tap automatically

### 3. Version Command
- **File**: `up/main.go`
- **Features**:
  - `up version` - Show version info
  - `up version --short` - Show version number only
  - `up version --json` - JSON output
  - Version set at build time via ldflags

### 4. Documentation
- **RELEASE.md**: Complete release process guide
- **HOMEBREW.md**: User-facing Homebrew installation guide
- Updated README files with Homebrew installation instructions

## Creating Your First Release

Once the setup is complete:

```bash
# Tag a version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions will automatically:
1. Build binaries for all platforms
2. Create a GitHub Release
3. Update the Homebrew tap

## Verifying the Setup

### 1. Check Secrets
Ensure the secret is set:
```bash
# Go to: https://github.com/uplang/tools/settings/secrets/actions
# Verify: HOMEBREW_TAP_TOKEN is listed
```

### 2. Test GoReleaser Locally
GoReleaser is included as a tool dependency in `go.mod`:
```bash
cd up

# Check the configuration
go tool goreleaser check

# Test a release build without publishing
go tool goreleaser release --snapshot --clean

# Check the generated artifacts
ls -lh dist/
```

### 3. Verify Homebrew Tap Access
Make sure the token has write access:
```bash
# Clone the brew repo with the token
git clone https://TOKEN@github.com/uplang/brew.git test-brew
cd test-brew
echo "test" > test.txt
git add test.txt
git commit -m "Test commit"
git push
cd ..
rm -rf test-brew

# If this works, the token has proper access
# Don't forget to delete the test commit from GitHub if needed
```

## Troubleshooting

### Release Workflow Fails
- Check GitHub Actions logs
- Verify secrets are set correctly
- Ensure tests pass before releasing

### Homebrew Formula Not Created
- Verify `HOMEBREW_TAP_TOKEN` has proper permissions:
  - **Fine-grained token**: `Contents: Read and write` for `uplang/brew`
  - **Classic token**: `repo` scope
- Check that `uplang/brew` repository exists
- Review GoReleaser logs in the workflow

### Token Expired
- Create a new token (steps above)
- Update the `HOMEBREW_TAP_TOKEN` secret

## Maintenance

### Token Rotation
If the token needs to be rotated:
1. Create a new token (see step 2)
2. Update the `HOMEBREW_TAP_TOKEN` secret (see step 3)
3. The next release will use the new token

### Updating GoReleaser Config
Edit `up/.goreleaser.yaml` and test with:
```bash
cd up
go tool goreleaser check
```

### Adding New Tools
To add Homebrew support for other tools (e.g., `language-server`, `repl`):
1. Copy `up/.goreleaser.yaml` to the tool directory
2. Update project name and binary name
3. Add release job in `.github/workflows/release.yml`

## Security Notes

- The `HOMEBREW_TAP_TOKEN` has write access to the `brew` repository
- **Recommended**: Use fine-grained tokens limited to only `uplang/brew` repository
- Never commit the token in code
- Use GitHub Secrets for storing the token
- Rotate the token periodically for security
- Token should have minimal necessary permissions (Contents: Read and write)

## Support

For issues with setup:
- Open an issue at https://github.com/uplang/tools/issues
- Check GoReleaser docs: https://goreleaser.com/
- Check Homebrew docs: https://docs.brew.sh/

## License

This setup and documentation is part of the UP Tools project, licensed under GNU GPL v3.0.

