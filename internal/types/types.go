package types

import (
	"fmt"
	"math/rand"
)

// Schema defines the structure of generated data
type Schema struct {
	Columns []Column `yaml:"columns" json:"columns"`
	Count   int      `yaml:"count" json:"count"`
	Seed    int64    `yaml:"seed,omitempty" json:"seed,omitempty"`
}

// Column defines a single column's generation rules
type Column struct {
	Name     string            `yaml:"name" json:"name"`
	Type     ColumnType        `yaml:"type" json:"type"`
	Nullable bool              `yaml:"nullable,omitempty" json:"nullable,omitempty"`
	NullPct  float64           `yaml:"null_pct,omitempty" json:"null_pct,omitempty"`
	Params   map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
	Unique   bool              `yaml:"unique,omitempty" json:"unique,omitempty"`
}

// ColumnType represents the type of data to generate
type ColumnType string

const (
	TypeString     ColumnType = "string"
	TypeInt        ColumnType = "int"
	TypeFloat      ColumnType = "float"
	TypeBool       ColumnType = "bool"
	TypeName       ColumnType = "name"
	TypeFirstName  ColumnType = "first_name"
	TypeLastName   ColumnType = "last_name"
	TypeEmail      ColumnType = "email"
	TypePhone      ColumnType = "phone"
	TypeAddress    ColumnType = "address"
	TypeCity       ColumnType = "city"
	TypeState      ColumnType = "state"
	TypeCountry    ColumnType = "country"
	TypeZipCode    ColumnType = "zipcode"
	TypeDate       ColumnType = "date"
	TypeDateTime   ColumnType = "datetime"
	TypeTime       ColumnType = "time"
	TypeUUID       ColumnType = "uuid"
	TypeURL        ColumnType = "url"
	TypeIPv4       ColumnType = "ipv4"
	TypeIPv6       ColumnType = "ipv6"
	TypeColor      ColumnType = "color"
	TypeLorem      ColumnType = "lorem"
	TypeEnum       ColumnType = "enum"
	TypeRegex      ColumnType = "regex"
	TypeRange      ColumnType = "range"
	TypeJSON       ColumnType = "json"
	TypeBoolWeight ColumnType = "bool_weight"
)

// Generator is the interface for all data generators
type Generator interface {
	Generate(rng *rand.Rand, params map[string]string) interface{}
}

// GeneratedRow represents a single row of generated data
type GeneratedRow map[string]interface{}

// String returns a formatted representation
func (c Column) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Type)
}
