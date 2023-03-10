package types

import (
	"strings"
)

var (
	codeToProvince = map[string]Province{
		"ab": 0,   // Alberta
		"bc": 1,   // British Columbia
		"mb": 2,   // Manitoba
		"nb": 3,   // New Brunswick
		"nl": 4,   // Newfoundland and Labrador
		"nt": 5,   // Northwest Territories
		"ns": 6,   // Nova Scotia
		"nu": 7,   // Nunavut
		"on": 8,   // Ontario
		"pe": 9,   // Prince Edward Island
		"qc": 10,  // Quebec
		"sk": 11,  // Saskatchewan
		"yt": 12,  // Yukon
		"xx": 255, // Testing
	}
)

// A type representing a Canadian province
type Province uint8

// Returns a province if the code is recognized,
// otherwise returns false
func ParseProvince(code string) (Province, bool) {
	code = strings.ToLower(code)
	province, ok := codeToProvince[code]
	return province, ok
}

// Returns a string representation of the province
// Returns "??" if unknown
func (p Province) String() string {
	for code, prov := range codeToProvince {
		if prov == p {
			return strings.ToUpper(code)
		}
	}
	return "??"
}

// String, but for JSON marshalling
func (p Province) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}
