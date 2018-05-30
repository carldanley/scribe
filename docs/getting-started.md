# Getting Started

## Prerequisites

* A working vault server
* `scribe` command-line tool

## Installation

`scribe` is built in 4 distributions: Linux, Darwin, FreeBSD and Windows. These binaries, along with their checksums, can be found in the "Releases" section of this repository. Download the appropiate binary and install it by copying it to your system's binary path.

Once you have copied the binary into a folder recognized in your system path, you should be able to run `scribe`; after doing so, you'll have an ouput similar to this:

```
scribe is a secret compositor for HashiCorp's Vault

Usage:
  scribe [command]

Available Commands:
  compose     Uses a compendium file (otherwise known as a config file) to compose tomes full of secrets
  help        Help about any command
  version     Print the version number

Flags:
  -h, --help   help for scribe

Use "scribe [command] --help" for more information about a command.
```

Now, that `scribe` has been installed, you need to configure [approle](https://www.vaultproject.io/docs/auth/approle.html) access to the vault server of your choice.

## Enabling Approle Authentication

In order to have `scribe` authenticate with vault server, you need to configure an [approle](https://www.vaultproject.io/docs/auth/approle.html) for it. Here is an example of the commands you can run to setup a basic approle for `scribe`.

```
$ vault auth-enable approle
Successfully enabled 'approle' at 'approle'!

$ nano ./scribe.hcl
path "secrets/*" {
  capabilities = ["list", "read"]
}

$ vault policy-write scribe ./scribe.hcl
Policy 'scribe' written.

$ vault write auth/approle/role/scribe token_ttl=5m token_max_ttl=30m secret_id_num_uses=0 secret_id_ttl=0 token_num_uses=0 policies=default,scribe
Success! Data written to: auth/approle/role/scribe

$ vault read auth/approle/role/scribe/role-id
Key    	Value
---    	-----
role_id	16c0b983-c4f4-69ed-1017-bf720f1919f6

$ vault write -f auth/approle/role/scribe/secret-id
Key               	Value
---               	-----
secret_id         	f6dc09d6-16bc-75ed-ed76-303ccab5af4c
secret_id_accessor	4dcb00ee-55c8-0308-d4b6-90e6aaa023c1
```

## Using a Compendium

`scribe` uses [https://github.com/spf13/viper](https://github.com/spf13/viper) so it is possible to load compendiums in any of the following formats: `JSON`, `TOML`, `YAML`, `HCL` and `Java properties config files`. A compendium is needed to specify how `scribe` will compose tomes of secrets.

Here is a fully-annotated example compendium, written in `YAML`.

```yaml
# The server object contains information about how scribe will connect
# to the vault server and authenticate itself
server:
  # the address of the vault server to connect to
  address: "https://my.vault.address"

  # the roleID that scribe will use to acquire tokens - this can be retrieved
  # by running: `vault read auth/approle/role/scribe/role-id` (depending on
  # how you configured your approle)
  roleID: "16c0b983-c4f4-69ed-1017-bf720f1919f6"

  # the secretID that scribe will use to acquire tokens - this can be retrieved
  # by running: `vault write -f auth/approle/role/scribe/secret-id` (which will
  # generate a new secret each time you run this)
  secretID: "f6dc09d6-16bc-75ed-ed76-303ccab5af4c"

# the tomes that scribe will attempt to compose
tomes:

# each new tome needs to specify an instrument; remember that an instrument is
# a tool that handles sending secret compositions to a target (depending on the type
# of instrument chosen)
- instrument:

    # this specifies the type of instrument that will be used to write the
    # tome of composited secrets
    type: "file"

    # depending on the type, you can specify additional settings (see the
    # supported instruments for further information on additional options)
    # for this case, "path" is a known option for the "file" instrument
    path: "./file-name-a.env"

  # each tome also specifies a list of secrets to composite and the rules for
  # how each of the secrets should be composed
  secrets:

  # for each secret, a path is required. this path represents the destination
  # for where where scribe should pull its secrets from
  - path: "secrets/scribe/general-app-settings"

    # each secret can be watched for changes, if this is specified, scribe will
    # periodically fetch the secret at the specified path and conduct a detection
    # about changes to that path
    watchForChanges: true

    # tells scribe how often to watch vault for changes (in seconds; default: 5)
    watchInterval: 1

    # for each secret, it is possible to specify a list of fields to include or
    # omit (depending on your use case)
    fields:

    # the key this rule refers to
    - name: "AVOID_THIS_KEY"

      # will omit the specified key from the transcription process
      omit: true

  # yet, another path for scribe to research
  - path: "secrets/scribe/globals/aws"

    # fields for the specified path
    fields:

    # the key this rule refers to
    - name: "DEFAULT_REGION"

      # the name will be mapped to this key name (instead of the one that exists)
      mapTo: "AWS_DEFAULT_REGION"
```

The compendium above specifies a tome with 2 secrets and a file instrument. You can instruct scribe to use a compendium like this:

`$ scribe compose --compendium="./my/compendium/config.yaml"`

## Environment Variables

The following environment variables are supported:

| Environment Variable | Description |
|:---------------------|:------------|
| `VAULT_ADDRESS` | The address of the Vault server to connect to |
| `VAULT_ROLE_ID` | The role ID to use when authenticating to the Vault server |
| `VAULT_SECRET_ID` | The secret ID to use when authenticating to the Vault server |
