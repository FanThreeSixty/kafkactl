package commands

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"

	kafkaGroup "github.com/bsm/sarama-cluster"
)

type kafkactlConfig struct {
	brokers           []string
	clusterClient     *kafkaGroup.Config
	clusterClientConn kafkaGroup.Client
	client            *sarama.Config
	clientConn        sarama.Client
	environment       string
	expand            bool
	port              string
	verbose           bool
	version           sarama.KafkaVersion
}

type kafkactlArguments struct {
	consumerGroup string
	offset        int64
	partition     int32
	topic         string
	topics        []string
}

const (
	clientid = "kafkactl"
)

var (
	kafkactlCfg = kafkactlConfig{
		brokers:     []string{"127.0.0.1:9092"},
		environment: "",
		verbose:     false,
	}

	kafkactlArgs kafkactlArguments
)

func globalCmdFlags() {
	// environment flag to determine IF we need to load a config specific
	// to said environment in yaml
	// ie:
	// config.yaml
	// brokers:
	//   - kafka:9092
	// production:
	//   brokers:
	//     - prod-kafka:9092
	KafkaCmd.PersistentFlags().StringVarP(&kafkactlCfg.environment, "environment", "e", "", "Change configuration based on Environment.")
	// flag override for brokers
	KafkaCmd.PersistentFlags().StringSliceVarP(&kafkactlCfg.brokers, "brokers", "b", kafkactlCfg.brokers, "Space delimited list of Kafka Brokers to connect to.")
	KafkaCmd.PersistentFlags().BoolVarP(&kafkactlCfg.expand, "expand", "", false, "Expand given Output.")
	KafkaCmd.PersistentFlags().BoolVarP(&kafkactlCfg.verbose, "verbose", "v", false, "Enable verbose debugging of script.")
}

func initConfig() {
	kafkactlCfg.client = sarama.NewConfig()
	kafkactlCfg.clusterClient = kafkaGroup.NewConfig()

	// configure viper to look into the following directories for
	// a 'config.yaml' configuration file.
	// 1. $HOME/.kafkactl, 2. /etc/kafkactl, 3. $PWD
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.kafkactl")
	viper.AddConfigPath("/etc/kafkactl")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	// assign FILE configuration values first
	err := viper.ReadInConfig()
	if err != nil {
		if kafkactlCfg.verbose {
			log.Printf("[VERBOSE] No configuration file found. Loading defaults. (err :%s)", err)
		}
	} else {
		if kafkactlCfg.verbose {
			log.Printf("[VERBOSE] Configuration file found.")
		}
		var prefix string
		if kafkactlCfg.environment != "" {
			prefix = kafkactlCfg.environment + "."
		} else {
			prefix = "default."
		}

		if len(viper.GetStringSlice(prefix+"brokers")) != 0 {
			kafkactlCfg.brokers = viper.GetStringSlice(prefix + "brokers")
		}

		switch viper.GetString(prefix + "version") {
		case "V0_10_2_0":
			kafkactlCfg.client.Version = sarama.V0_10_2_0
		default:
			kafkactlCfg.client.Version = sarama.V0_10_2_0
		}
	}

	// assign PersistentFlags (take priority over FILE/ENV)

	kafkactlCfg.configKafkaClients()
}

func (cfg *kafkactlConfig) configKafkaClients() {
	cfg.client.ClientID = clientid
	cfg.clusterClient.ClientID = clientid
	cfg.client.Consumer.Return.Errors = true
	cfg.clusterClient.Consumer.Return.Errors = true
}
