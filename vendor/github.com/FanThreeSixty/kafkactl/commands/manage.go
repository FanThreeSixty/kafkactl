package commands

import (
	"github.com/spf13/cobra"
)

var (
	manageCmd = &cobra.Command{
		Use:   "manage",
		Short: "command for managing resources in Kafka.",
	}

	manageConsumerGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "consumer group resources",
	}

	manageOffsetsCmd = &cobra.Command{
		Use:   "offsets",
		Short: "Manage Offsets",
	}

	manageOffsetFastForwardCmd = &cobra.Command{
		Use:   "fastforward",
		Short: "fast forward offsets for a given group/topic",
		Long:  `manages all partition offsets of a given topic (-t) to the latest offset for a given consumer group (-g)`,
		Run:   fastForwardOffset,
	}
	manageOffsetRewindCmd = &cobra.Command{
		Use:   "rewind",
		Short: "rewind offsets for a given group/topic",
		Long:  `manages all partition offsets of a given topic (-t) to the latest offset for a given consumer group (-g)`,
		Run:   rewindOffset,
	}

	// move to config or elsewhere
	offsetString string
)

func init() {
	// manage offset cmd
	manageOffsetFastForwardCmd.PersistentFlags().StringVarP(&kafkactlArgs.consumerGroup, "group", "g", "", "Consumer Group to manage offsets For.")
	manageOffsetFastForwardCmd.PersistentFlags().StringVarP(&kafkactlArgs.topic, "topic", "t", "", "Topic to manage offsets For with Consumer Group.")
	manageOffsetsCmd.AddCommand(manageOffsetFastForwardCmd)

	manageOffsetRewindCmd.PersistentFlags().StringVarP(&kafkactlArgs.consumerGroup, "group", "g", "", "Consumer Group to manage offsets For.")
	manageOffsetRewindCmd.PersistentFlags().StringVarP(&kafkactlArgs.topic, "topic", "t", "", "Topic to manage offsets For with Consumer Group.")
	manageOffsetsCmd.AddCommand(manageOffsetRewindCmd)

	manageCmd.AddCommand(manageOffsetsCmd)
}
