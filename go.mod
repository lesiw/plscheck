module lesiw.io/plscheck

go 1.24.2

require (
	github.com/google/go-cmp v0.7.0
	golang.org/x/telemetry v0.0.0-20251001141935-4eae98a72453
	golang.org/x/tools v0.36.0
	golang.org/x/tools/gopls v0.20.0
)

require (
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

// Pin to the exact x/tools version used by gopls v0.20.0
replace golang.org/x/tools => golang.org/x/tools v0.35.1-0.20250728180453-01a3475a31bc
