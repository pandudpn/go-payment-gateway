package pg

import (
	"strings"
)

type Credentials struct {
	// key client of Payment Gateway
	//
	// Default: ""
	ClientId string `json:"client_id"`

	// secret client of Payment Gateway
	//
	// Default: ""
	ClientSecret string `json:"client_secret"`
}

// checkCredentials is valid for the environment
//
// @param v is value. it can be clientId, clientSecret, apiKey, or secretKey
// @param g is payment gateway. e.g: Midtrans, Xendit, oy
func checkCredentials(en env, v, g string) bool {
	if v == "" {
		return false
	}

	// force value to lower
	v = strings.ToLower(v)

	if en == Production {
		switch g {
		case Midtrans:
			vs := strings.Split(v, "-")
			if len(vs) > 0 {
				return vs[0] != "sb"
			}
			return false
		case Xendit:
			return strings.Contains(v, "xnd_production")
		default:
			return false
		}
	}

	switch g {
	case Midtrans:
		vs := strings.Split(v, "-")
		if len(vs) > 0 {
			return vs[0] == "sb"
		}
		return false
	case Xendit:
		return strings.Contains(v, "xnd_development")
	default:
		return false
	}
}
