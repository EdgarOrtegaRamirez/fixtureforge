package formatters

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

var testColumns = []types.Column{
	{Name: "id", Type: types.TypeInt},
	{Name: "name", Type: types.TypeString},
	{Name: "email", Type: types.TypeEmail},
}

var testRows = []types.GeneratedRow{
	{"id": int64(1), "name": "Alice Smith", "email": "alice@example.com"},
	{"id": int64(2), "name": "Bob Jones", "email": "bob@example.com"},
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("csv")
	if err != nil {
		t.Fatalf("New(csv) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Errorf("CSV output has %d lines, want 3", len(lines))
	}
	if lines[0] != "id,name,email" {
		t.Errorf("CSV header = %q, want %q", lines[0], "id,name,email")
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("json")
	if err != nil {
		t.Fatalf("New(json) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("JSON output has %d items, want 2", len(result))
	}
}

func TestJSONLFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("jsonl")
	if err != nil {
		t.Fatalf("New(jsonl) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("JSONL output has %d lines, want 2", len(lines))
	}

	for _, line := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Errorf("Invalid JSONL line: %v", err)
		}
	}
}

func TestYAMLFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("yaml")
	if err != nil {
		t.Fatalf("New(yaml) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id:") {
		t.Error("YAML output missing 'id:' key")
	}
	if !strings.Contains(output, "name:") {
		t.Error("YAML output missing 'name:' key")
	}
	if !strings.Contains(output, "email:") {
		t.Error("YAML output missing 'email:' key")
	}
}

func TestMarkdownFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("markdown")
	if err != nil {
		t.Fatalf("New(markdown) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "| id | name | email |") {
		t.Error("Markdown output missing header row")
	}
	if !strings.Contains(output, "| --- | --- | --- |") {
		t.Error("Markdown output missing separator row")
	}
}

func TestHTMLFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("html")
	if err != nil {
		t.Fatalf("New(html) error = %v", err)
	}

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range testRows {
		if err := formatter.WriteRow(&buf, row, testColumns); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<table>") {
		t.Error("HTML output missing <table> tag")
	}
	if !strings.Contains(output, "</table>") {
		t.Error("HTML output missing </table> tag")
	}
	if !strings.Contains(output, "<th>id</th>") {
		t.Error("HTML output missing <th>id</th>")
	}
}

func TestSQLFormatter(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := New("sql")
	if err != nil {
		t.Fatalf("New(sql) error = %v", err)
	}

	cols := []types.Column{
		{Name: "id", Type: types.TypeInt},
		{Name: "name", Type: types.TypeString},
		{Name: "active", Type: types.TypeBool},
	}
	rows := []types.GeneratedRow{
		{"id": int64(1), "name": "Alice", "active": true},
		{"id": int64(2), "name": "Bob", "active": false},
	}

	if err := formatter.WriteHeader(&buf, cols); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	for _, row := range rows {
		if err := formatter.WriteRow(&buf, row, cols); err != nil {
			t.Fatalf("WriteRow() error = %v", err)
		}
	}

	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "INSERT INTO") {
		t.Error("SQL output missing INSERT statement")
	}
	if !strings.Contains(output, "VALUES") {
		t.Error("SQL output missing VALUES clause")
	}
	if !strings.Contains(output, "TRUE") || !strings.Contains(output, "FALSE") {
		t.Error("SQL output missing boolean values")
	}
}

func TestSQLFormatterNulls(t *testing.T) {
	var buf bytes.Buffer
	formatter, _ := New("sql")

	if err := formatter.WriteHeader(&buf, testColumns); err != nil {
		t.Fatal(err)
	}
	if err := formatter.WriteRow(&buf, types.GeneratedRow{"id": int64(1), "name": nil, "email": "test@test.com"}, testColumns); err != nil {
		t.Fatal(err)
	}
	if err := formatter.WriteFooter(&buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "NULL") {
		t.Error("SQL output missing NULL for nil values")
	}
}

func TestNewFormatterInvalid(t *testing.T) {
	_, err := New("invalid")
	if err == nil {
		t.Error("New(invalid) should return error")
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSQLValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, "NULL"},
		{"bool true", true, "TRUE"},
		{"bool false", false, "FALSE"},
		{"string", "hello", "'hello'"},
		{"string with quotes", "it's", "'it''s'"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlValue(tt.input)
			if result != tt.expected {
				t.Errorf("sqlValue(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
