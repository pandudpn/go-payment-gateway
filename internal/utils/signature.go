package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
)

// SignatureVerifier defines the interface for signature verification
type SignatureVerifier interface {
	Verify(r *http.Request) bool
}

// SHA512Verifier verifies SHA512-based signatures (used by Midtrans)
type SHA512Verifier struct {
	secretKey string
	orderIDKey   string
	statusKey   string
}

// NewSHA512Verifier creates a new SHA512 signature verifier
func NewSHA512Verifier(secretKey, orderIDKey, statusKey string) *SHA512Verifier {
	return &SHA512Verifier{
		secretKey: secretKey,
		orderIDKey: orderIDKey,
		statusKey: statusKey,
	}
}

// Verify verifies the SHA512 signature from the request
// Midtrans signature: SHA512(orderID + status + serverKey)
func (v *SHA512Verifier) Verify(r *http.Request) bool {
	// Get signature from header
	signature := r.Header.Get("X-Signature")
	if signature == "" {
		return false
	}

	// Parse body to get order_id and status
	if err := r.ParseForm(); err != nil {
		return false
	}

	orderID := r.FormValue(v.orderIDKey)
	status := r.FormValue(v.statusKey)

	// Calculate expected signature
	expected := v.calculate(orderID, status)

	// Use hmac.Equal to prevent timing attacks
	return hmacEqual([]byte(signature), []byte(expected))
}

// calculate calculates the SHA512 signature
func (v *SHA512Verifier) calculate(orderID, status string) string {
	payload := fmt.Sprintf("%s%s%s", orderID, status, v.secretKey)
	h := sha512.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyString verifies a signature string directly
func (v *SHA512Verifier) VerifyString(orderID, status, signature string) bool {
	expected := v.calculate(orderID, status)
	return hmacEqual([]byte(signature), []byte(expected))
}

// HMACVerifier verifies HMAC-based signatures
type HMACVerifier struct {
	secretKey  string
	headerKey  string
	hashFunc   func() hash.Hash
}

// NewHMACVerifier creates a new HMAC signature verifier
func NewHMACVerifier(secretKey, headerKey string, hashType string) *HMACVerifier {
	var hf func() hash.Hash
	switch strings.ToLower(hashType) {
	case "sha256":
		hf = sha256.New
	case "sha512":
		hf = sha512.New
	default:
		hf = sha256.New
	}

	return &HMACVerifier{
		secretKey: secretKey,
		headerKey: headerKey,
		hashFunc:  hf,
	}
}

// Verify verifies the HMAC signature from the request
func (v *HMACVerifier) Verify(r *http.Request) bool {
	// Get signature from header
	signature := r.Header.Get(v.headerKey)
	if signature == "" {
		return false
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}
	// Restore body for subsequent reads
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	// Calculate expected signature
	h := hmac.New(v.hashFunc, []byte(v.secretKey))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))

	return hmacEqual([]byte(signature), []byte(expected))
}

// VerifyBytes verifies HMAC signature for raw bytes
func (v *HMACVerifier) VerifyBytes(data []byte, signature string) bool {
	h := hmac.New(v.hashFunc, []byte(v.secretKey))
	h.Write(data)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmacEqual([]byte(signature), []byte(expected))
}

// WebhookVerifier is a general-purpose webhook signature verifier
type WebhookVerifier struct {
	verifier SignatureVerifier
}

// NewWebhookVerifier creates a new webhook verifier
func NewWebhookVerifier(verifier SignatureVerifier) *WebhookVerifier {
	return &WebhookVerifier{
		verifier: verifier,
	}
}

// Verify verifies the webhook signature
func (w *WebhookVerifier) Verify(r *http.Request) bool {
	if w.verifier == nil {
		return false
	}
	return w.verifier.Verify(r)
}

// CallbackTokenVerifier verifies callback token from headers (used by Xendit)
type CallbackTokenVerifier struct {
	token string
}

// NewCallbackTokenVerifier creates a new callback token verifier
func NewCallbackTokenVerifier(token string) *CallbackTokenVerifier {
	return &CallbackTokenVerifier{
		token: token,
	}
}

// Verify verifies the callback token from the request header
func (v *CallbackTokenVerifier) Verify(r *http.Request) bool {
	callbackToken := r.Header.Get("X-Callback-Token")
	return callbackToken == v.token
}

// hmacEqual compares two HMAC strings in constant time to prevent timing attacks
func hmacEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// CalculateSHA512 calculates SHA512 hash of a string
func CalculateSHA512(data string) string {
	h := sha512.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// CalculateSHA256 calculates SHA256 hash of a string
func CalculateSHA256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// CalculateHMACSHA512 calculates HMAC-SHA512 of data with key
func CalculateHMACSHA512(key, data string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// CalculateHMACSHA256 calculates HMAC-SHA256 of data with key
func CalculateHMACSHA256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SetBasicAuthorization creates a Basic Auth header value
func SetBasicAuthorization(username, password string) string {
	if username == "" {
		return ""
	}
	const basicAuthPrefix = "Basic "
	auth := username + ":" + password
	return basicAuthPrefix + encodeBase64(auth)
}

// encodeBase64 encodes a string to base64
func encodeBase64(data string) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	var result strings.Builder
	remaining := len(data)

	// Process 3 bytes at a time
	for i := 0; i < len(data); i += 3 {
		remaining = len(data) - i

		var n uint32
		switch remaining {
		case 1:
			n = uint32(data[i]) << 16
		case 2:
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8
		default:
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8 | uint32(data[i+2])
		}

		// Write 4 base64 characters
		result.WriteByte(base64Chars[(n>>18)&0x3F])
		result.WriteByte(base64Chars[(n>>12)&0x3F])

		if remaining >= 2 {
			result.WriteByte(base64Chars[(n>>6)&0x3F])
		} else {
			result.WriteByte('=')
		}

		if remaining >= 3 {
			result.WriteByte(base64Chars[n&0x3F])
		} else {
			result.WriteByte('=')
		}
	}

	return result.String()
}
