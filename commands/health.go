package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var (
	healthCmd = &cobra.Command{
		Use:     "health",
		Aliases: []string{"ping"},
		Short:   "Check health of configured Kafka Cluster",
		Long:    `Run connection checks via Sarama and HTTP calls on the configured Kafka Cluster.`,
		Run:     healthCheck,
	}
)

type endpointHealth struct {
	Endpoint     string `json:"endpoint"`
	EndpointAddr string `json:"endpoint_addr"`
	Healthy      bool   `json:"healthy"`
	Message      string `json:"message"`
}

type health struct {
	Brokers []endpointHealth `json:"brokers"`
	Healthy bool             `json:"healthy"`
}

func healthCheck(cmd *cobra.Command, args []string) {
	returnHealth := health{
		Healthy: true,
	}
	for _, broker := range kafkactlCfg.brokers {
		conf := sarama.NewConfig()
		conf.Version = kafkactlCfg.version
		clientConn := sarama.NewBroker(broker)

		err := clientConn.Open(conf)
		if err != nil {
			log.Fatalf("[ERROR] Unable to 'Open' a connection using Sarama: %s", err)
		}
		if connected, err := clientConn.Connected(); !connected {
			returnHealth.Brokers = append(returnHealth.Brokers, endpointHealth{
				Endpoint:     broker,
				EndpointAddr: "N/A",
				Healthy:      false,
				Message:      fmt.Sprintf("%s", err),
			})
			returnHealth.Healthy = false
		} else {
			returnHealth.Brokers = append(returnHealth.Brokers, endpointHealth{
				Endpoint:     broker,
				EndpointAddr: clientConn.Addr(),
				Healthy:      true,
				Message:      "",
			})
			clientConn.Close()
		}
	}

	json, _ := json.Marshal(returnHealth)
	fmt.Println(string(json))
}
