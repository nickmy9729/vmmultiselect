package api

import (
	"io/ioutil"
	"net/http"
)

func MakeAPIRequest(endpoint string, headers http.Header) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
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

