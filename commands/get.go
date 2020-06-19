package commands

import "github.com/spf13/cobra"

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "command for getting resources in Kafka.",
	}

	getGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "get consumer groups in a particular Kafka cluster.",
		Long:  `retrieve a get of consumer groups in a particular Kafka cluster.`,
		Run:   getGroups,
	}

	getTopicsCmd = &cobra.Command{
		Use:   "topics",
		Short: "get topics available in a particular Kafka cluster.",
		Long:  `retrieve a get of topics that are available in a particular Kafka cluster.`,
		Run:   getTopics,
	}
)

func init() {
	getCmd.AddCommand(getGroupsCmd)
	getCmd.AddCommand(getTopicsCmd)
}
