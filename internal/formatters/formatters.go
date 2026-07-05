package formatters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

// Formatter outputs generated data in various formats
type Formatter interface {
	WriteHeader(w io.Writer, columns []types.Column) error
	WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error
	WriteFooter(w io.Writer) error
}

// New creates a formatter for the given format
func New(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case "csv":
		return &CSVFormatter{}, nil
	case "json":
		return &JSONFormatter{array: true}, nil
	case "jsonl", "ndjson":
		return &JSONFormatter{array: false}, nil
	case "yaml", "yml":
		return &YAMLFormatter{}, nil
	case "markdown", "md":
		return &MarkdownFormatter{}, nil
	case "html":
		return &HTMLFormatter{}, nil
	case "sql":
		return &SQLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// --- CSV Formatter ---

type CSVFormatter struct {
	writer *csv.Writer
}

func (f *CSVFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	f.writer = csv.NewWriter(w)
	header := make([]string, len(columns))
	for i, col := range columns {
		header[i] = col.Name
	}
	return f.writer.Write(header)
}

func (f *CSVFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	record := make([]string, len(columns))
	for i, col := range columns {
		record[i] = formatValue(row[col.Name])
	}
	return f.writer.Write(record)
}

func (f *CSVFormatter) WriteFooter(w io.Writer) error {
	f.writer.Flush()
	return f.writer.Error()
}

// --- JSON Formatter ---

type JSONFormatter struct {
	array      bool
	columns    []types.Column
	rows       []types.GeneratedRow
}

func (f *JSONFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	f.columns = columns
	return nil
}

func (f *JSONFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	f.rows = append(f.rows, row)
	return nil
}

func (f *JSONFormatter) WriteFooter(w io.Writer) error {
	if f.array {
		_, err := fmt.Fprint(w, "[\n")
		if err != nil {
			return err
		}
		for i, row := range f.rows {
			data := make(map[string]interface{})
			for _, col := range f.columns {
				data[col.Name] = row[col.Name]
			}
			encoded, err := json.Marshal(data)
			if err != nil {
				return err
			}
			suffix := ","
			if i == len(f.rows)-1 {
				suffix = ""
			}
			_, err = fmt.Fprintf(w, "  %s%s\n", string(encoded), suffix)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(w, "]\n")
		return err
	}

	for _, row := range f.rows {
		data := make(map[string]interface{})
		for _, col := range f.columns {
			data[col.Name] = row[col.Name]
		}
		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s\n", string(encoded))
		if err != nil {
			return err
		}
	}
	return nil
}

// --- YAML Formatter ---

type YAMLFormatter struct {
	columns []types.Column
	rows    []types.GeneratedRow
}

func (f *YAMLFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	f.columns = columns
	return nil
}

func (f *YAMLFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	f.rows = append(f.rows, row)
	return nil
}

func (f *YAMLFormatter) WriteFooter(w io.Writer) error {
	for _, row := range f.rows {
		first := true
		for _, col := range f.columns {
			if first {
				if _, err := fmt.Fprintf(w, "- %s: %v\n", col.Name, formatValue(row[col.Name])); err != nil {
					return err
				}
				first = false
			} else {
				if _, err := fmt.Fprintf(w, "  %s: %v\n", col.Name, formatValue(row[col.Name])); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// --- Markdown Formatter ---

type MarkdownFormatter struct {
	rows []types.GeneratedRow
}

func (f *MarkdownFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	// Header row
	parts := make([]string, len(columns))
	for i, col := range columns {
		parts[i] = col.Name
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(parts, " | ")); err != nil {
		return err
	}

	// Separator
	seps := make([]string, len(columns))
	for i := range seps {
		seps[i] = "---"
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | ")); err != nil {
		return err
	}
	return nil
}

func (f *MarkdownFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	f.rows = append(f.rows, row)
	return nil
}

func (f *MarkdownFormatter) WriteFooter(w io.Writer) error {
	if len(f.rows) == 0 {
		return nil
	}
	// Get consistent column order from first row
	cols := make([]string, 0, len(f.rows[0]))
	for k := range f.rows[0] {
		cols = append(cols, k)
	}
	for _, row := range f.rows {
		parts := make([]string, len(cols))
		for i, col := range cols {
			parts[i] = formatValue(row[col])
		}
		if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(parts, " | ")); err != nil {
			return err
		}
	}
	return nil
}

// --- HTML Formatter ---

type HTMLFormatter struct {
	rows   []types.GeneratedRow
	columns []types.Column
}

func (f *HTMLFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	f.columns = columns
	if _, err := fmt.Fprint(w, "<table>\n<thead>\n<tr>\n"); err != nil {
		return err
	}
	for _, col := range columns {
		if _, err := fmt.Fprintf(w, "  <th>%s</th>\n", col.Name); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, "</tr>\n</thead>\n<tbody>\n"); err != nil {
		return err
	}
	return nil
}

func (f *HTMLFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	f.rows = append(f.rows, row)
	return nil
}

func (f *HTMLFormatter) WriteFooter(w io.Writer) error {
	for _, row := range f.rows {
		if _, err := fmt.Fprint(w, "<tr>\n"); err != nil {
			return err
		}
		for _, col := range f.columns {
			if _, err := fmt.Fprintf(w, "  <td>%s</td>\n", formatValue(row[col.Name])); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, "</tr>\n"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, "</tbody>\n</table>\n"); err != nil {
		return err
	}
	return nil
}

// --- SQL Formatter ---

type SQLFormatter struct {
	tableName string
	rows      []types.GeneratedRow
}

func (f *SQLFormatter) WriteHeader(w io.Writer, columns []types.Column) error {
	f.tableName = "generated_data"
	return nil
}

func (f *SQLFormatter) WriteRow(w io.Writer, row types.GeneratedRow, columns []types.Column) error {
	f.rows = append(f.rows, row)
	return nil
}

func (f *SQLFormatter) WriteFooter(w io.Writer) error {
	if len(f.rows) == 0 {
		return nil
	}

	// Get column names from first row
	colNames := make([]string, 0, len(f.rows[0]))
	for k := range f.rows[0] {
		colNames = append(colNames, k)
	}

	// INSERT statement
	fmt.Fprintf(w, "INSERT INTO %s (%s) VALUES\n", f.tableName, strings.Join(colNames, ", "))

	for i, row := range f.rows {
		vals := make([]string, len(colNames))
		for j, col := range colNames {
			vals[j] = sqlValue(row[col])
		}
		prefix := "  "
		suffix := ","
		if i == len(f.rows)-1 {
			suffix = ";"
		}
		fmt.Fprintf(w, "%s(%s)%s\n", prefix, strings.Join(vals, ", "), suffix)
	}
	return nil
}

func sqlValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// --- Helpers ---

func formatValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case bool:
		if v {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		encoded, _ := json.Marshal(v)
		return string(encoded)
	default:
		return fmt.Sprintf("%v", v)
	}
}
