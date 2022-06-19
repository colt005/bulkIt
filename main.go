package main

import (
	"fmt"
	"net/http"
)

type ImageSave struct {
	client *HTTPClient
}

type HTTPClient struct {
	client *http.Client
}

func NewImageSave(c *http.Client) *ImageSave {
	if c == nil {
		c = &http.Client{}
	}

	h := &HTTPClient{
		client: c,
	}

	return &ImageSave{
		client: h,
	}
}

func (c *HTTPClient) GET(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestRunFunc() {
	fmt.Println("Running")
}

func main() {
	TestRunFunc()
}
