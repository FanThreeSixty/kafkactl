package commands

import "github.com/spf13/cobra"

var (
	describeCmd = &cobra.Command{
		Use:     "describe",
		Aliases: []string{"desc"},
		Short:   "command for describing resources in Kafka.",
	}

	describeGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "describe group in a particular Kafka cluster.",
		Long:  `describe an existing group in a particular Kafka cluster.`,
		Run:   describeGroup,
	}

	describeTopicsCmd = &cobra.Command{
		Use:   "topic",
		Short: "describe topics available in a particular Kafka cluster.",
		Long:  `retrieve a describe of topics that are available in a particular Kafka cluster.`,
		Run:   describeTopic,
	}
)

func init() {
	//	describeCmd.AddCommand(describeGroupsCmd)
	describeCmd.AddCommand(describeTopicsCmd)
	describeGroupCmd.PersistentFlags().StringVarP(&kafkactlArgs.topic, "topic", "t", "", "Topic(s) to Search For (space delimited)")
	describeCmd.AddCommand(describeGroupCmd)
}
