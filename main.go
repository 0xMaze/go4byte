package main

import (
	"fbyte/cli"
	"log"
)

func main() {
	err := cli.Execute()
	if err != nil {
		log.Fatalf("Command could not be executed: %v", err)
	}
}
