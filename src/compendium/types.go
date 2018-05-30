package compendium

type Compendium struct {
	Server ServerSpec
	Tomes  []TomeSpec
}

type ServerSpec struct {
	Address  string
	RoleID   string
	SecretID string
}

type TomeSpec struct {
	Instrument map[string]interface{}
	Secrets    []SecretSpec
}

type SecretSpec struct {
	Path   string
	Fields []SecretField

	WatchForChanges bool
	WatchInterval   int32
}

type SecretField struct {
	Name         string
	MapTo        string
	DefaultValue string
	Omit         bool

	ForceUpper bool
	ForceLower bool
}
