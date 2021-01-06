package main

import (
	"log"

	"github.com/FanThreeSixty/kafkactl/commands"
)

func main() {
	if err := commands.KafkaCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
