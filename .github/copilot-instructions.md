# Copilot Instructions for tf-summarize

## Project Overview
`tf-summarize` is a Go CLI tool that summarizes Terraform plan output into human-readable formats: table, tree, JSON, HTML, and separate-tree views. It reads plans from stdin or file (JSON or binary tfplan), parses them using `hashicorp/terraform-json`, and writes formatted summaries.

## Architecture
- **main.go** — CLI entry point with flag parsing, orchestrates reader → parser → writer pipeline
- **reader/** — `Reader` interface (`Read() ([]byte, error)`, `Name() string`) with `FileReader` and `StdinReader`
- **parser/** — `Parser` interface (`Parse() (tfjson.Plan, error)`) with `JSONParser` and `BinaryParser`
- **writer/** — `Writer` interface (`Write(io.Writer) error`) with `TableWriter`, `TreeWriter`, `SeparateTree`, `JSONWriter`, `HTMLWriter`, `JSONSumWriter`
- **terraformstate/** — Utility functions for categorizing resource changes by action type (add/delete/update/recreate/move/import)
- **tree/** — Tree data structure for hierarchical resource representation

## Coding Conventions
- Go 1.26, no CGO
- Use `testify/assert` for test assertions
- Tests are co-located in `*_test.go` files within each package
- Use `bytes.Buffer` for capturing writer output in tests
- Interfaces are minimal (1-2 methods)
- Error wrapping uses `fmt.Errorf` (not `errors.Wrap`)
- File permissions use `0600` for created files
- `terraformstate.ResourceChanges` is a type alias for `[]*tfjson.ResourceChange`

## Build & Test
```bash
make build    # goreleaser snapshot build + test
make test     # lint + go test -v ./... -count=1
make lint     # golangci-lint (goimports + gofmt)
make gosec    # security scan (excludes G204, G705)
```

## Key Patterns
- Writer constructors: `NewXxxWriter(changes, ...) Writer` — always return the `Writer` interface
- `CreateWriter()` in writer.go is the factory that selects the right writer based on flags
- `CreateReader()` and `CreateParser()` are similar factories for their respective interfaces
- HTML templates use Go's `text/template` with embedded FS (`embed.FS`)
- Resource changes are grouped by action type: `map[string]terraformstate.ResourceChanges`

## When Writing Tests
- Construct `map[string]terraformstate.ResourceChanges` with mock `tfjson.ResourceChange` entries
- Use `bytes.Buffer` to capture output, then `assert.Equal` on the string
- For tree-based writers, mock `NewTreeWriterFunc` to isolate from drawing logic
- Integration tests in `main_test.go` compile the binary and run it as a subprocess
