package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type Doer interface {
	Do(request *http.Request) (*http.Response, error)
}

type Client struct {
	HttpClient *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{HttpClient: client}
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	token, err := ReadToken()
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", token)
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, errors.New("Server is down\n")
	}
	if response.StatusCode > 399 {
		defer response.Body.Close()
		result, _ := ioutil.ReadAll(response.Body)
		return nil, errors.New(string(result))
	}
	return response, nil
}
