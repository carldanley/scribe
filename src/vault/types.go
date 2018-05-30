package vault

import (
	"github.com/carldanley/scribe/src/compendium"
	"github.com/carldanley/scribe/src/instruments"
	vault "github.com/hashicorp/vault/api"
)

type Vault struct {
	Address               string
	RoleID                string
	SecretID              string
	ShouldWatchForChanges bool

	Compendium *compendium.Compendium
	VaultAuth  *VaultAuth

	TomeCache   map[*compendium.TomeSpec]*Tome
	SecretCache map[string]*map[string]string
}

type VaultAuth struct {
	SecretAuth   *vault.SecretAuth
	TimeAcquired int32
	Client       *vault.Client
}

type Tome struct {
	Spec       *compendium.TomeSpec
	Instrument *instruments.Instrument
	Secrets    map[string]*Secret
}

type Secret struct {
	Spec           *compendium.SecretSpec
	Cache          map[string]string
	LastCacheCheck int32
}
