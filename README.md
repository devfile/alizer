# Alizer
![Go](https://img.shields.io/badge/Go-1.18-blue)
[![Build status](https://github.com/devfile/alizer/actions/workflows/CI.yml/badge.svg)](https://github.com/devfile/alizer/actions/workflows/CI.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-orange.svg)](./LICENSE)

Alizer (which stands for Application Analyzer) is a utility whose goal is to extract information about an application source code.
Such information are:

- Programming languages.
- Frameworks.
- Tools used to build the application.

Additionally, Alizer can also select one devfile (cloud workspace file) from a list of available devfiles and/or
detect components (the concept of component is taken from Odo and its definition can be read on [odo.dev](https://odo.dev/docs/getting-started/basics/#component)).

## 🚨 BACKPORTED VERSION WARNING

This version of Alizer is backported and uses `go 1.18`. For latest version of alizer please use the [main](https://github.com/devfile/alizer) branch.

## Usage

### CLI

The Go CLI can be built with the below command:

```bash
$ go build alizer.go
```

### CLI Arguments

#### alizer analyze

```shell
./alizer analyze [OPTION]... [PATH]...
```

```sh
  --log {debug|info|warning}    sets the logging level of the CLI. The arg accepts only 3 values [`debug`, `info`, `warning`]. The default value is `warning` and the logging level is `ErrorLevel`.
```

#### alizer component

```shell
./alizer component [OPTION]... [PATH]...
```

```sh
  --log {debug|info|warning}    sets the logging level of the CLI. The arg accepts only 3 values [`debug`, `info`, `warning`]. The default value is `warning` and the logging level is `ErrorLevel`.
  --no-port-detection if this flag exists then no port detection is applied on the given application. If this flag doesn't exist then we are applying port detection as normal. In case we have both --no-port-detection and --port-detection the --no-port-detection overrides everything.
  --port-detection {docker|compose|source}    port detection strategy to use when detecting a port. Currently supported strategies are 'docker', 'compose' and 'source'. You can pass more strategies at the same time. They will be executed in order. By default Alizer will execute docker, compose and source.
```
**Deprecation Warning:** The `--port-detection` flag soon will be deprecated.

#### alizer devfile

```shell
./alizer devfile [OPTION]... [PATH]...
```

```sh
  --log {debug|info|warning}    sets the logging level of the CLI. The arg accepts only 3 values [`debug`, `info`, `warning`]. The default value is `warning` and the logging level is `ErrorLevel`.
  --registry strings    registry where to download the devfiles. Default value: https://registry.devfile.io
  --min-schema-version strings the minimum SchemaVersion of the matched devfile(s). The minimum accepted value is `2.0.0`, otherwise an error is returned.
  --max-schema-version strings the maximum SchemaVersion of the matched devfile(s). The minimum accepted value is `2.0.0`, otherwise an error is returned.
```

### Library Package

#### Language Detection

To analyze your source code with Alizer, just import it and use the recognizer:

```go
import "github.com/devfile/alizer/pkg/apis/recognizer"

languages, err := recognizer.Analyze("your/project/path")
```

#### Component Detection

It detects all components which are found in the source tree where each component consists of:

- _Name_: name of the component
- _Path_: root of the component
- _Languages_: list of languages belonging to the component ordered by their relevance.
- _Ports_: list of ports used by the component

```go
import "github.com/devfile/alizer/pkg/apis/recognizer"

// In case port detection is needed.
components, err := recognizer.DetectComponents("your/project/path")

// If there is no need for port detection
components, err := recognizer.DetectComponentsWithoutPortDetection("your/project/path")
```

For more info about name detection, see the [name detection](docs/public/name_detection.md) doc.

For more info about port detection, see the [port detection](docs/public/port_detection.md) doc.

#### Devfile Detection

It selects a devfile from a list of devfiles (from a devfile registry or other storage) based on the information found in the source tree.

```go
import "github.com/devfile/alizer/pkg/apis/recognizer"
import "github.com/devfile/alizer/pkg/apis/model"

// In case you want specific range of schemaVersion for matched devfiles
devifileFilter := model.DevfileFilter {
	MinSchemaVersion: "<minimum-schema-version>",
	MaxSchemaVersion: "<maximum-schema-version>",
}

// If you don't want a specific range of schemaVersion values
devfileFilter := model.DevfileFilter{}

// Call match devfiles func
devfiles, err := recognizer.MatchDevfiles("myproject", devfiles, devifileFilter)
```

## Outputs

Example of `analyze` command:

```json
[
  {
    "Name": "Go",
    "Aliases": ["golang"],
    "Weight": 94.72,
    "Frameworks": [],
    "Tools": ["1.18"],
    "CanBeComponent": true
  }
]
```

Example of `component` command:

```json
[
  {
    "Name": "spring4mvc-jpa",
    "Path": "path-of-the-component",
    "Languages": [
      {
        "Name": "Java",
        "Aliases": null,
        "Weight": 100,
        "Frameworks": ["Spring"],
        "Tools": ["Maven"],
        "CanBeComponent": true
      }
    ],
    "Ports": null
  }
]
```

Example of `devfile` command:

```json
[
  {
    "Name": "nodejs",
    "Language": "JavaScript",
    "ProjectType": "Node.js",
    "Tags": ["Node.js", "Express", "ubi8"]
  }
]
```

## Contributing

This is an open source project open to anyone. This project welcomes contributions and suggestions!

For information on getting started, refer to the [CONTRIBUTING instructions](CONTRIBUTING.md).

## Release process

An Alizer release is created each time a PR having updates on code is merged. You can create a new release [here](https://github.com/devfile/alizer/releases/new).

- A _tag_ should be created with the version of the release as name. `Alizer` follows the `v{major}.{minor}.{bugfix}` format (e.g `v0.1.0`)
- The _title_ of the release has to be the equal to the new tag created for the release.
- The _description_ of the release is optional. You may add a description if there were outstanding updates in the project, not mentioned in the issues or PRs of this release.

### Release Binaries
For each release a group of binary files is generated. More detailed we have the following types:
- `linux/amd64`
- `linux/ppc65le`
- `linux/s390x`
- `windows/amd64`
- `darwin/amd64`

In order to download a binary file:
* Go to the release you are interested for `https://github.com/devfile/alizer/releases/tag/<release-tag>`
* In the **Assets** section you will see the list of generated binaries for the release.

## Feedback & Questions

If you discover an issue please file a bug, and we will fix it as soon as possible.

Issues are tracked in the [devfile/api](https://github.com/devfile/api) repo with the label [area/alizer](https://github.com/devfile/api/issues?q=is%3Aopen+is%3Aissue+label%3Aarea%2Falizer)
