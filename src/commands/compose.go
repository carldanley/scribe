package commands

import (
	"os"

	"github.com/carldanley/scribe/src/compendium"
	"github.com/carldanley/scribe/src/vault"
	"github.com/spf13/cobra"
)

var compendiumFilePath string
var vaultAddress string
var vaultRoleID string
var vaultSecretID string

func init() {
	command := &cobra.Command{
		Use:   "compose",
		Short: "Uses a compendium file (otherwise known as a config file) to compose tomes full of secrets",
		Run:   runComposeCommand,
	}

	// add flags for this command
	command.Flags().StringVarP(&compendiumFilePath, "compendium", "c", "", "contains knowledge about the composition of tomes and the secrets contained within each of them")
	command.Flags().StringVarP(&vaultAddress, "address", "a", "", "The server address of the Vault server to connect to")
	command.Flags().StringVarP(&vaultRoleID, "role-id", "r", "", "The role ID to use when authenticating with the Vault server")
	command.Flags().StringVarP(&vaultSecretID, "secret-id", "s", "", "The secret ID to use when authenticating with the Vault server")

	// register this command
	RootCmd.AddCommand(command)
}

func runComposeCommand(md *cobra.Command, args []string) {
	vault := vault.Vault{}
	compendium := compendium.GetFromFile(compendiumFilePath)

	// if vault address is explicitly set, honor it
	if vaultAddress != "" {
		vault.Address = vaultAddress
	} else if compendium.Server.Address != "" {
		vault.Address = compendium.Server.Address
	} else {
		vault.Address = os.Getenv("VAULT_ADDRESS")
	}

	// if vault role ID is explicitly set, honor it
	if vaultRoleID != "" {
		vault.RoleID = vaultRoleID
	} else if compendium.Server.Address != "" {
		vault.RoleID = compendium.Server.RoleID
	} else {
		vault.RoleID = os.Getenv("VAULT_ROLE_ID")
	}

	// if vault secret ID is explicitly set, honor it
	if vaultSecretID != "" {
		vault.SecretID = vaultSecretID
	} else if compendium.Server.Address != "" {
		vault.SecretID = compendium.Server.SecretID
	} else {
		vault.SecretID = os.Getenv("VAULT_SECRET_ID")
	}

	vault.RegisterTomesFromCompendium(compendium)
	vault.Update()
}
