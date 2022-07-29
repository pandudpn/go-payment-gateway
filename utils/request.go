package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pandudpn/go-payment-gateway"
)

type request struct {
	// http request
	request *http.Request

	// client http client options
	client *http.Client
}

// RequestInterface define all request method
type RequestInterface interface {
	// SetBasicAuth will set your Authorization Basic
	SetBasicAuth(username, password string)

	// SetClient customize client
	SetClient(client *http.Client)

	// SetHeader adding a custom header
	//
	// Default:
	// - Content-Type : application/json
	// - Accept : application/json
	SetHeader(header map[string]string)

	// DoRequest request to target endpoint with context timeout
	// with Header and body payload (if not nil)
	DoRequest(ctx context.Context) ([]byte, int, error)
}

// NewRequest create an instance of http request
func NewRequest(method, url string, payload []byte) (RequestInterface, error) {
	// create instance of http.Request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		pg.Log.Errorf("error create http.Request with message %s", err)
		return nil, err
	}

	// default header
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	// set into headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// create default http.Client
	// default client has timeout 10s and skip the certificate tls
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &request{request: req, client: client}, nil
}

func (r *request) SetBasicAuth(username, password string) {
	r.request.SetBasicAuth(username, password)
}

func (r *request) SetClient(client *http.Client) {
	r.client = client
}

func (r *request) SetHeader(header map[string]string) {
	for k, v := range header {
		r.request.Header.Set(k, v)
	}
}

func (r *request) DoRequest(ctx context.Context) ([]byte, int, error) {
	if r.request == nil {
		return nil, http.StatusInternalServerError, ErrHttpRequest
	}

	if r.client == nil {
		return nil, http.StatusInternalServerError, ErrHttpClient
	}

	pg.Log.Printf("method=%s url=%s", r.request.Method, r.request.URL.String())
	pg.Log.Printf("headers=%v", r.request.Header)
	if r.request.Body != nil {
		reqBody, _ := ioutil.ReadAll(r.request.Body)
		pg.Log.Printf("request_body=%s", string(reqBody))

		// put body
		r.request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	}

	resBody, err := r.client.Do(r.request.WithContext(ctx))
	if err != nil {
		pg.Log.Error(err)
		return nil, http.StatusInternalServerError, err
	}
	defer resBody.Body.Close()

	pg.Log.Printf("response_status_code=%d", resBody.StatusCode)

	body, err := ioutil.ReadAll(resBody.Body)
	if err != nil {
		pg.Log.Error("parse response failed")
		return nil, resBody.StatusCode, err
	}

	pg.Log.Printf("response_body=%s", string(body))

	return body, resBody.StatusCode, nil
}
