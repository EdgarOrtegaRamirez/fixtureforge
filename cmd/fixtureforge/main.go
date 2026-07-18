package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
	"github.com/EdgarOrtegaRamirez/fixtureforge/pkg/fixtureforge"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var version = "0.1.0"

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "fixtureforge",
		Short: "Generate realistic test data for development and testing",
		Long:  `FixtureForge generates realistic test data in CSV, JSON, JSONL, YAML, Markdown, HTML, and SQL formats.`,
	}

	root.AddCommand(
		newGenerateCmd(),
		newDetectCmd(),
		newValidateCmd(),
		newSchemaCmd(),
		newVersionCmd(),
	)

	return root
}

func newGenerateCmd() *cobra.Command {
	var (
		count      int
		format     string
		seed       int64
		columns    []string
		schemaFile string
	)

	cmd := &cobra.Command{
		Use:   "generate [flags]",
		Short: "Generate test data",
		Long:  `Generate test data based on column definitions or auto-detection from column names.`,
		Example: `  # Generate 10 rows of CSV with auto-detected columns
  fixtureforge generate -n 10 -f csv --columns name,email,age

  # Generate JSON with a schema file
  fixtureforge generate -n 100 -f json -s schema.yaml

  # Generate with a seed for reproducible output
  fixtureforge generate -n 5 -f csv --columns id,name,email --seed 12345`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var schema *types.Schema

			if schemaFile != "" {
				// Load from schema file
				data, err := os.ReadFile(schemaFile)
				if err != nil {
					return fmt.Errorf("reading schema file: %w", err)
				}
				schema = &types.Schema{}
				if strings.HasSuffix(schemaFile, ".json") {
					if err := json.Unmarshal(data, schema); err != nil {
						return fmt.Errorf("parsing JSON schema: %w", err)
					}
				} else {
					if err := yaml.Unmarshal(data, schema); err != nil {
						return fmt.Errorf("parsing YAML schema: %w", err)
					}
				}
			} else if len(columns) > 0 {
				// Build schema from column names
				detectedCols := fixtureforge.DetectSchema(columns)
				schema = &types.Schema{
					Columns: detectedCols,
					Count:   count,
					Seed:    seed,
				}
			} else {
				return fmt.Errorf("specify --columns or --schema")
			}

			if count > 0 && schemaFile != "" {
				schema.Count = count
			}
			if seed != 0 {
				schema.Seed = seed
			}

			// Validate
			if issues := fixtureforge.ValidateSchema(schema); len(issues) > 0 {
				return fmt.Errorf("schema validation failed:\n  - %s", strings.Join(issues, "\n  - "))
			}

			// Generate
			engine := fixtureforge.NewEngine()
			n, err := engine.Generate(schema, os.Stdout, format)
			if err != nil {
				return fmt.Errorf("generating data: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Generated %d rows\n", n)
			return nil
		},
	}

	cmd.Flags().IntVarP(&count, "count", "n", 10, "Number of rows to generate")
	cmd.Flags().StringVarP(&format, "format", "f", "csv", "Output format (csv, json, jsonl, yaml, markdown, html, sql)")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for reproducible output")
	cmd.Flags().StringSliceVar(&columns, "columns", nil, "Comma-separated column names (auto-detects types)")
	cmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "Path to schema file (YAML or JSON)")

	return cmd
}

func newDetectCmd() *cobra.Command {
	var columns []string

	cmd := &cobra.Command{
		Use:   "detect [flags]",
		Short: "Detect column types from names",
		Long:  `Auto-detect column types based on column names. Useful for seeing what types would be generated.`,
		Example: `  # Detect types from column names
  fixtureforge detect --columns name,email,age,created_at,price

  # Detect from a CSV header
  head -1 data.csv | tr ',' '\n' | xargs fixtureforge detect --columns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(columns) == 0 {
				return fmt.Errorf("specify --columns")
			}

			detected := fixtureforge.DetectSchema(columns)
			for _, col := range detected {
				fmt.Printf("%-20s → %s\n", col.Name, col.Type)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&columns, "columns", nil, "Column names to analyze")

	return cmd
}

func newValidateCmd() *cobra.Command {
	var schemaFile string

	cmd := &cobra.Command{
		Use:   "validate [flags]",
		Short: "Validate a schema file",
		Long:  `Check a schema file for common issues and report them.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if schemaFile == "" {
				return fmt.Errorf("specify --schema")
			}

			data, err := os.ReadFile(schemaFile)
			if err != nil {
				return fmt.Errorf("reading schema file: %w", err)
			}

			schema := &types.Schema{}
			if strings.HasSuffix(schemaFile, ".json") {
				if err := json.Unmarshal(data, schema); err != nil {
					return fmt.Errorf("parsing JSON schema: %w", err)
				}
			} else {
				if err := yaml.Unmarshal(data, schema); err != nil {
					return fmt.Errorf("parsing YAML schema: %w", err)
				}
			}

			issues := fixtureforge.ValidateSchema(schema)
			if len(issues) == 0 {
				fmt.Println("✓ Schema is valid")
			} else {
				fmt.Printf("✗ Found %d issue(s):\n", len(issues))
				for _, issue := range issues {
					fmt.Printf("  - %s\n", issue)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "Path to schema file (YAML or JSON)")

	return cmd
}

func newSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Schema management commands",
	}

	cmd.AddCommand(newSchemaSampleCmd())

	return cmd
}

func newSchemaSampleCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "sample [flags]",
		Short: "Generate a sample schema file",
		Long:  `Generate a sample schema file that you can customize.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sample := `# FixtureForge Schema
columns:
  - name: id
    type: int
    params:
      min: "1"
      max: "1000"
  - name: name
    type: name
  - name: email
    type: email
  - name: age
    type: range
    params:
      min: "18"
      max: "80"
  - name: active
    type: bool
  - name: created_at
    type: datetime
count: 10
seed: 42
`
			if format == "json" {
				sampleJSON := map[string]interface{}{
					"columns": []map[string]interface{}{
						{"name": "id", "type": "int", "params": map[string]string{"min": "1", "max": "1000"}},
						{"name": "name", "type": "name"},
						{"name": "email", "type": "email"},
						{"name": "age", "type": "range", "params": map[string]string{"min": "18", "max": "80"}},
						{"name": "active", "type": "bool"},
						{"name": "created_at", "type": "datetime"},
					},
					"count": 10,
					"seed":  42,
				}
				data, _ := json.MarshalIndent(sampleJSON, "", "  ")
				sample = string(data) + "\n"
			}

			_, err := io.WriteString(os.Stdout, sample)
			return err
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (yaml, json)")

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("fixtureforge v%s\n", version)
		},
	}
}
