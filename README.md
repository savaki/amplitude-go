amplitude-go
============

Amplitude client for Go. For additional documentation visit https://amplitude.com/docs or view the godocs.

## Installation
---

	$ go get github.com/savaki/amplitude-go

##Examples
---

### Basic Client

Full example of a simple event tracker.

```go
package main

import (
	"github.com/savaki/amplitude-go"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	client := amplitude.DefaultClient(apiKey)

	// send your event to amplitude
	client.Event(map[string]interface{}{
		"user_id":    "123",
		"event_type": "abc",
	})

	client.Close()
}
```

### Customized Client

Example with custom options.

```go
package main

import (
	"github.com/savaki/amplitude-go"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	
	// allow 1024 messages to be buffered
	// use 10 concurrent goroutines to send messages
	client := amplitude.New(apiKey, 1024).Workers(10)
	
	// send your event to amplitude
	client.Event(map[string]interface{}{
		"user_id":    "123",
		"event_type": "abc",
	})
	
	client.Close()
}
```

### Flushing on Shutdown

The call to `client.Close()` will flush and wait for pending calls to be sent to Amplitude.

### Debugging

Enable debug output via the `DEBUG` environment variables:

```
export DEBUG="*"
```

The value of `DEBUG` provides a pattern match to allow you to selectively display certain lines.

```
22:27:39.302 291us  289us  amplitude - Client.DefaultClient
22:27:39.302 1us    1us    amplitude - Client.New
22:27:39.302 642us  642us  amplitude - Client.Workers(1) => current: 0
22:27:39.303 230us  230us  amplitude - Worker.Started
22:27:39.303 7us    7us    amplitude - Client.Workers(2) => current: 1
22:27:39.303 4us    4us    amplitude - Worker.Started
22:27:39.303 7us    7us    amplitude - Client.Event
22:27:39.303 66us   66us   amplitude - Worker.Flush
22:27:39.303 91us   91us   amplitude - Worker.Closed
22:27:40.844 1s     1s     amplitude - Worker.Closed
22:27:40.844 3us    3us    amplitude - Client.Close
```



