package fixtureforge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

func TestGenerateCSV(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeName},
			{Name: "email", Type: types.TypeEmail},
		},
		Count: 5,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	n, err := engine.Generate(schema, &buf, "csv")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if n != 5 {
		t.Errorf("Generate() returned %d rows, want 5", n)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 6 { // header + 5 rows
		t.Errorf("CSV output has %d lines, want 6", len(lines))
	}
	if lines[0] != "id,name,email" {
		t.Errorf("CSV header = %q, want %q", lines[0], "id,name,email")
	}
}

func TestGenerateJSON(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeString},
		},
		Count: 3,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "json")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("JSON output has %d items, want 3", len(result))
	}
}

func TestGenerateJSONL(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "email", Type: types.TypeEmail},
		},
		Count: 5,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "jsonl")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Errorf("JSONL output has %d lines, want 5", len(lines))
	}

	for _, line := range lines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Errorf("Invalid JSONL line: %v", err)
		}
	}
}

func TestGenerateYAML(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeString},
		},
		Count: 3,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "yaml")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id:") {
		t.Error("YAML output missing 'id:' key")
	}
	if !strings.Contains(output, "name:") {
		t.Error("YAML output missing 'name:' key")
	}
}

func TestGenerateSQL(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeString},
		},
		Count: 3,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "sql")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "INSERT INTO") {
		t.Error("SQL output missing INSERT statement")
	}
	if !strings.Contains(output, "VALUES") {
		t.Error("SQL output missing VALUES clause")
	}
}

func TestGenerateHTML(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeString},
		},
		Count: 2,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "html")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<table>") {
		t.Error("HTML output missing <table> tag")
	}
	if !strings.Contains(output, "<th>id</th>") {
		t.Error("HTML output missing <th>id</th>")
	}
}

func TestGenerateMarkdown(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "name", Type: types.TypeString},
		},
		Count: 2,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "markdown")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "| id | name |") {
		t.Error("Markdown output missing header")
	}
	if !strings.Contains(output, "| --- | --- |") {
		t.Error("Markdown output missing separator")
	}
}

func TestDeterministicOutput(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "name", Type: types.TypeName},
			{Name: "email", Type: types.TypeEmail},
		},
		Count: 5,
		Seed:  42,
	}

	var buf1, buf2 bytes.Buffer
	engine := NewEngine()

	_, _ = engine.Generate(schema, &buf1, "csv")
	_, _ = engine.Generate(schema, &buf2, "csv")

	if buf1.String() != buf2.String() {
		t.Error("Deterministic output should be identical with same seed")
	}
}

func TestDetectSchema(t *testing.T) {
	names := []string{"id", "first_name", "email", "created_at", "price", "active"}
	cols := DetectSchema(names)

	expected := []struct {
		name string
		typ  types.ColumnType
	}{
		{"id", types.TypeInt},
		{"first_name", types.TypeFirstName},
		{"email", types.TypeEmail},
		{"created_at", types.TypeDateTime},
		{"price", types.TypeFloat},
		{"active", types.TypeBool},
	}

	for i, exp := range expected {
		if cols[i].Name != exp.name {
			t.Errorf("column %d name = %q, want %q", i, cols[i].Name, exp.name)
		}
		if cols[i].Type != exp.typ {
			t.Errorf("column %d type = %q, want %q", i, cols[i].Type, exp.typ)
		}
	}
}

func TestValidateSchema(t *testing.T) {
	tests := []struct {
		name    string
		schema  *types.Schema
		wantErr bool
	}{
		{
			name: "valid",
			schema: &types.Schema{
				Columns: []types.Column{{Name: "id", Type: types.TypeInt}},
				Count:   10,
			},
			wantErr: false,
		},
		{
			name: "no columns",
			schema: &types.Schema{
				Columns: []types.Column{},
				Count:   10,
			},
			wantErr: true,
		},
		{
			name: "zero count",
			schema: &types.Schema{
				Columns: []types.Column{{Name: "id", Type: types.TypeInt}},
				Count:   0,
			},
			wantErr: true,
		},
		{
			name: "duplicate columns",
			schema: &types.Schema{
				Columns: []types.Column{
					{Name: "id", Type: types.TypeInt},
					{Name: "id", Type: types.TypeInt},
				},
				Count: 10,
			},
			wantErr: true,
		},
		{
			name: "empty column name",
			schema: &types.Schema{
				Columns: []types.Column{{Name: "", Type: types.TypeInt}},
				Count:   10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := ValidateSchema(tt.schema)
			if tt.wantErr && len(issues) == 0 {
				t.Error("ValidateSchema() expected errors, got none")
			}
			if !tt.wantErr && len(issues) > 0 {
				t.Errorf("ValidateSchema() unexpected errors: %v", issues)
			}
		})
	}
}

func TestAllColumnTypes(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{Name: "string_col", Type: types.TypeString},
			{Name: "int_col", Type: types.TypeInt},
			{Name: "float_col", Type: types.TypeFloat},
			{Name: "bool_col", Type: types.TypeBool},
			{Name: "name_col", Type: types.TypeName},
			{Name: "first_name_col", Type: types.TypeFirstName},
			{Name: "last_name_col", Type: types.TypeLastName},
			{Name: "email_col", Type: types.TypeEmail},
			{Name: "phone_col", Type: types.TypePhone},
			{Name: "address_col", Type: types.TypeAddress},
			{Name: "city_col", Type: types.TypeCity},
			{Name: "state_col", Type: types.TypeState},
			{Name: "country_col", Type: types.TypeCountry},
			{Name: "zipcode_col", Type: types.TypeZipCode},
			{Name: "date_col", Type: types.TypeDate},
			{Name: "datetime_col", Type: types.TypeDateTime},
			{Name: "time_col", Type: types.TypeTime},
			{Name: "uuid_col", Type: types.TypeUUID},
			{Name: "url_col", Type: types.TypeURL},
			{Name: "ipv4_col", Type: types.TypeIPv4},
			{Name: "ipv6_col", Type: types.TypeIPv6},
			{Name: "color_col", Type: types.TypeColor},
			{Name: "lorem_col", Type: types.TypeLorem},
			{Name: "enum_col", Type: types.TypeEnum},
			{Name: "range_col", Type: types.TypeRange},
			{Name: "json_col", Type: types.TypeJSON},
			{Name: "bool_weight_col", Type: types.TypeBoolWeight},
		},
		Count: 1,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	n, err := engine.Generate(schema, &buf, "json")
	if err != nil {
		t.Fatalf("Generate() with all types error = %v", err)
	}
	if n != 1 {
		t.Errorf("Generate() returned %d rows, want 1", n)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}
	if len(result) != 1 {
		t.Fatal("Expected 1 row")
	}

	// Verify each column has a non-nil value
	for _, col := range schema.Columns {
		if result[0][col.Name] == nil {
			t.Errorf("column %q is nil", col.Name)
		}
	}
}

func TestRegexGenerator(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{
				Name:   "code",
				Type:   types.TypeRegex,
				Params: map[string]string{"pattern": "[A-Z][a-z][a-z]"},
			},
		},
		Count: 5,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "csv")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		if len(parts[0]) != 3 {
			t.Errorf("regex output %q has length %d, want 3", parts[0], len(parts[0]))
		}
	}
}

func TestEnumGenerator(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{
				Name:   "status",
				Type:   types.TypeEnum,
				Params: map[string]string{"values": "active,inactive,pending"},
			},
		},
		Count: 10,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "csv")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	validValues := map[string]bool{"active": true, "inactive": true, "pending": true}
	for i := 1; i < len(lines); i++ {
		if !validValues[lines[i]] {
			t.Errorf("enum value %q not in expected set", lines[i])
		}
	}
}

func TestStringGenerator(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{
				Name:   "code",
				Type:   types.TypeString,
				Params: map[string]string{"min_length": "5", "max_length": "10"},
			},
		},
		Count: 10,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "csv")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		length := len(lines[i])
		if length < 5 || length > 10 {
			t.Errorf("string length %d not in range [5, 10]", length)
		}
	}
}

func TestFloatGenerator(t *testing.T) {
	schema := &types.Schema{
		Columns: []types.Column{
			{
				Name:   "price",
				Type:   types.TypeFloat,
				Params: map[string]string{"min": "10", "max": "100", "decimals": "2"},
			},
		},
		Count: 20,
		Seed:  42,
	}

	var buf bytes.Buffer
	engine := NewEngine()
	_, err := engine.Generate(schema, &buf, "csv")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		var val float64
		_, err := fmt.Sscanf(lines[i], "%f", &val)
		if err != nil {
			t.Errorf("invalid float value %q: %v", lines[i], err)
		}
	}
}
