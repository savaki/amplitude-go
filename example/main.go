package main

import (
	"fmt"
	"os"

	"github.com/savaki/amplitude-go"
)

func main() {
	apiKey := os.Getenv("AMPLITUDE_API_KEY")
	client := amplitude.New(apiKey, amplitude.OnPublishFunc(func(status int, err error) {
		fmt.Fprintf(os.Stderr, "status: %v, err: %v\n", status, err)
	}))
	client.Publish(amplitude.Event{
		UserId:    "123",
		EventType: "sample",
	})
	client.Flush()
	client.Close()
}
