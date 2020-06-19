package commands

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "producer commands for Kafka",
	Long:  "commands to produce content inside of the Kafka Node/Cluster",
	Run:   produceContent,
}

var (
	message string
)

func init() {
	produceCmd.PersistentFlags().StringVarP(&kafkactlArgs.topic, "topic", "t", "", "Topic to Produce Message To.")
	produceCmd.PersistentFlags().StringVarP(&message, "message", "m", "", "Message to Produce.")

}

func produceContent(cmd *cobra.Command, args []string) {
	kafkactlCfg.client.Producer.Return.Successes = true
	kafkactlCfg.client.Producer.RequiredAcks = sarama.WaitForAll
	kafkactlCfg.client.Producer.Retry.Max = 5

	if kafkactlArgs.topic == "" {
		log.Fatalf("[ERROR] Unable to write to a blank topic.  Please use '--topic' or '-t' to set a Topic Name.")
	}
	if message == "" {
		log.Fatalf("[ERROR] Unable to write a blank message.  Please use '--message' or '-m' to set a Message.")
	}

	if kafkactlCfg.verbose {
		log.Printf("[VERBOSE] Sending Message: '%s' to Topic %s", message, kafkactlArgs.topic)
	}

	producer, err := sarama.NewSyncProducer(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to broker as Sync Producer (err: %s)", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: kafkactlArgs.topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("-> Produced message '%s' in topic '%s'/partition '%d'/offset '%d'\n", message, kafkactlArgs.topic, partition, offset)
	}

}
