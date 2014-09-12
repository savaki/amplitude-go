package main

import (
	"github.com/savaki/amplitude-go"
	"log"
	"os"
)

func ok(err error) {
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {
	apiKey := os.Getenv("API_KEY")
	client := amplitude.DefaultClient(apiKey).Workers(2)
	err := client.Event(map[string]interface{}{
		"user_id":    "123",
		"event_type": "abc",
	})
	ok(err)
	client.Close()
}
