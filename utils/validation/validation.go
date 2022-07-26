package validation

type countryCode int

const (
	// ID is Indonesian Country
	ID countryCode = iota

	// PH is Philippines Country
	PH
)

// getMapCountryCode mapping to value
var getMapCountryCode = map[string]countryCode{
	"IDR": ID,
	"PHP": PH,
}

// instance for package validation
type validation struct {
	// c is countryCode
	// used for validation phoneNumber
	//
	// Default: ID
	c countryCode
}

// Validation defines all method field Validation
type Validation interface {
	// PhoneNumber validation phone number based on country code
	// e.g: ID (Indonesia) start with +62
	PhoneNumber(phone string) bool
}

// New create an instance of Validation
func New(c string) Validation {
	// get from map
	if v, ok := getMapCountryCode[c]; ok {
		return &validation{
			c: v,
		}
	}

	// if c is not exists (empty string) or not found in map
	// default value: ID
	return &validation{
		c: ID,
	}
}
