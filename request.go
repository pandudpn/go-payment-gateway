package pg

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var defaultTimeout = time.Duration(10) * time.Second

type ApiRequest struct {
	// HttpClient http client request
	HttpClient *http.Client
	
	// Logging logging interface
	Logging Logging
}

// ApiRequestInterface defines all api request
type ApiRequestInterface interface {
	// Call do request to target
	Call(ctx context.Context, httpMethod, url string, header http.Header, body, result interface{}) error
	
	// DoRequest to target
	DoRequest(req *http.Request, result interface{}) error
}

// DefaultApiRequest is default of ApiRequest instance
func DefaultApiRequest() ApiRequestInterface {
	client := &http.Client{Timeout: defaultTimeout, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	
	return &ApiRequest{
		HttpClient: client,
		Logging:    Log,
	}
}

func (a *ApiRequest) Call(ctx context.Context, httpMethod, url string, header http.Header, body, result interface{}) error {
	var (
		err     error
		reqBody []byte
	)
	
	// validation body when method is Post or Put
	// hasBody := body != nil || (reflect.ValueOf(body).Kind() != reflect.Ptr && !reflect.ValueOf(body).IsNil())
	hasBody := body != nil
	if hasBody {
		// handling error panic
		switch body.(type) {
		case []byte:
			reqBody = body.([]byte)
		default:
			reqBody, err = json.Marshal(body)
			if err != nil {
				a.Logging.Errorf("error marshal body %s", err)
				return err
			}
		}
	}
	
	// generate New Http.Request
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewBuffer(reqBody))
	if err != nil {
		a.Logging.Errorf("cannot create midtrans request: %s", err)
		return err
	}
	
	// set header
	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("go-payment-gateway/%s", Version))
	
	// log request
	a.Logging.Println("================================= REQUEST =================================")
	a.Logging.Println(req.Method, req.URL, req.Proto)
	a.logHttpHeader(req.Header)
	a.Logging.Println("BODY:", string(reqBody), "\n")
	
	return a.DoRequest(req, result)
}

// DoRequest to client with some params
func (a *ApiRequest) DoRequest(req *http.Request, result interface{}) error {
	start := time.Now()
	
	// request to client target
	res, err := a.HttpClient.Do(req)
	if err != nil {
		a.Logging.Errorf("cannot send request %s", err)
		return err
	}
	defer res.Body.Close()
	
	a.Logging.Println("================================= END =================================")
	a.Logging.Printf("Request completed in %.4fs\n", time.Since(start).Seconds())
	
	// read body response
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		a.Logging.Errorf("Cannot read response body: %s", err)
		return err
	}
	
	// log response
	a.Logging.Println("================================= RESPONSE =================================")
	a.Logging.Println(res.Proto, res.Status)
	a.logHttpHeader(res.Header)
	a.Logging.Println("BODY:", string(resBody))
	
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		a.Logging.Errorf("parse body response into struct failed: %s", err)
		return err
	}
	
	return nil
}

func (a *ApiRequest) logHttpHeader(header http.Header) {
	a.Logging.Println("HEADERS:")
	
	for name, value := range header {
		name = strings.ToLower(name)
		
		a.Logging.Printf("%s: %v", name, value)
	}
}
