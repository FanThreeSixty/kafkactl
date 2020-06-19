package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type topicStruct struct {
	Name       string            `json:"name"`
	Partitions []partitionStruct `json:"partitions"`
}

type partitionStruct struct {
	ID         int32        `json:"id"`
	LeaderAddr string       `json:"leader_addr"`
	LeaderID   int32        `json:"leader_id"`
	Offsets    offsetStruct `json:"offsets"`
}

type offsetStruct struct {
	Group  int64 `json:"group"`
	LogEnd int64 `json:"log_end"`
	Lag    int64 `json:"lag"`
}

// move this to getPartitionsByTopic -> partitions.go
func getPartitionsByTopic(kafkaConnection sarama.Client, topic string) ([]int32, error) {
	var returnPartitions []int32
	partitions, err := kafkaConnection.Partitions(topic)
	if err != nil {
		errorMessage := fmt.Sprintf("[ERROR] Unable to find Partitions for Topic %s (err: %s).", topic, err)
		log.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}
	for _, partition := range partitions {
		returnPartitions = append(returnPartitions, partition)
	}

	return returnPartitions, nil
}

func getTopics(cmd *cobra.Command, args []string) {
	var err error

	kafkactlCfg.clientConn, err = sarama.NewClient(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err)
	}
	defer kafkactlCfg.clientConn.Close()

	table := uitable.New()
	table.Wrap = true
	table.AddRow("TOPIC", "NUM OF PARTITIONS")

	topics := retrieveTopics(kafkactlCfg.clientConn)
	topicChan := make(chan topicStruct)
	defer close(topicChan)

	for _, topic := range topics {
		go func(top topicStruct, conn sarama.Client) {
			top.retrievePartitions(conn)
			topicChan <- top
		}(topic, kafkactlCfg.clientConn)
	}

	for i := 0; i < len(topics); i++ {
		t := <-topicChan
		table.AddRow(t.Name, len(t.Partitions))
	}

	fmt.Println(table)

}

func describeTopic(cmd *cobra.Command, args []string) {
	var err error
	if args[0] == "" {
		fmt.Println("[ERROR] Missing Topic to Describe.\n\tUsage: kafkactl topic describe example")
		return
	}

	kafkactlCfg.clientConn, err = sarama.NewClient(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err)
		return
	}
	defer kafkactlCfg.clientConn.Close()

	topic := topicStruct{
		Name: args[0],
	}
	topic.retrievePartitions(kafkactlCfg.clientConn)

	brokerMap := make(map[int]string)
	table := uitable.New()
	table.Wrap = true

	partitionChan := make(chan partitionStruct, len(topic.Partitions))
	defer close(partitionChan)

	for _, partition := range topic.Partitions {
		go func(top string, part int32, conn sarama.Client) {
			offsetNewest, err := conn.GetOffset(top, part, sarama.OffsetNewest)
			if err != nil {
				log.Println(err)
			}
			leaderOfPartition, err := conn.Leader(top, part)
			if err != nil {
				log.Println(err)
			}

			partitionChan <- partitionStruct{
				ID:         part,
				LeaderAddr: leaderOfPartition.Addr(),
				LeaderID:   leaderOfPartition.ID(),
				Offsets: offsetStruct{
					LogEnd: offsetNewest,
				},
			}
		}(topic.Name, partition.ID, kafkactlCfg.clientConn)
	}

	// TODO:
	// need to add ability to sort through this array based on partition ID
	var parts []partitionStruct
	for i := 0; i < len(topic.Partitions); i++ {
		parts = append(parts, <-partitionChan)
	}

	table.AddRow(topic.Name, "PARTITION", "LEADER", "NEWEST OFFSET")
	for _, part := range parts {
		brokerMap[int(part.LeaderID)] = part.LeaderAddr
		table.AddRow("", part.ID, fmt.Sprintf("ID: %v", part.LeaderID), part.Offsets.LogEnd)
	}
	if kafkactlCfg.expand {
		table.AddRow("--LEGEND--")
		table.AddRow("BrokerID", "HostAddr")
		for brokerID, brokerAddr := range brokerMap {
			table.AddRow(brokerID, brokerAddr)
		}
		table.AddRow("----------")
	}
	fmt.Println(table)
}

func retrieveTopics(clientConn sarama.Client) (topicStructs []topicStruct) {
	topics, _ := clientConn.Topics()
	for _, topic := range topics {
		if !kafkactlCfg.verbose && topic == "__consumer_offsets" {
			continue
		}
		topicStructs = append(topicStructs, topicStruct{
			Name: topic,
		})
	}
	return
}

func (topic *topicStruct) retrievePartitions(clientConn sarama.Client) {
	partitions, err := clientConn.Partitions(topic.Name)
	if err != nil {
		log.Printf("[WARNING] Unable to get Partitions for topic %s (err: %s", topic.Name, err)
	}
	for _, partition := range partitions {
		leaderOfPartition, err := clientConn.Leader(topic.Name, partition)
		if err != nil {
			log.Println(err)
		}
		topic.Partitions = append(topic.Partitions, partitionStruct{
			ID:         partition,
			LeaderAddr: leaderOfPartition.Addr(),
			LeaderID:   leaderOfPartition.ID(),
		})
	}
	return
}

func (partition *partitionStruct) retrieveOffsets(clientConn sarama.Client, partitionTopic string) (oStruct []offsetStruct) {
	return
}

func collectTopics() (tStructs []topicStruct) {
	clientConn, err := sarama.NewClient(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err)
	}
	defer clientConn.Close()

	topics, _ := clientConn.Topics()
	for _, topic := range topics {
		tStruct := topicStruct{Name: topic}
		partitions, _ := clientConn.Partitions(topic)
		for _, partition := range partitions {
			tStruct.Partitions = append(tStruct.Partitions, partitionStruct{
				ID: partition,
			})
		}
		tStructs = append(tStructs, tStruct)
	}
	return
}
