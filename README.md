![scribe](./docs/images/scribe-logo-small.png)

# What is scribe?

`scribe` is a command-line tool that uses a compendium to research scripts from Vault and transcribe them into "tomes" using different instruments.

# Glossary of Terms

| Term | Definition |
|:-----|:-----------|
| script | A secret from HashiCorp's Vault; contains one or many key/value pairs |
| instrument | A tool used for transcribing secrets and writing them to a target (file, socket, etc) |
| tome | A composition of scripts that were transcribed with different instruments |
| compendium | A configuration file that contains knowledge about how scripts should be bound together with an instrument to create a tome |

# Features

`scribe` boasts a handsome set of features that enable you to manipulate secrets in HashiCorp's Vault with ease. Here are some of the main features that `scribe` offers:

* Can write one or many tomes with a single compendium
* Efficient with queries to Vault to keep TCP traffic as light as possible
* Capable of watching Vault for changes and re-writing the tomes affected by the corresponding change(s)
* Capable of transcribing individual fields contained in a secret; each field's transcription can be controlled individually
* Comes with ready-made, [officially supported instruments](./docs/supported-instruments.md)
* Supports extension of scribe through the use of custom instrumentation

# Additional Documentation

* [Getting Started](./docs/getting-started.md)
* [Running `scribe` in a Docker Container](./docs/docker.md)
* [Officially Supported Instruments](./docs/supported-instruments.md)
* [Writing Custom Instruments](./docs/writing-custom-instruments.md)
