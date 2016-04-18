package main

import (
	"os"

	"github.com/savaki/amplitude-go"
)

func main() {
	apiKey := os.Getenv("AMPLITUDE_API_KEY")
	client := amplitude.New(apiKey)
	client.Publish(amplitude.Event{
		UserId:    "123",
		EventType: "sample",
	})
	client.Flush()
	client.Close()
}
