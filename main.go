package main

import (
	"log"

	"github.com/sporting-innovations/kafkactl/commands"
)

func main() {
	if err := commands.KafkaCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
