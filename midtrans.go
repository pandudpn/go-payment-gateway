package pg

const (
	mdUriSandbox    string = "https://api.sandbox.midtrans.com"
	mdUriProduction string = "https://api.midtrans.com"
)

// Midtrans configuration
type Midtrans struct {
	// uri is base url of Midtrans Core API
	uri string
	
	// credentials key for access Midtrans Core API
	credentials *Credentials
}
