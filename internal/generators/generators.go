package generators

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

// Registry holds all registered generators
type Registry struct {
	generators map[types.ColumnType]types.Generator
}

// NewRegistry creates a new generator registry with all built-in generators
func NewRegistry() *Registry {
	r := &Registry{
		generators: make(map[types.ColumnType]types.Generator),
	}
	r.registerDefaults()
	return r
}

// Get returns a generator for the given column type
func (r *Registry) Get(t types.ColumnType) (types.Generator, bool) {
	g, ok := r.generators[t]
	return g, ok
}

// Register adds a custom generator
func (r *Registry) Register(t types.ColumnType, g types.Generator) {
	r.generators[t] = g
}

func (r *Registry) registerDefaults() {
	r.generators[types.TypeString] = &StringGen{}
	r.generators[types.TypeInt] = &IntGen{}
	r.generators[types.TypeFloat] = &FloatGen{}
	r.generators[types.TypeBool] = &BoolGen{}
	r.generators[types.TypeName] = &NameGen{}
	r.generators[types.TypeFirstName] = &FirstNameGen{}
	r.generators[types.TypeLastName] = &LastNameGen{}
	r.generators[types.TypeEmail] = &EmailGen{}
	r.generators[types.TypePhone] = &PhoneGen{}
	r.generators[types.TypeAddress] = &AddressGen{}
	r.generators[types.TypeCity] = &CityGen{}
	r.generators[types.TypeState] = &StateGen{}
	r.generators[types.TypeCountry] = &CountryGen{}
	r.generators[types.TypeZipCode] = &ZipCodeGen{}
	r.generators[types.TypeDate] = &DateGen{}
	r.generators[types.TypeDateTime] = &DateTimeGen{}
	r.generators[types.TypeTime] = &TimeGen{}
	r.generators[types.TypeUUID] = &UUIDGen{}
	r.generators[types.TypeURL] = &URLGen{}
	r.generators[types.TypeIPv4] = &IPv4Gen{}
	r.generators[types.TypeIPv6] = &IPv6Gen{}
	r.generators[types.TypeColor] = &ColorGen{}
	r.generators[types.TypeLorem] = &LoremGen{}
	r.generators[types.TypeEnum] = &EnumGen{}
	r.generators[types.TypeRegex] = &RegexGen{}
	r.generators[types.TypeRange] = &RangeGen{}
	r.generators[types.TypeJSON] = &JSONGen{}
	r.generators[types.TypeBoolWeight] = &BoolWeightGen{}
}

// --- String Generator ---

type StringGen struct{}

func (g *StringGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	minLen := getParamInt(params, "min_length", 3)
	maxLen := getParamInt(params, "max_length", 20)
	prefix := getParamString(params, "prefix", "")
	length := minLen + rng.Intn(maxLen-minLen+1)

	chars := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return prefix + string(b)
}

// --- Int Generator ---

type IntGen struct{}

func (g *IntGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	min := getParamInt(params, "min", 0)
	max := getParamInt(params, "max", 1000)
	return int64(min + rng.Intn(max-min+1))
}

// --- Float Generator ---

type FloatGen struct{}

func (g *FloatGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	min := getParamFloat(params, "min", 0)
	max := getParamFloat(params, "max", 1000)
	decimals := getParamInt(params, "decimals", 2)
	val := min + rng.Float64()*(max-min)
	mult := 1.0
	for i := 0; i < decimals; i++ {
		mult *= 10
	}
	return float64(int(val*mult)) / mult
}

// --- Bool Generator ---

type BoolGen struct{}

func (g *BoolGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return rng.Intn(2) == 0
}

// --- BoolWeight Generator (configurable true/false probability) ---

type BoolWeightGen struct{}

func (g *BoolWeightGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	truePct := getParamFloat(params, "true_pct", 50)
	return rng.Float64()*100 < truePct
}

// --- Name Generator ---

type NameGen struct{}

var firstNames = []string{
	"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda",
	"David", "Elizabeth", "William", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Christopher", "Karen", "Charles", "Lisa", "Daniel", "Nancy",
	"Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra", "Donald", "Ashley",
	"Amanda", "Emily", "Andrew", "Donna", "Joshua", "Michelle", "Kenneth", "Carol",
	"Kevin", "Amanda", "Brian", "Dorothy", "George", "Melissa", "Timothy", "Deborah",
	"Ronald", "Stephanie", "Edward", "Rebecca", "Jason", "Sharon", "Jeffrey", "Laura",
	"Ryan", "Cynthia", "Jacob", "Kathleen", "Gary", "Amy", "Nicholas", "Angela",
	"Eric", "Shirley", "Jonathan", "Anna", "Stephen", "Brenda", "Larry", "Pamela",
	"Justin", "Emma", "Scott", "Nicole", "Brandon", "Helen", "Benjamin", "Samantha",
	"Samuel", "Katherine", "Raymond", "Christine", "Gregory", "Debra", "Frank", "Rachel",
	"Alexander", "Carolyn", "Patrick", "Janet", "Jack", "Catherine", "Dennis", "Maria",
	"Jerry", "Heather", "Tyler", "Diane", "Aaron", "Ruth", "Jose", "Julie",
	"Adam", "Olivia", "Nathan", "Joyce", "Henry", "Virginia",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
	"Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill",
	"Flores", "Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell",
	"Mitchell", "Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz",
	"Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales",
	"Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson",
	"Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward",
	"Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray",
	"Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel",
	"Myers", "Long", "Ross", "Foster", "Jimenez",
}

func (g *NameGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return firstNames[rng.Intn(len(firstNames))] + " " + lastNames[rng.Intn(len(lastNames))]
}

// --- FirstName Generator ---

type FirstNameGen struct{}

func (g *FirstNameGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return firstNames[rng.Intn(len(firstNames))]
}

// --- LastName Generator ---

type LastNameGen struct{}

func (g *LastNameGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return lastNames[rng.Intn(len(lastNames))]
}

// --- Email Generator ---

type EmailGen struct{}

var emailDomains = []string{
	"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "aol.com",
	"icloud.com", "mail.com", "protonmail.com", "zoho.com", "yandex.com",
}

func (g *EmailGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	first := strings.ToLower(firstNames[rng.Intn(len(firstNames))])
	last := strings.ToLower(lastNames[rng.Intn(len(lastNames))])
	domain := emailDomains[rng.Intn(len(emailDomains))]

	sep := []string{".", "_", ""}[rng.Intn(3)]
	num := ""
	if rng.Intn(3) == 0 {
		num = fmt.Sprintf("%d", rng.Intn(1000))
	}

	return first + sep + last + num + "@" + domain
}

// --- Phone Generator ---

type PhoneGen struct{}

func (g *PhoneGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	format := getParamString(params, "format", "us")
	switch format {
	case "us":
		return fmt.Sprintf("(%03d) %03d-%03d",
			200+rng.Intn(800), rng.Intn(1000), rng.Intn(10000))
	case "intl":
		return fmt.Sprintf("+1-%03d-%03d-%04d",
			200+rng.Intn(800), rng.Intn(1000), rng.Intn(10000))
	case "e164":
		return fmt.Sprintf("+1%03d%03d%04d",
			200+rng.Intn(800), rng.Intn(1000), rng.Intn(10000))
	default:
		return fmt.Sprintf("%03d-%03d-%04d",
			200+rng.Intn(800), rng.Intn(1000), rng.Intn(10000))
	}
}

// --- Address Generator ---

type AddressGen struct{}

var streetNames = []string{
	"Main", "Oak", "Maple", "Cedar", "Elm", "Pine", "Walnut", "Chestnut",
	"Birch", "Spruce", "Willow", "Ash", "Poplar", "Cypress", "Hickory",
	"First", "Second", "Third", "Fourth", "Fifth", "Sixth", "Seventh",
	"Park", "Lake", "Hill", "River", "Valley", "Forest", "Meadow",
}

var streetTypes = []string{"St", "Ave", "Blvd", "Dr", "Ln", "Rd", "Way", "Ct", "Pl", "Ter"}

func (g *AddressGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	num := 1 + rng.Intn(9999)
	street := streetNames[rng.Intn(len(streetNames))]
	stType := streetTypes[rng.Intn(len(streetTypes))]
	return fmt.Sprintf("%d %s %s", num, street, stType)
}

// --- City Generator ---

type CityGen struct{}

var cities = []string{
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
	"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
	"Austin", "Jacksonville", "Fort Worth", "Columbus", "Charlotte",
	"Indianapolis", "San Francisco", "Seattle", "Denver", "Washington",
	"Nashville", "Oklahoma City", "El Paso", "Boston", "Portland",
	"Las Vegas", "Memphis", "Louisville", "Baltimore", "Milwaukee",
	"Albuquerque", "Tucson", "Fresno", "Sacramento", "Mesa",
	"Kansas City", "Atlanta", "Omaha", "Colorado Springs", "Raleigh",
	"Long Beach", "Virginia Beach", "Miami", "Oakland", "Minneapolis",
	"Tulsa", "Tampa", "Arlington", "New Orleans", "Wichita",
}

func (g *CityGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return cities[rng.Intn(len(cities))]
}

// --- State Generator ---

type StateGen struct{}

var states = []string{
	"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
	"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
	"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
	"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
	"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
}

func (g *StateGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return states[rng.Intn(len(states))]
}

// --- Country Generator ---

type CountryGen struct{}

var countries = []string{
	"United States", "Canada", "United Kingdom", "Germany", "France",
	"Japan", "Australia", "Brazil", "India", "Mexico",
	"Spain", "Italy", "South Korea", "Netherlands", "Sweden",
	"Switzerland", "Norway", "Denmark", "Finland", "Belgium",
	"Austria", "Ireland", "New Zealand", "Singapore", "Portugal",
}

func (g *CountryGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return countries[rng.Intn(len(countries))]
}

// --- ZipCode Generator ---

type ZipCodeGen struct{}

func (g *ZipCodeGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	format := getParamString(params, "format", "us5")
	switch format {
	case "us5":
		return fmt.Sprintf("%05d", rng.Intn(100000))
	case "us9":
		return fmt.Sprintf("%05d-%04d", rng.Intn(100000), rng.Intn(10000))
	case "uk":
		letters := "ABCDEFGHJKLMNPQRSTUWXYZ"
		return fmt.Sprintf("%s%d%s %d%s%s",
			string(letters[rng.Intn(len(letters))]),
			rng.Intn(10),
			string(letters[rng.Intn(len(letters))]),
			rng.Intn(10),
			string(letters[rng.Intn(len(letters))]),
			string(letters[rng.Intn(len(letters))]))
	default:
		return fmt.Sprintf("%05d", rng.Intn(100000))
	}
}

// --- Date Generator ---

type DateGen struct{}

func (g *DateGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	year := 2000 + rng.Intn(26)
	month := 1 + rng.Intn(12)
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		daysInMonth[1] = 29
	}
	day := 1 + rng.Intn(daysInMonth[month-1])

	format := getParamString(params, "format", "iso")
	switch format {
	case "iso":
		return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	case "us":
		return fmt.Sprintf("%02d/%02d/%04d", month, day, year)
	case "eu":
		return fmt.Sprintf("%02d.%02d.%04d", day, month, year)
	default:
		return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}
}

// --- DateTime Generator ---

type DateTimeGen struct{}

func (g *DateTimeGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	year := 2000 + rng.Intn(26)
	month := 1 + rng.Intn(12)
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		daysInMonth[1] = 29
	}
	day := 1 + rng.Intn(daysInMonth[month-1])
	hour := rng.Intn(24)
	minute := rng.Intn(60)
	second := rng.Intn(60)

	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", year, month, day, hour, minute, second)
}

// --- Time Generator ---

type TimeGen struct{}

func (g *TimeGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	hour := rng.Intn(24)
	minute := rng.Intn(60)
	second := rng.Intn(60)
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

// --- UUID Generator ---

type UUIDGen struct{}

func (g *UUIDGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	b := make([]byte, 16)
	_, _ = rng.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// --- URL Generator ---

type URLGen struct{}

var urlSchemes = []string{"https", "http"}
var urlDomains = []string{
	"example.com", "sample.org", "test.io", "demo.dev", "mock.app",
	"fake.co", "dummy.net", "placeholder.com", "virtual.io", "stub.dev",
}
var urlPaths = []string{
	"", "/index.html", "/about", "/contact", "/products",
	"/api/v1/users", "/api/v1/items", "/docs", "/help", "/login",
	"/dashboard", "/settings", "/profile", "/search", "/blog",
}

func (g *URLGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	scheme := urlSchemes[rng.Intn(len(urlSchemes))]
	domain := urlDomains[rng.Intn(len(urlDomains))]
	path := urlPaths[rng.Intn(len(urlPaths))]
	return fmt.Sprintf("%s://%s%s", scheme, domain, path)
}

// --- IPv4 Generator ---

type IPv4Gen struct{}

func (g *IPv4Gen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return fmt.Sprintf("%d.%d.%d.%d",
		1+rng.Intn(223), rng.Intn(256), rng.Intn(256), 1+rng.Intn(254))
}

// --- IPv6 Generator ---

type IPv6Gen struct{}

func (g *IPv6Gen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	return fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x",
		rng.Intn(0x10000), rng.Intn(0x10000), rng.Intn(0x10000), rng.Intn(0x10000),
		rng.Intn(0x10000), rng.Intn(0x10000), rng.Intn(0x10000), rng.Intn(0x10000))
}

// --- Color Generator ---

type ColorGen struct{}

func (g *ColorGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	format := getParamString(params, "format", "hex")
	switch format {
	case "hex":
		return fmt.Sprintf("#%06x", rng.Intn(0x1000000))
	case "rgb":
		return fmt.Sprintf("rgb(%d, %d, %d)", rng.Intn(256), rng.Intn(256), rng.Intn(256))
	case "hsl":
		return fmt.Sprintf("hsl(%d, %d%%, %d%%)", rng.Intn(360), 30+rng.Intn(71), 20+rng.Intn(61))
	default:
		return fmt.Sprintf("#%06x", rng.Intn(0x1000000))
	}
}

// --- Lorem Generator ---

type LoremGen struct{}

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum", "perspiciatis", "unde",
	"omnis", "iste", "natus", "error", "voluptatem", "accusantium", "doloremque",
	"laudantium", "totam", "rem", "aperiam", "eaque", "ipsa", "quae", "ab", "illo",
	"inventore", "veritatis", "quasi", "architecto", "beatae", "vitae", "dicta",
}

func (g *LoremGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	wordCount := getParamInt(params, "words", 10)
	words := make([]string, wordCount)
	for i := range words {
		words[i] = loremWords[rng.Intn(len(loremWords))]
	}
	return strings.Join(words, " ")
}

// --- Enum Generator ---

type EnumGen struct{}

func (g *EnumGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	values := getParamString(params, "values", "a,b,c")
	parts := strings.Split(values, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts[rng.Intn(len(parts))]
}

// --- Regex Generator ---

type RegexGen struct{}

func (g *RegexGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	pattern := getParamString(params, "pattern", "[a-z]{5}")
	return generateFromPattern(rng, pattern)
}

// --- Range Generator ---

type RangeGen struct{}

func (g *RangeGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	min := getParamFloat(params, "min", 1)
	max := getParamFloat(params, "max", 100)
	step := getParamFloat(params, "step", 1)
	if step <= 0 {
		step = 1
	}
	val := min + rng.Float64()*(max-min)
	val = float64(int(val/step)) * step
	decimals := getParamInt(params, "decimals", 0)
	if decimals == 0 {
		return int64(val)
	}
	mult := 1.0
	for i := 0; i < decimals; i++ {
		mult *= 10
	}
	return float64(int(val*mult)) / mult
}

// --- JSON Generator ---

type JSONGen struct{}

func (g *JSONGen) Generate(rng *rand.Rand, params map[string]string) interface{} {
	// Generate a simple nested JSON object
	obj := map[string]interface{}{
		"id":     rng.Intn(10000),
		"name":   firstNames[rng.Intn(len(firstNames))],
		"email":  firstNames[rng.Intn(len(firstNames))] + "@example.com",
		"active": rng.Intn(2) == 0,
	}
	return obj
}

// --- Helper functions ---

func getParamInt(params map[string]string, key string, defaultVal int) int {
	if params == nil {
		return defaultVal
	}
	if v, ok := params[key]; ok {
		var n int
		_, err := fmt.Sscanf(v, "%d", &n)
		if err == nil {
			return n
		}
	}
	return defaultVal
}

func getParamFloat(params map[string]string, key string, defaultVal float64) float64 {
	if params == nil {
		return defaultVal
	}
	if v, ok := params[key]; ok {
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		if err == nil {
			return f
		}
	}
	return defaultVal
}

func getParamString(params map[string]string, key string, defaultVal string) string {
	if params == nil {
		return defaultVal
	}
	if v, ok := params[key]; ok {
		return v
	}
	return defaultVal
}

// generateFromPattern generates a string from a simple regex-like pattern
// Supports: [a-z], [0-9], [A-Z], (a|b|c)
func generateFromPattern(rng *rand.Rand, pattern string) string {
	var result strings.Builder
	i := 0
	for i < len(pattern) {
		if pattern[i] == '[' && i+2 < len(pattern) && pattern[i+1] == '^' {
			end := strings.IndexByte(pattern[i:], ']')
			if end > 0 {
				chars := pattern[i+2 : i+end]
				for j := 0; j < 5; j++ {
					result.WriteByte(chars[rng.Intn(len(chars))])
				}
				i += end + 1
				continue
			}
		} else if pattern[i] == '[' {
			end := strings.IndexByte(pattern[i:], ']')
			if end > 0 {
				chars := pattern[i+1 : i+end]
				result.WriteByte(chars[rng.Intn(len(chars))])
				i += end + 1
				continue
			}
		} else if pattern[i] == '(' {
			end := strings.IndexByte(pattern[i:], ')')
			if end > 0 {
				options := strings.Split(pattern[i+1:i+end], "|")
				result.WriteString(options[rng.Intn(len(options))])
				i += end + 1
				continue
			}
		}
		result.WriteByte(pattern[i])
		i++
	}
	return result.String()
}
