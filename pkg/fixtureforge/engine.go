package fixtureforge

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/detectors"
	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/formatters"
	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/generators"
	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

// Engine generates test data based on a schema
type Engine struct {
	registry *generators.Registry
}

// NewEngine creates a new generation engine
func NewEngine() *Engine {
	return &Engine{
		registry: generators.NewRegistry(),
	}
}

// Generate produces rows of test data based on the schema and writes to output
func (e *Engine) Generate(schema *types.Schema, output io.Writer, format string) (int, error) {
	// Auto-detect column types if not specified
	for i := range schema.Columns {
		if schema.Columns[i].Type == "" {
			schema.Columns[i].Type = detectors.DetectColumnType(schema.Columns[i].Name)
		}
	}

	// Create RNG with seed
	var rng *rand.Rand
	if schema.Seed != 0 {
		rng = rand.New(rand.NewSource(schema.Seed))
	} else {
		rng = rand.New(rand.NewSource(42))
	}

	// Create formatter
	formatter, err := formatters.New(format)
	if err != nil {
		return 0, fmt.Errorf("creating formatter: %w", err)
	}

	// Write header
	if err := formatter.WriteHeader(output, schema.Columns); err != nil {
		return 0, fmt.Errorf("writing header: %w", err)
	}

	// Generate rows
	count := 0
	for i := 0; i < schema.Count; i++ {
		row := e.generateRow(rng, schema.Columns)
		if err := formatter.WriteRow(output, row, schema.Columns); err != nil {
			return 0, fmt.Errorf("writing row %d: %w", i, err)
		}
		count++
	}

	// Write footer
	if err := formatter.WriteFooter(output); err != nil {
		return 0, fmt.Errorf("writing footer: %w", err)
	}

	return count, nil
}

// GenerateRow generates a single row of test data
func (e *Engine) GenerateRow(schema *types.Schema) types.GeneratedRow {
	var rng *rand.Rand
	if schema.Seed != 0 {
		rng = rand.New(rand.NewSource(schema.Seed))
	} else {
		rng = rand.New(rand.NewSource(42))
	}
	return e.generateRow(rng, schema.Columns)
}

func (e *Engine) generateRow(rng *rand.Rand, columns []types.Column) types.GeneratedRow {
	row := make(types.GeneratedRow)
	for _, col := range columns {
		// Handle nulls
		if col.Nullable && rng.Float64()*100 < col.NullPct {
			row[col.Name] = nil
			continue
		}

		gen, ok := e.registry.Get(col.Type)
		if !ok {
			// Fallback to string
			gen, _ = e.registry.Get(types.TypeString)
		}
		row[col.Name] = gen.Generate(rng, col.Params)
	}
	return row
}

// DetectSchema auto-detects column types from names
func DetectSchema(names []string) []types.Column {
	columns := make([]types.Column, len(names))
	for i, name := range names {
		columns[i] = types.Column{
			Name: name,
			Type: detectors.DetectColumnType(name),
		}
	}
	return columns
}

// ValidateSchema checks a schema for issues
func ValidateSchema(schema *types.Schema) []string {
	var issues []string
	if len(schema.Columns) == 0 {
		issues = append(issues, "schema has no columns")
	}
	if schema.Count <= 0 {
		issues = append(issues, "count must be positive")
	}
	if schema.Count > 10000000 {
		issues = append(issues, "count exceeds maximum (10 million)")
	}
	seen := make(map[string]bool)
	for _, col := range schema.Columns {
		if seen[col.Name] {
			issues = append(issues, fmt.Sprintf("duplicate column name: %s", col.Name))
		}
		seen[col.Name] = true
		if col.Name == "" {
			issues = append(issues, "column name cannot be empty")
		}
	}
	return issues
}
