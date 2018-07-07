package techo

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type Marshaler func(interface{}) ([]byte, error)

// Response represents HTTP response.
type Response struct {
	Code   int
	Body   []byte
	Header http.Header
}

// Client represents test echo client.
type Client struct {
	e         *echo.Echo
	marshaler Marshaler
	Header    http.Header
}

// New creates new echo client.
func New(e *echo.Echo, marshaler Marshaler) *Client {
	return &Client{
		e:         e,
		marshaler: marshaler,
		Header:    map[string][]string{},
	}
}

func (c *Client) request(method string, url string, headers http.Header, body interface{}) Response {
	var br io.Reader
	if body != nil {
		b, e := c.marshaler(body)
		if e != nil {
			panic(e)
		}
		br = bytes.NewReader(b)
	}

	req := httptest.NewRequest(echo.GET, url, br)
	rec := httptest.NewRecorder()

	for hn, hv := range c.Header {
		req.Header[hn] = hv
	}

	if headers != nil {
		for hn, hv := range headers {
			req.Header[hn] = hv
		}
	}

	c.e.ServeHTTP(rec, req)

	return Response{
		Code:   rec.Code,
		Body:   rec.Body.Bytes(),
		Header: req.Header,
	}
}

// Get send GET request.
func (c *Client) Get(url string, header http.Header) Response {
	return c.request(echo.GET, url, header, nil)
}

// Post send POST request.
func (c *Client) Post(url string, header http.Header, body interface{}) Response {
	return c.request(echo.POST, url, header, body)
}

// Put send PUT request.
func (c *Client) Put(url string, header http.Header, body interface{}) Response {
	return c.request(echo.PUT, url, header, body)
}

// Patch send PATCH request.
func (c *Client) Patch(url string, header http.Header, body interface{}) Response {
	return c.request(echo.PATCH, url, header, body)
}

// Delete send DELETE request.
func (c *Client) Delete(url string, header http.Header, body interface{}) Response {
	return c.request(echo.DELETE, url, header, body)
}
