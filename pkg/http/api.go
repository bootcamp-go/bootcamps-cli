package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Baseurl = "https://api.github.com"
)

type ApiManager interface {
	Get(url string, resp interface{}) error
	Post(url string, body []byte, resp interface{}) error
	Put(url string, body []byte, resp interface{}) error
}

type apiManager struct {
	client *http.Client
	token  string
}

func NewApiManager(token string) ApiManager {
	return &apiManager{
		client: &http.Client{},
		token:  token,
	}
}

// Get
func (m *apiManager) Get(url string, resp interface{}) error {
	return m.do(url, "GET", nil, resp)
}

// Post
func (m *apiManager) Post(url string, body []byte, resp interface{}) error {
	return m.do(url, "POST", body, resp)
}

// Put
func (m *apiManager) Put(url string, body []byte, resp interface{}) error {
	return m.do(url, "PUT", body, resp)
}

func (m *apiManager) do(url string, method string, body []byte, resp interface{}) error {
	url = Baseurl + url
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+m.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, err := m.client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return fmt.Errorf(fmt.Sprintf("bad response from github api: %d", response.StatusCode))
	}

	// unmarhal response to resp
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return err
	}

	return nil
}
