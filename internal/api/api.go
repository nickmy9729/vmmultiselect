package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// MakeAPIRequest makes an HTTP GET request to the specified endpoint with the provided headers
// and returns the response body as a byte slice.
func MakeAPIRequest(endpoint string, headers http.Header) ([]byte, error) {
	return MakeAPIRequestWithTimeout(endpoint, headers, 0)
}

// MakeAPIRequestWithTimeout makes an HTTP GET request to the specified endpoint with the provided headers
// and a specified timeout. It returns the response body as a byte slice.
func MakeAPIRequestWithTimeout(endpoint string, headers http.Header, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

