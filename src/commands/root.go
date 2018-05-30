package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "scribe",
	Short: "scribe is a secret compositor for HashiCorp's Vault",
}
