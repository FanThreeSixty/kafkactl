package commands

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// this all seems super messy...
type consumerGroupStruct struct {
	Name            string                      `json:"name"`
	Members         []consumerGroupMemberStruct `json:"members"`
	CoordinatorID   int32                       `json:"coordinator_id"`
	CoordinatorAddr string                      `json:"coordinator_addr"`
	State           string                      `json:"state"`
	Topics          []topicStruct               `json:"topics"`
}

type consumerGroupMemberStruct struct {
	ID        string `json:"id"`
	Topic     string `json:"consuming_topic"`
	Partition int32  `json:"consumer_partition"`
}

type consumerGroupReturn struct {
	Status   string                         `json:"status"`
	Message  string                         `json:"message"`
	Metadata map[string]consumerGroupStruct `json:"metadata"`
}

func (cg *consumerGroupStruct) findConsumerID(topic string, partition int32) consumerGroupMemberStruct {
	for _, member := range cg.Members {
		if member.Topic == topic {
			if member.Partition == partition {
				return member
			}
		}
	}
	return consumerGroupMemberStruct{}
}

func getGroups(cmd *cobra.Command, args []string) {
	var getGroupsReturn consumerGroupReturn
	getGroupsReturn.Metadata = getConsumerGroups()

	table := uitable.New()
	table.Wrap = true
	table.AddRow("GROUP", "MEMBERS", "COORDINATOR.ID", "COORDINATOR.ADDR")

	for groupName, groupMeta := range getGroupsReturn.Metadata {
		table.AddRow(groupName, len(groupMeta.Members), groupMeta.CoordinatorID, groupMeta.CoordinatorAddr)
	}
	fmt.Println(table)

}

func describeGroup(cmd *cobra.Command, args []string) {
	consumerGroups := getConsumerGroups()
	//TODO: handle this error better
	if len(args) == 0 {
		log.Printf("[ERROR] Missing Group to search for.  Please use the following command with an available Kafka Consumer Group:")
		log.Printf("\tkafkactl -e %s desc group $GROUP", kafkactlCfg.environment)
		for groupName := range consumerGroups {
			log.Printf("\t\t- %s", groupName)
		}
		return
	}

	myGroup := consumerGroups[args[0]]

	table := uitable.New()
	table.Wrap = true
	table.AddRow("-----")
	table.AddRow("Group:", myGroup.Name)
	table.AddRow("Coordinator:")
	table.AddRow("\tAddr:", myGroup.CoordinatorAddr)
	table.AddRow("\tID:", myGroup.CoordinatorID)

	// TODO: GoRoutine this -> needs connection pooling
	if kafkactlArgs.topic != "" {
		//TODO: Verify topic exists first, don't just create one
		myGroup.getConsumerGroupTopicDetails(kafkactlArgs.topic)
		for _, topicMeta := range myGroup.Topics {
			table.AddRow("---" + topicMeta.Name + "---")
			table.AddRow("PARTITION", "CONSUMER", "OFFSET", "LOG END OFFSET", "LAG")
			for _, partition := range topicMeta.Partitions {
				consumer := myGroup.findConsumerID(topicMeta.Name, partition.ID)
				if consumer.ID != "" {
					table.AddRow(partition.ID, consumer.ID, partition.Offsets.LogEnd, partition.Offsets.Group, partition.Offsets.Lag)
				} else {
					table.AddRow(partition.ID, "-", partition.Offsets.LogEnd, partition.Offsets.Group, partition.Offsets.Lag)
				}
			}
		}
	} else {
		table.AddRow("---ACTIVE MEMBERS---")
		table.AddRow("ID", "TOPIC", "PARTITION")
		for _, memberMeta := range myGroup.Members {
			table.AddRow(memberMeta.ID, memberMeta.Topic, memberMeta.Partition)
		}
	}

	//TODO: Support JSON formatting for Output
	fmt.Println(table)
}

func getConsumerGroups() (consumerGroups map[string]consumerGroupStruct) {
	consumerGroups = make(map[string]consumerGroupStruct)
	for _, broker := range kafkactlCfg.brokers {
		clientConn := sarama.NewBroker(broker)

		gRequest := sarama.ListGroupsRequest{}
		err := clientConn.Open(kafkactlCfg.client)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			return nil
		}
		defer clientConn.Close()

		groups, err := clientConn.ListGroups(&gRequest)
		if err != nil {
			log.Fatalf("[ERROR] Cannot list groups (err: %s).", err)
			return nil
		}

		for group := range groups.Groups {
			returningConsumerGroup := consumerGroupStruct{
				Name: group,
			}
			returningConsumerGroup.getConsumerGroupDetails(clientConn)
			consumerGroups[group] = returningConsumerGroup
		}

	}
	return
}

func (cg *consumerGroupStruct) getConsumerGroupDetails(broker *sarama.Broker) error {
	groupDescs, _ := broker.DescribeGroups(&sarama.DescribeGroupsRequest{Groups: []string{cg.Name}})
	var members []string
	for member, memberMeta := range groupDescs.Groups[0].Members {
		members = append(members, member)
		assignmentMeta, _ := memberMeta.GetMemberAssignment()
		for topic, partition := range assignmentMeta.Topics {
			for _, p := range partition {
				cg.Members = append(cg.Members, consumerGroupMemberStruct{
					ID:        member,
					Topic:     topic,
					Partition: p,
				})
			}
			if kafkactlCfg.verbose {
				log.Printf("%s assigned to: %s-%v (%s)", member, topic, partition, assignmentMeta.UserData)
			}
		}
	}
	groupMetadata, err := broker.GetConsumerMetadata(&sarama.ConsumerMetadataRequest{ConsumerGroup: cg.Name})
	if err != nil {
		return displayError(fmt.Sprintf("[ERROR] Unable to retreive Metadata for Group %s", cg.Name))
	}
	defer groupMetadata.Coordinator.Close()

	cg.CoordinatorID = groupMetadata.CoordinatorID
	cg.CoordinatorAddr = groupMetadata.CoordinatorHost
	return nil
}

func (cg *consumerGroupStruct) getConsumerGroupTopicDetails(lfTopic string) error {
	clientConn, err := sarama.NewClient(kafkactlCfg.brokers, kafkactlCfg.client)
	if err != nil {
		return displayError(fmt.Sprintf("[ERROR] Unable to connect to Kafka Brokers (err: %s)", err))
	}
	defer clientConn.Close()

	om, err := sarama.NewOffsetManagerFromClient(cg.Name, clientConn)
	if err != nil {
		return displayError(fmt.Sprintf("[ERROR] Unable to Manage Offsets (err: %s)", err))
	}

	myPartitions, err := getPartitionsByTopic(clientConn, lfTopic)
	if err != nil {
		return displayError(fmt.Sprintf("[ERROR] Unable to get Partitions for Topic %s (err: %s", lfTopic, err))
	}

	cgTopic := topicStruct{Name: lfTopic}

	for _, partition := range myPartitions {
		leaderOfPartition, err := clientConn.Leader(lfTopic, partition)
		if err != nil {
			log.Printf("[WARNING] Unable to find Leader for Partition %v (err: %s)", partition, err)
		}
		cgPartition := partitionStruct{
			ID:         partition,
			LeaderAddr: leaderOfPartition.Addr(),
			LeaderID:   leaderOfPartition.ID(),
		}

		pom, err := om.ManagePartition(cgTopic.Name, partition)
		if err != nil {
			return displayError(fmt.Sprintf("[ERROR] Unable to connect to Offset Manager for Partition %v", partition))
		}
		defer pom.Close()

		groupOffset, _ := pom.NextOffset()
		if groupOffset == sarama.OffsetOldest {
			return displayError(fmt.Sprintf("offset is oldest (%v is %v)", groupOffset, sarama.OffsetOldest))
		}
		currentOffset, err := clientConn.GetOffset(cgTopic.Name, partition, sarama.OffsetNewest)
		if err != nil {
			return displayError(fmt.Sprintf("[ERROR] Unable to get Current Offset for topic %s:%v", cgTopic.Name, partition))
		}
		if groupOffset == sarama.OffsetNewest {
			groupOffset = currentOffset
		}
		cgPartition.Offsets = offsetStruct{
			LogEnd: currentOffset,
			Group:  groupOffset,
			Lag:    currentOffset - groupOffset,
		}
		cgTopic.Partitions = append(cgTopic.Partitions, cgPartition)
	}
	cg.Topics = append(cg.Topics, cgTopic)

	return nil
}
