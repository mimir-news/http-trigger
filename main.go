package main

import (
	"log"
)

func main() {
	opts := getOptions()
	c := newClient(opts)

	err := c.trigger()
	if err != nil {
		log.Fatal("Trigger failed:", err)
	}

	log.Println("Trigger OK")
}
