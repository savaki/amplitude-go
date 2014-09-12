package amplitude

import (
	"time"
	"net/http"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	DefaultQueueSize = 250
	ApiEndpoint      = "https://api.amplitude.com/httpapi"
)

type event map[string]interface{}

type Client struct {
	ApiKey     string
	FlushAt    int
	FlushAfter time.Duration
	events     chan event
	workers    []worker
}

func DefaultClient(apiKey string) *Client {
	return New(apiKey, DefaultQueueSize)
}

func New(apiKey string, queueSize int) *Client {
	events := make(chan event, queueSize)
	workers := []worker{}

	client := &Client{
		ApiKey:  apiKey,
		events:  events,
		workers: workers,
	}

	return client.Workers(1)
}

func (c *Client) Event(e event) error {
	if _, ok := e["user_id"] ; ok {
		return fmt.Errorf("missing required parameter: user_id")
	}
	if _, ok := e["event_type"] ; ok {
		return fmt.Errorf("missing required parameter: event_type")
	}

	select {
	case c.events <- e:
	default:
		return fmt.Errorf("Unable to send event, queue is full.  Use a larger queue size or create more workers.")
	}

	return nil
}

func (c *Client) Workers(desired int) *Client {
	if current := len(c.workers); desired > current {
		for i := current; i < desired; i++ {
			w := worker{
				apiKey: c.ApiKey,
				events: c.events,
			}
			w.start()
			c.workers = append(c.workers, w)
		}

	} else if desired < current {
		for i := current; i > desired; i-- {
			c.workers[0].Close()
			c.workers = c.workers[1:]
		}

	}

	return c
}

func (c *Client) Stop() {
	select {
	case <-c.events:
	}
}

type worker struct {
	done   chan string
	apiKey string
	events <-chan event
}

func (w *worker) start() {
	if w.done == nil {
		w.done = make(chan string)
	}

	for {
		select {
		case e := <-w.events:
			w.flush(e)
		case <-w.done:
			return
		}
	}
}

func (w *worker) Close() {
	w.done <- "yep!"
}

func (w *worker) flush(e event) {
	data, _ := json.Marshal(e)

	params := url.Values{}
	params.Set("api_key", w.apiKey)
	params.Set("event", string(data))

	http.PostForm(ApiEndpoint, params)
}
