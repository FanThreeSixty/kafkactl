package commands

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

func rewindOffset(cmd *cobra.Command, args []string) {
	table := uitable.New()
	table.Wrap = true

	clientConn, err := sarama.NewClient(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err)
	}

	myPartitions, err := getPartitionsByTopic(clientConn, kafkactlArgs.topic)
	if err != nil {
		return
	}
	defer clientConn.Close()

	om, err := sarama.NewOffsetManagerFromClient(kafkactlArgs.consumerGroup, clientConn)
	if err != nil {
		log.Fatalf("[ERROR] Unable to Manage offsets (err: %s)", err)
	}
	defer om.Close()

	for _, partition := range myPartitions {
		pom, err := om.ManagePartition(kafkactlArgs.topic, partition)
		if err != nil {
			log.Fatalf("[ERROR] Unable to connect to offset Manager for Partition %v", partition)
		}
		defer pom.Close()

		offset, err := clientConn.GetOffset(kafkactlArgs.topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Fatalf("[ERROR] Unable to retrieve offset for topic/partition: %s/%v, err: %s", kafkactlArgs.topic, partition, err)
		}
		log.Printf("rewinding offset for topic/partition to oldest known offset: %s/%v -> %v", kafkactlArgs.topic, partition, offset)
		pom.MarkOffset(offset, "")
	}

	return

}

// fastForwardConsumerGroupOffset takes the given group and manages all of the resulting partitions of
// a given topic to their latest offsets to 'catch up' the topic
func fastForwardOffset(cmd *cobra.Command, args []string) {
	// TODO: verify topic exists

	clientConfig := sarama.NewConfig()
	clientConfig.ClientID = "kafkactl"
	table := uitable.New()
	table.Wrap = true

	clientConn, err := sarama.NewClient(kafkactlCfg.brokers, clientConfig)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err)
	}

	myPartitions, err := getPartitionsByTopic(clientConn, kafkactlArgs.topic)
	if err != nil {
		return
	}
	defer clientConn.Close()

	om, err := sarama.NewOffsetManagerFromClient(kafkactlArgs.consumerGroup, clientConn)
	if err != nil {
		log.Fatalf("[ERROR] Unable to Manage offsets (err: %s)", err)
	}
	defer om.Close()

	for _, partition := range myPartitions {
		pom, err := om.ManagePartition(kafkactlArgs.topic, partition)
		if err != nil {
			log.Fatalf("[ERROR] Unable to connect to offset Manager for Partition %v", partition)
		}
		defer pom.Close()

		newestOffset, err := clientConn.GetOffset(kafkactlArgs.topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("[ERROR] Unable to retrieve offset for topic/partition: %s/%v, err: %s", kafkactlArgs.topic, partition, err)
		}
		log.Printf("fast forwarding offset for topic/partition: %s/%v -> %v", kafkactlArgs.topic, partition, newestOffset)
		pom.MarkOffset(newestOffset, "")
	}

	return
}
