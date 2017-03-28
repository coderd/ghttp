package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
	BaseUri string
	Timeout time.Duration
}

type Options struct {
	Headers map[string]string
	Json    interface{}
	Body    io.ReadCloser
}

type Client struct {
	Config *Config
}

func NewClient(config *Config) *Client {
	return &Client{
		Config: config,
	}
}

func (c *Client) Get(uri string, options *Options) (Response, error) {
	return c.Request("GET", uri, options)
}

func (c *Client) Head(uri string, options *Options) (Response, error) {
	return c.Request("HEAD", uri, options)
}

func (c *Client) Post(uri string, options *Options) (Response, error) {
	return c.Request("POST", uri, options)
}

func (c *Client) Put(uri string, options *Options) (Response, error) {
	return c.Request("PUT", uri, options)
}

func (c *Client) Patch(uri string, options *Options) (Response, error) {
	return c.Request("PATCH", uri, options)
}

func (c *Client) Delete(uri string, options *Options) (Response, error) {
	return c.Request("DELETE", uri, options)
}

func (c *Client) Request(method, uri string, options *Options) (Response, error) {
	req, err := http.NewRequest(method, c.Config.BaseUri+uri, nil)
	if err != nil {
		return nil, err
	}

	if err = c.processOptions(req, options); err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: c.Config.Timeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp), nil
}

func (c *Client) processOptions(request *http.Request, options *Options) error {
	if options == nil {
		return nil
	}

	// Process body
	var body io.ReadCloser
	if options.Json != nil {
		b, err := json.Marshal(options.Json)
		if err != nil {
			return err
		}
		body = ioutil.NopCloser(bytes.NewReader(b))

		request.Header.Add("Content-Type", "application/json")
	} else if options.Body != nil {
		body = options.Body
	}
	request.Body = body

	// Process headers
	for key, value := range options.Headers {
		if key == "Host" {
			request.Host = value
		} else {
			request.Header.Add(key, value)
		}
	}

	return nil
}
