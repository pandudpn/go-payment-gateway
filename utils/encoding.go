package utils

import (
	"encoding/base64"
	"fmt"
)

// SetBasicAuthorization set header Basic Authorization
func SetBasicAuthorization(username, password string) string {
	data := fmt.Sprintf("%s:%s", username, password)
	enc := base64.StdEncoding.EncodeToString([]byte(data))

	return fmt.Sprintf("Basic %s", enc)
}
