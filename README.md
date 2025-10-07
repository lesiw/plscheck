# plscheck

Externalized versions of gopls internal analyzers from [golang.org/x/tools/gopls/internal/analysis](https://github.com/golang/tools/tree/master/gopls/internal/analysis).

## Purpose

This package allows you to run gopls analyzers in CI pipelines and other tooling outside of gopls itself. The analyzers in gopls are normally only available within the language server, but this package makes them available as standalone checkers.

## Versioning

Tags in this repository always match upstream gopls releases. For example, `v0.20.0` in this repository corresponds to `gopls/v0.20.0` in golang.org/x/tools.

All dependencies (golang.org/x/tools and golang.org/x/tools/gopls) are pinned to the same version to ensure consistency.

## Available Analyzers

- **deprecated** - checks for usage of deprecated identifiers
- **embeddirective** - validates //go:embed directives
- **fillreturns** - suggests fixes for incomplete return statements
- **infertypeargs** - suggests removing unnecessary type arguments
- **maprange** - suggests using maps.Keys/Values/All for map iteration
- **modernize** - suggests modern Go idioms (e.g., min/max builtins)
- **nonewvars** - detects := assignments that don't introduce new variables
- **noresultvalues** - checks for return statements with unexpected result values
- **recursiveiter** - detects recursive iterator patterns that may cause issues
- **simplifycompositelit** - simplifies composite literal types
- **simplifyrange** - simplifies range statement expressions
- **simplifyslice** - simplifies slice expressions
- **unusedfunc** - finds unused functions
- **unusedparams** - finds unused function parameters
- **unusedvariable** - finds unused variables
- **yield** - checks yield function calls in iterator functions

## Usage

```bash
# Install
go install lesiw.io/plscheck/cmd/plscheck@latest

# Run all analyzers
plscheck ./...

# Run specific analyzers
plscheck -unusedparams -unusedvariable ./...
```

## Limitations

Some gopls analyzers (fillstruct, fillswitch) are tightly coupled with gopls internals and are not included as they don't export standalone analyzers suitable for command-line use.
