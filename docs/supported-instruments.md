# Supported Instruments

## Overview

`scribe` comes with the following instruments (which are available for use inside your compendiums):

1. [file](#file)
1. [cat](#cat)
1. [executable](#executable)
1. [template](#template)

## Instrument Configuration

### file

The `file` instrument handles writing composited secrets to a file on disk.

*Please not that this instrument is not recommended for use in a production setting as storing potentially sensitive data on disk is considered bad practice.*

#### Settings

| Key | Type | Description | Default |
|:----|:-----|:------------|:--------|
| `path` | `string` | A path to the file that scribe will write; can be relative or absolute. If the directories leading up to the path do not exist, scribe will create all of them automatically. | `""` |

#### Example Compendium

```yaml
server:
  address: "https://my.vault.address"
  roleID: "16c0b983-c4f4-69ed-1017-bf720f1919f6"
  secretID: "f6dc09d6-16bc-75ed-ed76-303ccab5af4c"
tomes:
- instrument:
    type: "file"
    path: "/my/file/path/file.env"
  secrets:
  - path: "secrets/some/path"
```

### cat

The `cat` instrument handles writing composited secrets to `stdout`.

#### Settings

This instrument does not have any additional settings.

#### Example Compendium

```yaml
server:
  address: "https://my.vault.address"
  roleID: "16c0b983-c4f4-69ed-1017-bf720f1919f6"
  secretID: "f6dc09d6-16bc-75ed-ed76-303ccab5af4c"
tomes:
- instrument:
    type: "cat"
  secrets:
  - path: "secrets/some/path"
```

### executable

The `executable` instrument instructs `scribe` to run a specified command and inject its runtime with the composited secrets.

#### Settings

| Key | Type | Description | Default |
|:----|:-----|:------------|:--------|
| `command` | `string` | A command that `scribe` will execute when composited secrets have been retrieved. It is a good idea to specify the full path to the binary (as opposed to relying on `$PATH` to resolve where the binary exists). | `""` |
| `cleanEnvironment` | `bool` | Indicates whether or not `scribe` will allow pre-existing environment variables to be copied into the runtime for the specified command | `false` |
| `exitOnCommandError` | `bool` | Indicates whether or not `scribe` will exit if the launched child-process exits due to an error. | `false` |
| `restartOnExit` | `bool` | Indicates whether or not `scribe` will restart the launched child-process when it exits. | `false` |

#### Example Compendium

```yaml
server:
  address: "https://my.vault.address"
  roleID: "16c0b983-c4f4-69ed-1017-bf720f1919f6"
  secretID: "f6dc09d6-16bc-75ed-ed76-303ccab5af4c"
tomes:
- instrument:
    type: "executable"
    command: "/usr/bin/node /usr/src/app/index.js"
    cleanEnvironment: true
    exitOnCommandError: true
    restartOnExit: false
  secrets:
  - path: "secrets/some/path"
```

#### Additional Notes

When a tome in a `scribe` compendium instructs `scribe` to use the `executable` instrument **AND** the *same* tome instructs `scribe` to watch for changes to 1 or more secrets, `scribe` will fire a `SIGINT` signal in the child process. This is useful to periodically rotate secrets, make general application changes, or anything else you can think of.

To achieve rolling updates on a server application, it is possible to configure your application to watch for the `SIGINT` signal and carry out the following:

1. Update the health endpoint on your server application to reflect that this process is no longer healthy; thus, stopping new connections.
1. Finish processing all existing connections; thus, draining all connections to the server application.
1. After the server application has 0 active connections, have your process exit itself.
1. `scribe` will restart the process with the newly updated secrets it fetched.

This simple node.js application demonstrates `scribe`'s ability to restart the process:

```js
// console.log() the environment variables
Object.keys(process.env).forEach((key) => {
  console.log(`${key}: ${process.env[key]}`);
});

process.on('SIGINT', () => {
  console.log("killing...");
  setTimeout(() => {
    // drain all of your connections and do whatever you need to here
    console.log("killed");

    // when finished, exit
    process.exit(1)
  }, 5000);
})

// wait forever and ever and ever
const wait = () => setTimeout(wait, 1000);
wait();
```

With a properly configured `scribe` compendium and some random secrets thrown at it, `scribe` should generate some output similar to this:

```
$ godep go run main.go compose --compendium=$HOME/scribe.yaml
==== [scribe] ==== Executing command: /usr/bin/node /Users/carldanley/test-app.js ...
APP_NAME: "my-app"
APP_TITLE: "app-title"
killing...
killed
==== [scribe] ==== Executable finished with error: exit status 1
==== [scribe] ==== Executing command: /usr/bin/node /Users/carldanley/test-app.js ...
APP_TITLE: "updated-app-title"
APP_NAME: "my-app"
```

### template

The `template` instrument handles passing composited secrets through a template to generate an output file. This instrument uses a combination of [Golang templates](https://golang.org/pkg/text/template/) and [Sprig](http://masterminds.github.io/sprig/) to introduce a useful set of templating functionality.

#### Settings

| Key | Type | Description | Default |
|:----|:-----|:------------|:--------|
| `templatePath` | `string` | A path to the desired template file; can be relative or absolute. The contents of this file will be used to generate the template. | `""` |
| `outputPath` | `string` | A path to the file that scribe will write the generated template to; can be relative or absolute. If the directories leading up to the path do not exist, scribe will create all of them automatically. | `""` |

#### Example Compendium

```yaml
server:
  address: "https://my.vault.address"
  roleID: "16c0b983-c4f4-69ed-1017-bf720f1919f6"
  secretID: "f6dc09d6-16bc-75ed-ed76-303ccab5af4c"
tomes:
- instrument:
    type: "template"
    templatePath: "./example.tmpl"
    outputPath: "/some/path/example-output.yaml"
  secrets:
  - path: "secrets/my-secret"
```

#### Additional Notes

Here are a few examples of using Sprig/Go to template things out.

##### Iterating over Keys

```
{{- range $key, $value := . }}
export {{ $key | upper }}="{{ $value }}"
{{- end }}
```

##### Accessing Keys by Name

```
access_key="{{ .S3_ACCESS_KEY }}"
secret_key="{{ .S3_ACCESS_KEY }}"
```
