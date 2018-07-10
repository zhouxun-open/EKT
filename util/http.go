package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGet(url string) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func HttpPost(url string, body []byte) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header["Content-Type"] = []string{"application/json"}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
