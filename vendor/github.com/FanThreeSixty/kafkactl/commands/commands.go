package commands

import (
	"errors"
	"log"

	"github.com/spf13/cobra"
)

// KafkaCmd is the root command for Cobra to be exported
var KafkaCmd = &cobra.Command{
	Use:           "kafkactl",
	Short:         "a command line tool used to interact with Kafka",
	Long:          `kafkactl is a command line script to interact with (a) Kafka client(s)/cluster(s)`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information for Kafkactl",
	Run:   getVersion,
}

func init() {
	globalCmdFlags()

	// kafka utility commands
	KafkaCmd.AddCommand(describeCmd)
	KafkaCmd.AddCommand(getCmd)
	KafkaCmd.AddCommand(manageCmd)
	KafkaCmd.AddCommand(consumeCmd)
	KafkaCmd.AddCommand(produceCmd)
	KafkaCmd.AddCommand(healthCmd)
	KafkaCmd.AddCommand(versionCmd)

	cobra.OnInitialize(initConfig)
}

func displayError(message string) error {
	log.Printf(message)
	return errors.New(message)
}
