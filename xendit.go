package pg

const xndUri string = "https://api.xendit.co"

type Xendit struct {
	// uri is base url of Xendit API
	uri string
	
	// credentials key for access Xendit API
	credentials *Credentials
}
