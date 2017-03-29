package ghttp

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	Addr = "127.0.0.1:9362"
)

type httpHandler struct {
	t *testing.T
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	action := req.Method + " " + req.URL.Path

	switch action {
	case "GET /foo":
		w.WriteHeader(200)
		w.Write([]byte("bar"))

	case "HEAD /foo":
		w.WriteHeader(200)

	case "POST /foo":
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			h.t.Error(err)
		}
		w.WriteHeader(201)
		w.Write(b)

	case "PUT /foo":
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			h.t.Error(err)
		}
		w.WriteHeader(200)
		w.Write(b)

	case "PATCH /foo":
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			h.t.Error(err)
		}
		w.WriteHeader(200)
		w.Write(b)

	case "DELETE /foo":
		w.WriteHeader(204)

	case "GET /host":
		w.WriteHeader(200)
		w.Write([]byte(req.Host))

	case "GET /header":
		w.WriteHeader(200)
		w.Write([]byte(req.Header.Get("X-Foo")))
	}
}

func serveHttp(t *testing.T) {
	s := &http.Server{
		Addr:         Addr,
		Handler:      &httpHandler{t: t},
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	t.Fatal(s.ListenAndServe())
}

func TestRequest(t *testing.T) {
	go serveHttp(t)
	time.Sleep(time.Millisecond * 5)

	for _, test := range tests {
		config := &Config{
			BaseUri: "http://" + Addr,
			Timeout: 2 * time.Second,
		}
		c := NewClient(config)

		var (
			resp Response
			err  error
		)
		switch test.method {
		case "GET":
			resp, err = c.Get(test.path, test.options)
		case "HEAD":
			resp, err = c.Head(test.path, test.options)
		case "POST":
			resp, err = c.Post(test.path, test.options)
		case "PUT":
			resp, err = c.Put(test.path, test.options)
		case "PATCH":
			resp, err = c.Patch(test.path, test.options)
		case "DELETE":
			resp, err = c.Delete(test.path, test.options)
		}

		if err != nil {
			t.Error(err)
		}

		status := resp.GetStatusCode()
		if status != test.expectedStatus {
			t.Errorf("Got status %d, expected %d", status, test.expectedStatus)
		}

		body := resp.GetBody()
		b, err := ioutil.ReadAll(body)
		if err != nil {
			t.Error(err)
		}

		if string(b) != test.expectedBody {
			t.Errorf("Got body %s, expected %s", b, test.expectedBody)
		}
	}
}

var tests = []struct {
	method         string
	path           string
	options        *Options
	expectedStatus int
	expectedBody   string
}{
	{
		method:         "GET",
		path:           "/foo",
		options:        nil,
		expectedStatus: 200,
		expectedBody:   "bar",
	},
	{
		method:         "HEAD",
		path:           "/foo",
		options:        nil,
		expectedStatus: 200,
		expectedBody:   "",
	},
	{
		method: "POST",
		path:   "/foo",
		options: &Options{
			Json: struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
		expectedStatus: 201,
		expectedBody:   `{"Foo":"bar"}`,
	},
	{
		method: "POST",
		path:   "/foo",
		options: &Options{
			Body: ioutil.NopCloser(strings.NewReader("bar")),
		},
		expectedStatus: 201,
		expectedBody:   `bar`,
	},
	{
		method: "PUT",
		path:   "/foo",
		options: &Options{
			Json: struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
		expectedStatus: 200,
		expectedBody:   `{"Foo":"bar"}`,
	},
	{
		method: "PATCH",
		path:   "/foo",
		options: &Options{
			Json: struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
		expectedStatus: 200,
		expectedBody:   `{"Foo":"bar"}`,
	},
	{
		method:         "DELETE",
		path:           "/foo",
		options:        nil,
		expectedStatus: 204,
		expectedBody:   "",
	},

	{
		method: "GET",
		path:   "/host",
		options: &Options{
			Headers: map[string]string{
				"Host": "example.com",
			},
		},
		expectedStatus: 200,
		expectedBody:   "example.com",
	},
	{
		method: "GET",
		path:   "/header",
		options: &Options{
			Headers: map[string]string{
				"X-Foo": "bar",
			},
		},
		expectedStatus: 200,
		expectedBody:   "bar",
	},
}
