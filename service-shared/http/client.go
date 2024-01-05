package http

import (
	"io"
	"io/ioutil"
	"net/http"
)

type ClientResponse struct {
	StatusCode   int
	ResponseBody []byte
}

//Client acts as a wrapper interface to allow for easier unit testing
type Client interface {
	Post(url, contentType string, body io.Reader) (resp *ClientResponse, err error)
	Get(url string) (resp *ClientResponse, err error)
}

type DefaultClient struct {
	HttpClient *http.Client
}

func (client DefaultClient) Get(url string) (*ClientResponse, error) {
	resp, err := client.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return createResponse(resp), nil
}

func (client DefaultClient) Post(url, contentType string, body io.Reader) (*ClientResponse, error) {
	resp, err := client.HttpClient.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return createResponse(resp), nil
}

func createResponse(resp *http.Response) *ClientResponse {
	respBody, _ := ioutil.ReadAll(resp.Body)
	return &ClientResponse{
		StatusCode:   resp.StatusCode,
		ResponseBody: respBody,
	}
}
