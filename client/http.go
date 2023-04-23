// Package client contains a HTTP client.
package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Parameters provides the parameters used when creating a new HTTP client.
type Parameters struct {
	Timeout *time.Duration
}

// NewHTTPClient instantiates a new HTTPClient based on provided parameters.
func NewHTTPClient(parameters Parameters) HTTPClient {
	if parameters.Timeout == nil {
		timeout := 1 * time.Second
		parameters.Timeout = &timeout
	}

	client := &http.Client{
		Timeout: *parameters.Timeout,
	}

	return HTTPClient{client}
}

// HTTPRequestData contains the request data.
type HTTPRequestData struct {
	Method     string
	URL        string
	Headers    map[string]string
	Body       []byte
	GetPayload *url.Values
}

// HTTPClient contains the HTTP client.
type HTTPClient struct {
	*http.Client
}

// HTTPStatusCodeError is an error that occurs when receiving an unexpected status code (>= 400).
type HTTPStatusCodeError struct {
	URL        string
	StatusCode int
	Message    string
}

// Error return an error string.
func (e HTTPStatusCodeError) Error() string {
	return fmt.Sprintf("Error response from %s, got status: %d, message: %s", e.URL, e.StatusCode, e.Message)
}

// RequestBytes does the actual HTTP request.
// Returns a slice of bytes or an error.
func (client *HTTPClient) RequestBytes(reqData HTTPRequestData) ([]byte, error) {

	r, err := client.request(reqData)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode >= 400 {
		resp, _ := ioutil.ReadAll(r.Body)
		message := string(resp)
		return nil, HTTPStatusCodeError{
			URL:        reqData.URL,
			StatusCode: r.StatusCode,
			Message:    message,
		}
	}

	return ioutil.ReadAll(r.Body)
}

func (client *HTTPClient) request(reqData HTTPRequestData) (*http.Response, error) {
	req, err := http.NewRequest(reqData.Method, reqData.URL, bytes.NewBuffer(reqData.Body))
	if err != nil {
		return nil, err
	}

	if reqData.GetPayload != nil {
		req.URL.RawQuery = reqData.GetPayload.Encode()
	}

	for k, v := range reqData.Headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("User-Agent", "auth-api")

	resp, err := client.Do(req)

	if err != nil {
		if reqData.Method == http.MethodPost {
			return resp, fmt.Errorf("error making request: %v. Body: %s", err, reqData.Body)
		}

		return resp, fmt.Errorf("error making request: %v. Query: %v", err, req.URL.RawQuery)
	}

	return resp, nil
}
