package amplitude

import (
	"encoding/json"
	"fmt"
	. "github.com/visionmedia/go-debug"
	"net/http"
	"net/url"
	"sync"
)

const (
	DefaultQueueSize = 100
	ApiEndpoint      = "https://api.amplitude.com/httpapi"
)

var debug = Debug("amplitude")

type event map[string]interface{}

type Client struct {
	ApiKey  string
	wg      *sync.WaitGroup
	events  chan event
	workers []worker
}

func DefaultClient(apiKey string) *Client {
	debug("Client.DefaultClient")
	return New(apiKey, DefaultQueueSize)
}

func New(apiKey string, queueSize int) *Client {
	debug("Client.New")
	events := make(chan event, queueSize)
	workers := []worker{}

	client := &Client{
		ApiKey:  apiKey,
		events:  events,
		workers: workers,
		wg:      &sync.WaitGroup{},
	}

	return client.Workers(1)
}

func (c *Client) Event(e event) error {
	if _, ok := e["user_id"]; !ok {
		err := fmt.Errorf("missing required parameter: user_id")
		debug(err.Error())
		return err
	}
	if _, ok := e["event_type"]; !ok {
		err := fmt.Errorf("missing required parameter: event_type")
		debug(err.Error())
		return err
	}

	select {
	case c.events <- e:
		debug("Client.Event")
	default:
		err := fmt.Errorf("Unable to send event, queue is full.  Use a larger queue size or create more workers.")
		debug(err.Error())
		return err
	}

	return nil
}

func (c *Client) Workers(desired int) *Client {
	debug(fmt.Sprintf("Client.Workers(%d) => current: %d", desired, len(c.workers)))
	if current := len(c.workers); desired > current {
		started := make(chan bool)
		for i := current; i < desired; i++ {
			w := worker{
				apiKey: c.ApiKey,
				events: c.events,
				wg:     c.wg,
			}
			go w.start(started)
			<-started

			c.workers = append(c.workers, w)
		}
		close(started)

	} else if desired < current {
		for i := current; i > desired; i-- {
			c.workers[0].Close()
			c.workers = c.workers[1:]
		}

	}

	return c
}

func (c *Client) Close() {
	close(c.events)
	c.wg.Wait()
	debug("Client.Close")
}

type worker struct {
	done   chan struct{}
	wg     *sync.WaitGroup
	apiKey string
	events <-chan event
}

func (w *worker) start(started chan bool) {
	debug("Worker.Started")
	defer debug("Worker.Closed")

	// maintain worker reference counter
	w.wg.Add(1)
	defer w.wg.Done()

	started <- true

	// done is our control channel to allow this worker to be independently closed
	if w.done == nil {
		w.done = make(chan struct{})
	}

	for {
		select {
		case e, open := <-w.events:
			if e != nil {
				w.flush(e)
			} else if !open {
				return // channel closed
			}

		case <-w.done:
			return
		}
	}
}

func (w *worker) Close() {
	debug("Worker.Closing")
	close(w.done)
}

func (w *worker) flush(e event) {
	data, err := json.Marshal(e)
	if err != nil {
		debug(err.Error())
		return
	}

	params := url.Values{}
	params.Set("api_key", w.apiKey)
	params.Set("event", string(data))

	debug("Worker.Flush")
	_, err = http.PostForm(ApiEndpoint, params)
	if err != nil {
		debug(err.Error())
	}
}
