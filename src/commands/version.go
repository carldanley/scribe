package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run:   runVersionCommand,
	}

	RootCmd.AddCommand(versionCmd)
}

func runVersionCommand(md *cobra.Command, args []string) {
	fmt.Println("v0.3.10")
}
