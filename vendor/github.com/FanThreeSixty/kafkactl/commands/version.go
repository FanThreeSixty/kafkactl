package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const kafkactlVersion = "0.7.0"

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Printf(version())
}

func version() string {
	return fmt.Sprintf("Version: %s", kafkactlVersion)
}
