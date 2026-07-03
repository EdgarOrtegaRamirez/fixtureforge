package detectors

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

func TestDetectColumnType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.ColumnType
	}{
		{"id", "id", types.TypeInt},
		{"user_id", "user_id", types.TypeInt},
		{"email", "email", types.TypeEmail},
		{"email_address", "email_address", types.TypeEmail},
		{"phone", "phone", types.TypePhone},
		{"telephone", "telephone", types.TypePhone},
		{"name", "name", types.TypeName},
		{"full_name", "full_name", types.TypeName},
		{"first_name", "first_name", types.TypeFirstName},
		{"fname", "fname", types.TypeFirstName},
		{"last_name", "last_name", types.TypeLastName},
		{"surname", "surname", types.TypeLastName},
		{"address", "address", types.TypeAddress},
		{"street_address", "street_address", types.TypeAddress},
		{"city", "city", types.TypeCity},
		{"town", "town", types.TypeCity},
		{"state", "state", types.TypeState},
		{"province", "province", types.TypeState},
		{"country", "country", types.TypeCountry},
		{"zip_code", "zip_code", types.TypeZipCode},
		{"postal_code", "postal_code", types.TypeZipCode},
		{"date", "date", types.TypeDate},
		{"created_date", "created_date", types.TypeDate},
		{"datetime", "datetime", types.TypeDateTime},
		{"timestamp", "timestamp", types.TypeDateTime},
		{"created_at", "created_at", types.TypeDateTime},
		{"time", "time", types.TypeTime},
		{"uuid", "uuid", types.TypeUUID},
		{"guid", "guid", types.TypeUUID},
		{"url", "url", types.TypeURL},
		{"website", "website", types.TypeURL},
		{"link", "link", types.TypeURL},
		{"ip", "ip", types.TypeIPv4},
		{"ip_address", "ip_address", types.TypeIPv4},
		{"ipv4", "ipv4", types.TypeIPv4},
		{"ipv6", "ipv6", types.TypeIPv6},
		{"color", "color", types.TypeColor},
		{"hex_color", "hex_color", types.TypeColor},
		{"description", "description", types.TypeLorem},
		{"bio", "bio", types.TypeLorem},
		{"active", "active", types.TypeBool},
		{"enabled", "enabled", types.TypeBool},
		{"is_active", "is_active", types.TypeBool},
		{"price", "price", types.TypeFloat},
		{"amount", "amount", types.TypeFloat},
		{"salary", "salary", types.TypeFloat},
		{"quantity", "quantity", types.TypeInt},
		{"count", "count", types.TypeInt},
		{"age", "age", types.TypeRange},
		{"status", "status", types.TypeEnum},
		{"role", "role", types.TypeEnum},
		{"random_column", "random_column", types.TypeString},
		{"foo", "foo", types.TypeString},
		{"UPPER_CASE", "UPPER_CASE", types.TypeString},
		{"MyColumn", "MyColumn", types.TypeString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectColumnType(tt.input)
			if result != tt.expected {
				t.Errorf("DetectColumnType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetectSchema(t *testing.T) {
	names := []string{"id", "first_name", "email", "created_at", "active"}
	cols := make([]types.Column, len(names))
	for i, name := range names {
		cols[i] = types.Column{
			Name: name,
			Type: DetectColumnType(name),
		}
	}

	if len(cols) != 5 {
		t.Fatalf("DetectSchema() returned %d columns, want 5", len(cols))
	}

	expected := []struct {
		name string
		typ  types.ColumnType
	}{
		{"id", types.TypeInt},
		{"first_name", types.TypeFirstName},
		{"email", types.TypeEmail},
		{"created_at", types.TypeDateTime},
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
