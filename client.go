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

type client struct {
	config *Config
}

func NewClient(config *Config) *client {
	return &client{
		config: config,
	}
}

func (c *client) Get(uri string, options *Options) (Response, error) {
	return c.Request("GET", uri, options)
}

func (c *client) Head(uri string, options *Options) (Response, error) {
	return c.Request("HEAD", uri, options)
}

func (c *client) Post(uri string, options *Options) (Response, error) {
	return c.Request("POST", uri, options)
}

func (c *client) Put(uri string, options *Options) (Response, error) {
	return c.Request("PUT", uri, options)
}

func (c *client) Patch(uri string, options *Options) (Response, error) {
	return c.Request("PATCH", uri, options)
}

func (c *client) Delete(uri string, options *Options) (Response, error) {
	return c.Request("DELETE", uri, options)
}

func (c *client) Request(method, uri string, options *Options) (Response, error) {
	req, err := http.NewRequest(method, c.config.BaseUri+uri, nil)
	if err != nil {
		return nil, err
	}

	if err = c.processOptions(req, options); err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: c.config.Timeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp), nil
}

func (c *client) processOptions(request *http.Request, options *Options) error {
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
		request.Header.Add(key, value)
	}

	return nil
}
