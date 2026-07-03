# FixtureForge - AGENTS.md

## Project Overview
FixtureForge is a Go CLI tool and library for generating realistic test data. It supports 28 column types, 7 output formats, auto-detection of column types from names, schema files, and deterministic generation with seeds.

## Architecture
- `cmd/fixtureforge/` - CLI entry point using cobra
- `pkg/fixtureforge/` - Public API (Engine)
- `internal/types/` - Core data types (Schema, Column, ColumnType)
- `internal/generators/` - 28 data generators with a Registry pattern
- `internal/detectors/` - Auto-detection of column types from names
- `internal/formatters/` - Output formatters (CSV, JSON, JSONL, YAML, Markdown, HTML, SQL)

## Key Files
- `internal/generators/generators.go` - All 28 generators (StringGen, IntGen, EmailGen, UUIDGen, etc.)
- `internal/detectors/detectors.go` - DetectColumnType() maps names to types
- `internal/formatters/formatters.go` - Formatter interface and implementations
- `pkg/fixtureforge/engine.go` - Engine ties everything together

## Build & Test
```bash
go build ./cmd/fixtureforge/
go test ./... -v
```

## Adding a New Column Type
1. Add a new `ColumnType` constant to `internal/types/types.go`
2. Create a struct implementing the `Generator` interface in `internal/generators/generators.go`
3. Register it in `NewRegistry()` in `internal/generators/generators.go`
4. Add detection rules in `internal/detectors/detectors.go`
5. Add tests in `internal/generators/generators_test.go`

## Generator Interface
```go
type Generator interface {
    Generate(rng *rand.Rand, params map[string]string) interface{}
}
```

## Commit Conventions
- `feat:` new feature
- `fix:` bug fix
- `docs:` documentation
- `test:` tests
- `refactor:` code restructuring
- `chore:` dependency updates
