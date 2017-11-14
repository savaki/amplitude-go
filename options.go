package amplitude

import (
	"net/http"
	"time"
)

type Option func(*Client)

func Interval(v time.Duration) Option {
	return func(c *Client) {
		c.interval = v
	}
}

func OnPublishFunc(fn func(status int, err error)) Option {
	return func(c *Client) {
		c.onPublishFunc = fn
	}
}

func HTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
