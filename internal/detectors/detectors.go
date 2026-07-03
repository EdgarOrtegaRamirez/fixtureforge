package detectors

import (
	"strings"

	"github.com/EdgarOrtegaRamirez/fixtureforge/internal/types"
)

// DetectColumnType attempts to determine the column type from its name
func DetectColumnType(name string) types.ColumnType {
	lower := strings.ToLower(name)
	lower = strings.NewReplacer(
		"_", "", "-", "", " ", "", ".", "",
	).Replace(lower)

	// Exact matches first
	switch lower {
	case "id", "userid", "user_id", "orderid", "order_id", "itemid", "item_id":
		return types.TypeInt
	case "email", "emailaddress", "email_address", "e-mail":
		return types.TypeEmail
	case "phone", "phonenumber", "phone_number", "tel", "telephone", "mobile":
		return types.TypePhone
	case "name", "fullname", "full_name", "displayname", "display_name":
		return types.TypeName
	case "firstname", "first_name", "fname":
		return types.TypeFirstName
	case "lastname", "last_name", "lname", "surname":
		return types.TypeLastName
	case "address", "streetaddress", "street_address", "addr":
		return types.TypeAddress
	case "city", "town", "locality":
		return types.TypeCity
	case "state", "province", "region":
		return types.TypeState
	case "country", "nation":
		return types.TypeCountry
	case "zip", "zipcode", "zip_code", "postalcode", "postal_code":
		return types.TypeZipCode
	case "date", "createddate", "created_date", "updateddate", "updated_date":
		return types.TypeDate
	case "datetime", "timestamp", "createdat", "created_at", "updatedat", "updated_at":
		return types.TypeDateTime
	case "time", "createdtime", "created_time":
		return types.TypeTime
	case "uuid", "uid", "guid":
		return types.TypeUUID
	case "url", "website", "link", "homepage", "uri":
		return types.TypeURL
	case "ip", "ipaddress", "ip_address", "ipv4":
		return types.TypeIPv4
	case "ipv6":
		return types.TypeIPv6
	case "color", "colour", "hexcolor", "hex_color":
		return types.TypeColor
	case "description", "desc", "bio", "biography", "about":
		return types.TypeLorem
	case "active", "enabled", "verified", "approved", "visible", "deleted":
		return types.TypeBool
	case "amount", "price", "cost", "total", "balance", "salary", "revenue":
		return types.TypeFloat
	case "quantity", "qty", "count", "stock", "inventory":
		return types.TypeInt
	case "age":
		return types.TypeRange
	case "status", "type", "role", "level", "priority", "category":
		return types.TypeEnum
	}

	// Substring matches
	if containsAny(lower, "email", "e-mail", "mail") {
		return types.TypeEmail
	}
	if containsAny(lower, "phone", "tel", "mobile", "fax") {
		return types.TypePhone
	}
	if containsAny(lower, "url", "website", "link", "href") {
		return types.TypeURL
	}
	if containsAny(lower, "ip", "ipv4") && !containsAny(lower, "email") {
		return types.TypeIPv4
	}
	if containsAny(lower, "uuid", "guid") {
		return types.TypeUUID
	}
	if containsAny(lower, "date") && containsAny(lower, "time") {
		return types.TypeDateTime
	}
	if containsAny(lower, "datetime", "timestamp") {
		return types.TypeDateTime
	}
	if containsAny(lower, "time") {
		return types.TypeTime
	}
	if containsAny(lower, "date", "day", "year", "month") && !containsAny(lower, "updated", "created") {
		return types.TypeDate
	}
	if containsAny(lower, "first") && containsAny(lower, "name") {
		return types.TypeFirstName
	}
	if containsAny(lower, "last") && containsAny(lower, "name", "surname") {
		return types.TypeLastName
	}
	if containsAny(lower, "fullname", "full_name") {
		return types.TypeName
	}
	if containsAny(lower, "addr", "street", "address") {
		return types.TypeAddress
	}
	if containsAny(lower, "city", "town") {
		return types.TypeCity
	}
	if containsAny(lower, "state", "province") {
		return types.TypeState
	}
	if containsAny(lower, "country", "nation") {
		return types.TypeCountry
	}
	if containsAny(lower, "zip", "postal") {
		return types.TypeZipCode
	}
	if containsAny(lower, "color", "colour", "hex") {
		return types.TypeColor
	}
	if containsAny(lower, "desc", "bio", "about", "summary", "body", "content", "text", "message", "note") {
		return types.TypeLorem
	}
	if containsAny(lower, "price", "cost", "amount", "total", "balance", "salary", "revenue", "discount", "fee", "rate") {
		return types.TypeFloat
	}
	if containsAny(lower, "age") {
		return types.TypeRange
	}
	if containsAny(lower, "active", "enabled", "verified", "is_") || strings.HasPrefix(lower, "is") {
		return types.TypeBool
	}

	// Default to string
	return types.TypeString
}

// containsAny checks if s contains any of the substrings
func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
