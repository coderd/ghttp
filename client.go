package ghttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseUri string
	Timeout time.Duration
}

type Options struct {
	Headers map[string]string
	Json    interface{}
	Body    io.Reader
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
	return c.request(http.MethodGet, uri, Options)
}

func (c *client) Head(uri string, options *Options) (Response, error) {
	return c.request(http.MethodHead, uri, Options)
}

func (c *client) Post(uri string, options *Options) (Response, error) {
	return c.request(http.MethodPost, uri, Options)
}

func (c *client) Put(uri string, options *Options) (Response, error) {
	return c.request(http.MethodPut, uri, Options)
}

func (c *client) Patch(uri string, options *Options) (Response, error) {
	return c.request(http.MethodPatch, uri, Options)
}

func (c *client) Delete(uri string, options *Options) (Response, error) {
	return c.request(http.MethodDelete, uri, Options)
}

func (c *client) request(method, uri string, options *Options) (Response, error) {
	var (
		body io.Reader
	)
	if options.Json != nil {
		if json, err := json.Marshal(v); err != nil {
			return err
		}
		body = strings.NewReader(string(json))
	} else if options.Body != nil {
		body = options.Body
	}

	req, err := http.NewRequest(method, c.config.BaseUri+uri, body)
	if err != nil {
		return err
	}

	for key, value := range options.Headers {
		req.Header.Add(key, value)
	}

	httpClient := &http.Client{
		Timeout: c.config.Timeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	return newResponse(resp), nil
}
