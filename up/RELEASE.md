# Release Guide for UP CLI

This guide explains how to release a new version of the UP CLI tool.

## Prerequisites

### 1. GitHub Secret for Homebrew Tap

To enable automatic Homebrew formula updates, you need to set up a GitHub token:

1. Create a Personal Access Token (PAT) on GitHub:
   - Go to Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Click "Generate new token (classic)"
   - Give it a descriptive name like "GoReleaser Homebrew Tap"
   - Set expiration as needed (or no expiration for automation)
   - Select scopes:
     - `repo` (Full control of private repositories)
     - `write:packages` (if using GitHub Packages)
   - Click "Generate token" and copy the token

2. Add the token as a repository secret:
   - Go to your repository: `https://github.com/uplang/tools`
   - Navigate to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `HOMEBREW_TAP_TOKEN`
   - Value: Paste the token you copied
   - Click "Add secret"

### 2. Homebrew Tap Repository

Ensure the Homebrew tap repository exists:
- Repository: `https://github.com/uplang/homebrew-tap`
- Should have a `Formula/` directory
- The GoReleaser will automatically create/update the formula file

If it doesn't exist:
```bash
# Create the tap repository
mkdir homebrew-tap
cd homebrew-tap
mkdir Formula
echo "# Homebrew Tap for UP Language Tools" > README.md
git init
git add .
git commit -m "Initial commit"
# Push to GitHub as uplang/homebrew-tap
```

## Release Process

### 1. Prepare the Release

Before tagging, ensure:
- All tests pass: `go test ./...`
- Code is linted: `go tool golangci-lint run ./...`
- README and documentation are up to date
- CHANGELOG is updated (if maintained)

### 2. Create and Push a Tag

```bash
# Make sure you're on main branch and up to date
git checkout main
git pull origin main

# Create an annotated tag
git tag -a v1.0.0 -m "Release v1.0.0: Description of changes"

# Push the tag to GitHub
git push origin v1.0.0
```

### 3. Automated Release

Once you push the tag, GitHub Actions will automatically:

1. **Build** binaries for all platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)

2. **Create** a GitHub Release:
   - Release title: `up v1.0.0`
   - Includes all binaries as downloadable assets
   - Generates checksums for verification
   - Includes a changelog

3. **Update** the Homebrew tap:
   - Creates/updates `Formula/up.rb` in `uplang/homebrew-tap`
   - Includes SHA256 checksums for the archives
   - Sets the correct version and download URLs

### 4. Verify the Release

After the GitHub Action completes:

1. Check the [releases page](https://github.com/uplang/tools/releases)
   - Verify all binaries are present
   - Verify checksums.txt exists

2. Test the Homebrew installation:
   ```bash
   # Update tap
   brew update
   brew upgrade uplang/tap/up

   # Or fresh install
   brew install uplang/tap/up

   # Verify version
   up version
   ```

3. Test a binary download:
   ```bash
   # Download for your platform
   curl -LO https://github.com/uplang/tools/releases/download/v1.0.0/up_1.0.0_darwin_arm64.tar.gz

   # Extract and test
   tar xzf up_1.0.0_darwin_arm64.tar.gz
   ./up version
   ```

## Testing Releases Locally

Before creating a real release, test the GoReleaser configuration:

```bash
# GoReleaser is already available as a Go tool (defined in go.mod)
# No separate installation needed!

# Test the release build (creates dist/ but doesn't publish)
go tool goreleaser release --snapshot --clean

# Check the generated artifacts
ls -lh dist/
```

This creates binaries in the `dist/` directory without publishing anything.

**Note**: GoReleaser is included as a tool dependency in `go.mod`, so it's automatically available via `go tool goreleaser`. No separate installation is required!

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v2.0.0): Incompatible API changes
- **MINOR** version (v1.1.0): New functionality, backwards compatible
- **PATCH** version (v1.0.1): Bug fixes, backwards compatible

### Pre-releases

For pre-releases, use suffixes:
- Alpha: `v1.0.0-alpha.1`
- Beta: `v1.0.0-beta.1`
- Release Candidate: `v1.0.0-rc.1`

GoReleaser will automatically mark these as pre-releases on GitHub.

## Troubleshooting

### GoReleaser Fails

Check the GitHub Actions logs:
1. Go to Actions tab in your repository
2. Find the failed workflow run
3. Check the "Release UP CLI" job logs

Common issues:
- **Tests fail**: Fix failing tests before releasing
- **Missing secrets**: Ensure `HOMEBREW_TAP_TOKEN` is set
- **Permission denied**: Check token permissions include `repo` scope

### Homebrew Formula Not Updated

If the Homebrew tap doesn't update:

1. Check the token has write access to `uplang/homebrew-tap`
2. Verify the repository exists and has the correct name
3. Check GoReleaser logs for brew-specific errors
4. Manually verify in the tap repository if a PR was created

### Version Command Shows Wrong Version

The version is set at build time via ldflags. If building manually:

```bash
VERSION=v1.0.0
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build -ldflags "\
  -X main.version=$VERSION \
  -X main.commit=$COMMIT \
  -X main.date=$DATE \
  -X main.builtBy=manual" \
  -o up .

./up version
```

## Additional Resources

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)

