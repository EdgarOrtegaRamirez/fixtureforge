package generators

import (
	"math/rand"
	"testing"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

func TestStringGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &StringGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("StringGen.Generate() returned %T, want string", val)
		}
		if len(s) < 3 || len(s) > 20 {
			t.Errorf("string length %d not in range [3, 20]", len(s))
		}
	}
}

func TestStringGenCustomLength(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &StringGen{}
	params := map[string]string{"min_length": "5", "max_length": "10"}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, params).(string)
		if len(val) < 5 || len(val) > 10 {
			t.Errorf("string length %d not in range [5, 10]", len(val))
		}
	}
}

func TestStringGenWithPrefix(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &StringGen{}
	params := map[string]string{"prefix": "user_"}

	for i := 0; i < 10; i++ {
		val := gen.Generate(rng, params).(string)
		if len(val) < 5 || val[:5] != "user_" {
			t.Errorf("string %q doesn't have prefix 'user_'", val)
		}
	}
}

func TestIntGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &IntGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		n, ok := val.(int64)
		if !ok {
			t.Fatalf("IntGen.Generate() returned %T, want int64", val)
		}
		if n < 0 || n > 1000 {
			t.Errorf("int value %d not in range [0, 1000]", n)
		}
	}
}

func TestIntGenCustomRange(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &IntGen{}
	params := map[string]string{"min": "10", "max": "20"}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, params).(int64)
		if val < 10 || val > 20 {
			t.Errorf("int value %d not in range [10, 20]", val)
		}
	}
}

func TestFloatGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &FloatGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		f, ok := val.(float64)
		if !ok {
			t.Fatalf("FloatGen.Generate() returned %T, want float64", val)
		}
		if f < 0 || f > 1000 {
			t.Errorf("float value %f not in range [0, 1000]", f)
		}
	}
}

func TestBoolGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &BoolGen{}

	trueCount := 0
	falseCount := 0
	for i := 0; i < 1000; i++ {
		val := gen.Generate(rng, nil)
		b, ok := val.(bool)
		if !ok {
			t.Fatalf("BoolGen.Generate() returned %T, want bool", val)
		}
		if b {
			trueCount++
		} else {
			falseCount++
		}
	}

	// Should be roughly 50/50
	if trueCount < 400 || trueCount > 600 {
		t.Errorf("bool distribution %d/%d too skewed", trueCount, falseCount)
	}
}

func TestNameGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &NameGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("NameGen.Generate() returned %T, want string", val)
		}
		parts := splitSpace(s)
		if len(parts) != 2 {
			t.Errorf("name %q doesn't have exactly 2 parts", s)
		}
	}
}

func TestFirstNameGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &FirstNameGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("FirstNameGen.Generate() returned %T, want string", val)
		}
		if len(s) == 0 {
			t.Error("first name is empty")
		}
	}
}

func TestLastNameGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &LastNameGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("LastNameGen.Generate() returned %T, want string", val)
		}
		if len(s) == 0 {
			t.Error("last name is empty")
		}
	}
}

func TestEmailGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &EmailGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("EmailGen.Generate() returned %T, want string", val)
		}
		if !contains(s, "@") {
			t.Errorf("email %q doesn't contain @", s)
		}
		if !contains(s, ".") {
			t.Errorf("email %q doesn't contain .", s)
		}
	}
}

func TestPhoneGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &PhoneGen{}

	val := gen.Generate(rng, nil)
	s, ok := val.(string)
	if !ok {
		t.Fatalf("PhoneGen.Generate() returned %T, want string", val)
	}
	if len(s) < 10 {
		t.Errorf("phone number %q too short", s)
	}
}

func TestAddressGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &AddressGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("AddressGen.Generate() returned %T, want string", val)
		}
		if len(s) == 0 {
			t.Error("address is empty")
		}
	}
}

func TestCityGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &CityGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("CityGen.Generate() returned %T, want string", val)
		}
		if len(s) == 0 {
			t.Error("city is empty")
		}
	}
}

func TestStateGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &StateGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("StateGen.Generate() returned %T, want string", val)
		}
		if len(s) != 2 {
			t.Errorf("state code %q has length %d, want 2", s, len(s))
		}
	}
}

func TestCountryGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &CountryGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("CountryGen.Generate() returned %T, want string", val)
		}
		if len(s) == 0 {
			t.Error("country is empty")
		}
	}
}

func TestZipCodeGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &ZipCodeGen{}

	val := gen.Generate(rng, nil)
	s, ok := val.(string)
	if !ok {
		t.Fatalf("ZipCodeGen.Generate() returned %T, want string", val)
	}
	if len(s) != 5 {
		t.Errorf("zip code %q has length %d, want 5", s, len(s))
	}
}

func TestDateGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &DateGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("DateGen.Generate() returned %T, want string", val)
		}
		// ISO format: YYYY-MM-DD
		if len(s) != 10 || s[4] != '-' || s[7] != '-' {
			t.Errorf("date %q not in ISO format", s)
		}
	}
}

func TestDateTimeGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &DateTimeGen{}

	val := gen.Generate(rng, nil)
	s, ok := val.(string)
	if !ok {
		t.Fatalf("DateTimeGen.Generate() returned %T, want string", val)
	}
	if len(s) < 19 {
		t.Errorf("datetime %q too short", s)
	}
}

func TestTimeGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &TimeGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("TimeGen.Generate() returned %T, want string", val)
		}
		// Format: HH:MM:SS
		if len(s) != 8 || s[2] != ':' || s[5] != ':' {
			t.Errorf("time %q not in expected format", s)
		}
	}
}

func TestUUIDGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &UUIDGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("UUIDGen.Generate() returned %T, want string", val)
		}
		// UUID format: 8-4-4-4-12
		if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			t.Errorf("UUID %q not in expected format", s)
		}
		// Check version nibble
		if s[14] != '4' {
			t.Errorf("UUID version = %c, want 4", s[14])
		}
	}
}

func TestURLGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &URLGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("URLGen.Generate() returned %T, want string", val)
		}
		if !contains(s, "://") {
			t.Errorf("URL %q doesn't contain ://", s)
		}
	}
}

func TestIPv4Gen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &IPv4Gen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("IPv4Gen.Generate() returned %T, want string", val)
		}
		parts := splitDot(s)
		if len(parts) != 4 {
			t.Errorf("IPv4 %q doesn't have 4 parts", s)
		}
	}
}

func TestIPv6Gen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &IPv6Gen{}

	val := gen.Generate(rng, nil)
	s, ok := val.(string)
	if !ok {
		t.Fatalf("IPv6Gen.Generate() returned %T, want string", val)
	}
	parts := splitColon(s)
	if len(parts) != 8 {
		t.Errorf("IPv6 %q doesn't have 8 parts", s)
	}
}

func TestColorGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &ColorGen{}

	val := gen.Generate(rng, nil)
	s, ok := val.(string)
	if !ok {
		t.Fatalf("ColorGen.Generate() returned %T, want string", val)
	}
	if len(s) != 7 || s[0] != '#' {
		t.Errorf("color %q not in hex format", s)
	}
}

func TestLoremGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &LoremGen{}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, nil)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("LoremGen.Generate() returned %T, want string", val)
		}
		words := splitSpace(s)
		if len(words) != 10 {
			t.Errorf("lorem has %d words, want 10", len(words))
		}
	}
}

func TestEnumGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &EnumGen{}
	params := map[string]string{"values": "a,b,c"}

	validValues := map[string]bool{"a": true, "b": true, "c": true}
	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, params)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("EnumGen.Generate() returned %T, want string", val)
		}
		if !validValues[s] {
			t.Errorf("enum value %q not in expected set", s)
		}
	}
}

func TestRegexGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &RegexGen{}
	params := map[string]string{"pattern": "[a-z][a-z][a-z]"}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, params)
		s, ok := val.(string)
		if !ok {
			t.Fatalf("RegexGen.Generate() returned %T, want string", val)
		}
		if len(s) != 3 {
			t.Errorf("regex output %q has length %d, want 3", s, len(s))
		}
	}
}

func TestRangeGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &RangeGen{}
	params := map[string]string{"min": "10", "max": "20"}

	for i := 0; i < 100; i++ {
		val := gen.Generate(rng, params)
		n, ok := val.(int64)
		if !ok {
			t.Fatalf("RangeGen.Generate() returned %T, want int64", val)
		}
		if n < 10 || n > 20 {
			t.Errorf("range value %d not in [10, 20]", n)
		}
	}
}

func TestJSONGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &JSONGen{}

	val := gen.Generate(rng, nil)
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Fatalf("JSONGen.Generate() returned %T, want map[string]interface{}", val)
	}
	if _, exists := obj["id"]; !exists {
		t.Error("JSON object missing 'id' field")
	}
	if _, exists := obj["name"]; !exists {
		t.Error("JSON object missing 'name' field")
	}
}

func TestBoolWeightGen(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := &BoolWeightGen{}
	params := map[string]string{"true_pct": "90"}

	trueCount := 0
	for i := 0; i < 1000; i++ {
		val := gen.Generate(rng, params)
		b, ok := val.(bool)
		if !ok {
			t.Fatalf("BoolWeightGen.Generate() returned %T, want bool", val)
		}
		if b {
			trueCount++
		}
	}

	// With 90% true probability, should be roughly 900/1000
	if trueCount < 800 || trueCount > 1000 {
		t.Errorf("bool weight distribution %d/1000 too far from expected 900", trueCount)
	}
}

func TestRegistryDefaults(t *testing.T) {
	registry := NewRegistry()

	expectedTypes := []types.ColumnType{
		types.TypeString, types.TypeInt, types.TypeFloat, types.TypeBool,
		types.TypeName, types.TypeFirstName, types.TypeLastName,
		types.TypeEmail, types.TypePhone, types.TypeAddress,
		types.TypeCity, types.TypeState, types.TypeCountry, types.TypeZipCode,
		types.TypeDate, types.TypeDateTime, types.TypeTime,
		types.TypeUUID, types.TypeURL, types.TypeIPv4, types.TypeIPv6,
		types.TypeColor, types.TypeLorem, types.TypeEnum, types.TypeRegex,
		types.TypeRange, types.TypeJSON, types.TypeBoolWeight,
	}

	for _, ct := range expectedTypes {
		gen, ok := registry.Get(ct)
		if !ok {
			t.Errorf("registry missing generator for type %q", ct)
		}
		if gen == nil {
			t.Errorf("registry returned nil generator for type %q", ct)
		}
	}
}

func TestGenerateFromPattern(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	tests := []struct {
		pattern string
		minLen  int
		maxLen  int
	}{
		{"[abc]", 1, 1},
		{"[0-9]", 1, 1},
		{"[a-z]", 1, 1},
		{"(foo|bar)", 3, 3},
		{"hello", 5, 5},
		{"[abc][def]", 2, 2},
	}

	for _, tt := range tests {
		result := generateFromPattern(rng, tt.pattern)
		if len(result) < tt.minLen || len(result) > tt.maxLen {
			t.Errorf("pattern %q produced %q (len=%d), want len in [%d, %d]",
				tt.pattern, result, len(result), tt.minLen, tt.maxLen)
		}
	}
}

// Helpers

func splitSpace(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	return result
}

func splitDot(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '.' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	return result
}

func splitColon(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ':' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	return result
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
