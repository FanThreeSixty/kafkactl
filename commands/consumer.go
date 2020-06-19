package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"

	kafkaCluster "github.com/bsm/sarama-cluster"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "consume commands for Kafka",
	Long:  "commands to consume content inside of the Kafka Node/Cluster",
	Run:   consume,
}

func init() {
	consumeCmd.PersistentFlags().StringVarP(&kafkactlArgs.consumerGroup, "group", "g", "", "Utilize a Consumer Group for Consumption of a Particular Topic.")
	consumeCmd.PersistentFlags().Int32VarP(&kafkactlArgs.partition, "partition", "p", 0, "Partition to Consume From (defaults to 0).")
	consumeCmd.PersistentFlags().Int64VarP(&kafkactlArgs.offset, "offset", "", 0, "Offset to begin Consumption From (defaults to oldest).")

}

func consume(cmd *cobra.Command, args []string) {
	if args[0] != "" {
		kafkactlArgs.topic = args[0]
	} else {
		log.Fatalf("[ERROR] Please pass in a topic for consumption\n\tkafkactl consume topic <topic>\n")
		return
	}
	log.Printf("[DEBUG] consumerGroup: %s, partition: %v", kafkactlArgs.consumerGroup, kafkactlArgs.partition)
	if kafkactlArgs.consumerGroup == "" {
		if kafkactlArgs.partition == 0 {
			log.Printf("[DEBUG] Starting to consume from Partition 0. Use [-p|--partition] to change or use [-g|--group] to join a Consumer Group.")
		}
		consumeByPartition()
	} else {
		consumeByGroup()
	}
}

func consumeByPartition() {
	master, err := sarama.NewConsumer(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to brokers as Consumer (err: %s)", err)
	}
	defer master.Close()

	consumer, err := master.ConsumePartition(kafkactlArgs.topic, kafkactlArgs.partition, kafkactlArgs.offset)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect Individual Consumer to Topic %s (err: %s)", kafkactlArgs.topic, err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-consumer.Errors():
			log.Printf("[ConsumerErr] %s", err)
		case msg := <-consumer.Messages():
			fmt.Printf("[%s:p%d:o%d][@%v] %s %s\n", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp, msg.Key, msg.Value)
		case <-signals:
			log.Println("Interrupt is detected")
			return
		}
	}
}

func consumeByGroup() {
	consumer, err := kafkaCluster.NewConsumer(kafkactlCfg.brokers, kafkactlArgs.consumerGroup, []string{kafkactlArgs.topic}, kafkactlCfg.clusterClient)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect Consumer Group Consumer for Group %s to Topic %s (err: %s)", kafkactlArgs.consumerGroup, kafkactlArgs.topic, err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case err, more := <-consumer.Errors():
			if more {
				log.Printf("Error: %s\n", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		case <-signals:
			return
		}
	}
}
