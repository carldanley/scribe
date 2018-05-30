package vault

import (
	"log"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

func (v *Vault) GetClient() *vault.Client {
	currentTime := int32(time.Now().Unix())

	// make sure the vault auth pointer exists
	if v.VaultAuth == nil {
		v.VaultAuth = &VaultAuth{}
	}

	// check that the vault client exists
	if v.VaultAuth.Client == nil {
		client, err := vault.NewClient(vault.DefaultConfig())
		if err != nil {
			log.Println("Could not create vault client:", err)
			os.Exit(1)
		}

		if v.Address == "" {
			v.Address = "http://127.0.0.1:8200"
		}

		client.SetAddress(v.Address)
		v.VaultAuth.Client = client
	}

	// check to see if the secret auth has expired or not
	if v.VaultAuth.SecretAuth != nil && (currentTime-v.VaultAuth.TimeAcquired) >= int32(v.VaultAuth.SecretAuth.LeaseDuration) {
		v.VaultAuth.SecretAuth = nil
	}

	// check to see if the secret auth has been acquired
	if v.VaultAuth.SecretAuth == nil {
		secret, err := v.VaultAuth.Client.Logical().Write("auth/approle/login", map[string]interface{}{
			"role_id":   v.RoleID,
			"secret_id": v.SecretID,
		})

		if err != nil {
			log.Println("Could not login with specified app role:", err)
			os.Exit(1)
		}

		v.VaultAuth.SecretAuth = secret.Auth
		v.VaultAuth.TimeAcquired = currentTime
		v.VaultAuth.Client.SetToken(secret.Auth.ClientToken)
	}

	return v.VaultAuth.Client
}
