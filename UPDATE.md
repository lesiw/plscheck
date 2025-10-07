# Updating plscheck for New gopls Versions

This document describes the process for updating plscheck to match a new gopls release.

## Overview

plscheck externalizes gopls internal analyzers from `golang.org/x/tools/gopls/internal/analysis`. Each plscheck release corresponds exactly to a gopls release (e.g., plscheck v0.20.0 = gopls/v0.20.0).

## Update Process

### 1. Checkout Upstream gopls Version

```bash
# Assumes golang.org/x/tools is cloned locally
# Adjust TOOLS_PATH to your local clone location
TOOLS_PATH="${TOOLS_PATH:-$HOME/.local/src/github.com/golang/tools}"

cd "$TOOLS_PATH"
git fetch --tags
git checkout gopls/vX.Y.Z  # Replace with target version
```

### 2. Copy Analyzer Source Files

```bash
# From the plscheck repository root
cd "$PLSCHECK_ROOT"  # Or wherever lesiw.io/plscheck is located

# Remove old analyzers (preserve cmd/, internal/, .gitignore, README.md, go.mod, go.sum)
rm -rf deprecated embeddirective fillreturns fillstruct fillswitch infertypeargs \
       maprange modernize nonewvars noresultvalues recursiveiter \
       simplifycompositelit simplifyrange simplifyslice unusedfunc \
       unusedparams unusedvariable yield

# Copy all analyzers from upstream
for dir in "$TOOLS_PATH"/gopls/internal/analysis/*/; do
    cp -r "$dir"* "$(basename "$dir")/";
done

# Remove all test files and testdata directories
# Tests exist upstream and don't need to be duplicated here
find . -name "*_test.go" -type f -delete
find . -type d -name "testdata" -exec rm -rf {} + 2>/dev/null
```

### 3. Copy Internal Dependencies

Copy required internal packages from both gopls and x/tools:

```bash
# gopls internal dependencies
mkdir -p internal/{fuzzy,util/{astutil,bug,moreiters,moreslices,safetoken,tokeninternal}}

cp -r "$TOOLS_PATH"/gopls/internal/fuzzy/* internal/fuzzy/
cp -r "$TOOLS_PATH"/gopls/internal/util/astutil/* internal/util/astutil/
cp -r "$TOOLS_PATH"/gopls/internal/util/bug/* internal/util/bug/
cp -r "$TOOLS_PATH"/gopls/internal/util/moreiters/* internal/util/moreiters/
cp -r "$TOOLS_PATH"/gopls/internal/util/moreslices/* internal/util/moreslices/
cp -r "$TOOLS_PATH"/gopls/internal/util/safetoken/* internal/util/safetoken/
cp -r "$TOOLS_PATH"/gopls/internal/util/tokeninternal/* internal/util/tokeninternal/

# x/tools internal dependencies
mkdir -p internal/xtools/{aliases,analysisinternal/typeindex,astutil,stdlib,typeparams,typesinternal/typeindex,versions}

cp -r "$TOOLS_PATH"/internal/aliases/* internal/xtools/aliases/
cp -r "$TOOLS_PATH"/internal/analysisinternal/* internal/xtools/analysisinternal/
cp -r "$TOOLS_PATH"/internal/analysisinternal/typeindex/* internal/xtools/analysisinternal/typeindex/
cp -r "$TOOLS_PATH"/internal/astutil/* internal/xtools/astutil/
cp -r "$TOOLS_PATH"/internal/stdlib/* internal/xtools/stdlib/
cp -r "$TOOLS_PATH"/internal/typeparams/* internal/xtools/typeparams/
cp -r "$TOOLS_PATH"/internal/typesinternal/* internal/xtools/typesinternal/
cp -r "$TOOLS_PATH"/internal/typesinternal/typeindex/* internal/xtools/typesinternal/typeindex/
cp -r "$TOOLS_PATH"/internal/versions/* internal/xtools/versions/

# Remove test files from internal dependencies as well
find internal -name "*_test.go" -type f -delete
find internal -type d -name "testdata" -exec rm -rf {} + 2>/dev/null
```

Note: If new dependencies are added in future gopls versions, check for missing imports during build and add them. Test dependencies (testenv, testfiles) are not needed since we don't include tests.

### 4. Update Import Paths

Replace all gopls/x/tools internal import paths with plscheck paths:

```bash
# Update gopls internal imports
find . -name "*.go" -type f -exec sed -i '' 's|golang.org/x/tools/gopls/internal/|lesiw.io/plscheck/internal/|g' {} +

# Update analyzer cross-references
find . -name "*.go" -type f -exec sed -i '' 's|golang.org/x/tools/gopls/internal/analysis/|lesiw.io/plscheck/|g' {} +

# Update x/tools internal imports
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/analysisinternal"|"lesiw.io/plscheck/internal/xtools/analysisinternal"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/typesinternal"|"lesiw.io/plscheck/internal/xtools/typesinternal"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/typeparams"|"lesiw.io/plscheck/internal/xtools/typeparams"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/astutil"|"lesiw.io/plscheck/internal/xtools/astutil"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/aliases"|"lesiw.io/plscheck/internal/xtools/aliases"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/analysisinternal/typeindex"|"lesiw.io/plscheck/internal/xtools/analysisinternal/typeindex"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/typesinternal/typeindex"|"lesiw.io/plscheck/internal/xtools/typesinternal/typeindex"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/stdlib"|"lesiw.io/plscheck/internal/xtools/stdlib"|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|"golang.org/x/tools/internal/versions"|"lesiw.io/plscheck/internal/xtools/versions"|g' {} +

# Fix circular imports in internal/xtools
sed -i '' 's|"golang.org/x/tools/internal/typesinternal"|"lesiw.io/plscheck/internal/xtools/typesinternal"|g' internal/xtools/typesinternal/typeindex/typeindex.go
```

### 5. Update cmd/plscheck/main.go

Review the list of analyzers and update `cmd/plscheck/main.go`:

1. Check which analyzers export an `Analyzer` variable (some like fillstruct/fillswitch may not)
2. Update the import list
3. Update the `multichecker.Main()` call with the available analyzers

Example pattern:
```go
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"lesiw.io/plscheck/analyzer1"
	"lesiw.io/plscheck/analyzer2"
	// ... add all analyzers that export Analyzer
)

func main() {
	multichecker.Main(
		analyzer1.Analyzer,
		analyzer2.Analyzer,
		// ... list all analyzers
	)
}
```

### 6. Synchronize go.mod Versions

Critical: All versions must match gopls exactly.

```bash
# Check gopls go.mod for required versions
# The go mod download will cache it in your GOMODCACHE
go mod download -json golang.org/x/tools/gopls@vX.Y.Z | grep GoMod

# Or check directly from your local clone
cat "$TOOLS_PATH"/gopls/go.mod

# Update plscheck go.mod to match:
# 1. Go version (e.g., go 1.24.2)
# 2. golang.org/x/tools/gopls version (e.g., v0.20.0)
# 3. Add replace directive for golang.org/x/tools to match gopls's exact requirement
```

Example go.mod structure:
```go
module lesiw.io/plscheck

go 1.24.2  // Must match gopls

require (
	golang.org/x/tools/gopls vX.Y.Z
	// ... other deps from go mod tidy
)

// Pin to the exact x/tools version used by gopls
replace golang.org/x/tools => golang.org/x/tools v0.35.1-0.20250728180453-01a3475a31bc
```

### 7. Build and Test

```bash
go mod tidy
go build ./cmd/plscheck

# Verify it works
./plscheck -help
./plscheck ./cmd/plscheck
```

### 8. Handle Known Issues

- **fillstruct/fillswitch**: These don't export standalone Analyzers. Comment out their `SuggestedFix` functions if they use gopls-internal packages like `cache.Package` or `parsego.File`. Remove unused imports.

- **Missing dependencies**: If build fails with missing internal packages, check upstream for new dependencies and copy them following the same pattern.

### 9. Update Documentation

Update README.md if the analyzer list changed (additions/removals).

### 10. Git Commit and Tag

```bash
git add -A
git commit -m "Update to gopls vX.Y.Z

Synchronized with gopls/vX.Y.Z from golang.org/x/tools

<list any notable changes, new/removed analyzers, etc.>

Dependencies pinned to match gopls vX.Y.Z exactly."

git tag -a vX.Y.Z -m "plscheck vX.Y.Z

Matches gopls/vX.Y.Z from golang.org/x/tools

This release includes N analyzers externalized from gopls internal
analysis package, allowing them to be used in CI and other tooling."
```

## Version Synchronization Checklist

Before tagging, verify these are synchronized:

- [ ] Repository tag matches gopls version (e.g., v0.20.0)
- [ ] `go.mod` Go directive matches gopls (e.g., `go 1.24.2`)
- [ ] `go.mod` has `golang.org/x/tools/gopls` at target version
- [ ] `go.mod` has replace directive pinning `golang.org/x/tools` to gopls's exact requirement
- [ ] Build succeeds: `go build ./cmd/plscheck`
- [ ] CLI works: `./plscheck -help`
- [ ] All Go Authors copyright notices preserved
- [ ] No plscheck-specific copyright notices added

## Quick Reference: Key Files

- `cmd/plscheck/main.go` - Update analyzer list
- `go.mod` - Sync versions and add replace directive
- `README.md` - Update analyzer list if changed
- Individual analyzers - May need import fixes or function commenting

## Notes

- Always preserve Go Authors copyright notices in copied files
- Never add custom copyright notices to our wrapper code
- The replace directive in go.mod is essential for version consistency
- Test on actual Go code to ensure analyzers work correctly
