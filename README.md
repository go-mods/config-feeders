# Config feeders

[![Go Reference](https://pkg.go.dev/badge/github.com/go-mods/config-feeders.svg)](https://pkg.go.dev/github.com/go-mods/config-feeders)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mods/config-feeders)](https://goreportcard.com/report/github.com/go-mods/config-feeders)
[![Release](https://img.shields.io/github/release/go-mods/config-feeders.svg)](https://github.com/go-mods/config-feeders/releases/latest)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/go-mods/config-feeders/blob/master/LICENSE.md)

`Config Feeders` provides a collection of feeders for the [GoLobby Config](https://github.com/golobby/config) package. They load configuration data from environment variables, .env files, JSON, YAML and TOML files, as well as default values specified in struct tags.

## Installation

Prerequisite: Go 1.23 or later.

```bash
go get github.com/go-mods/config-feeders
```

## Usage

Import the desired feeder package, create an instance of the feeder and use it with the GoLobby Config package. Here's an example using the JSON feeder:

```go
package main

import (
	"fmt"

	json_feeder "github.com/go-mods/config-feeders/json"
	"github.com/golobby/config/v3"
)

type MyConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	var myConfig MyConfig

	jsonFeeder := json_feeder.Json{Path: "config.json"}

	err := config.New().AddFeeder(jsonFeeder).AddStruct(&myConfig).Feed()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Host: %s, Port: %d\n", myConfig.Host, myConfig.Port)
}
```

A runnable example of every feeder is available in [example/main.go](example/main.go). Run it from the root of the module with `go run ./example`.

## Feeders

The following feeders are available in this package:

| Feeder Name | File               | Struct tag | Description                                                                       |
|-------------|--------------------|------------|-----------------------------------------------------------------------------------|
| Default     | default/feeder.go  | `default`  | Fills a structure with default values specified in structure tags.                |
| DotEnv      | dotenv/feeder.go   | `env`      | Feeds a structure with data from a .env file.                                     |
| GlobEnvs    | globenvs/feeder.go | `env`      | Feeds a structure with data from multiple .env files matching specified patterns. |
| Env         | env/feeder.go      | `env`      | Feeds a structure with data from environment variables.                           |
| Json        | json/feeder.go     | `json`     | Feeds a structure with data from a JSON file.                                     |
| Yaml        | yaml/feeder.go     | `yaml`     | Feeds a structure with data from a YAML file.                                     |
| Toml        | toml/feeder.go     | `toml`     | Feeds a structure with data from a TOML file.                                     |

### Order of the feeders

Feeders run in the order they were added, and each one overrides the values set by the previous ones. `Default` always overwrites its tagged fields, even with an empty `default:""`: add it first.

```go
config.New().AddFeeder(default_feeder.Default{}, json_feeder.Json{Path: "config.json"}, env_feeder.Env{})
```

### Optional files

`GlobEnvs` treats its files as optional: a pattern that matches nothing is not an error, and the structure is left untouched when no file matches at all. Files are read in the order of the patterns, so a later file overrides an earlier one. The other file feeders return an error when their file is missing.

### Custom unmarshal functions

The file feeders use `json.Unmarshal`, `yaml.Unmarshal` and `toml.Unmarshal`. To use another function, set `json_feeder.JSONUnmarshal`, `yaml_feeder.YAMLUnmarshal` or `toml_feeder.TOMLUnmarshal`. These are package-level variables: set them once at startup, before any feeder runs.

### Limitations

- `Default` supports strings, booleans, integers, floats, named types based on them, `time.Duration`, and slices of these (comma-separated values).
- `DotEnv`, `GlobEnvs` and `Env` support the same built-in types and slices of them, but return an error on a `time.Duration` or a named type, as soon as a value is set for that field.
- With `Default`, `DotEnv`, `GlobEnvs` and `Env`:
  - an `int` or a `uint` is limited to 32 bits: use `int64` or `uint64` for larger values;
  - a nil pointer to a struct and the elements of a slice of structs are ignored;
  - unexported fields are set too.
- `Env` ignores an empty environment variable: it cannot reset a value to empty.
- `DotEnv`, `GlobEnvs` and `Env` do not support cyclic structures.
- When a feeder returns an error, the structure may already be partially filled.

## Development

Prerequisite: [Task](https://taskfile.dev/installation/).

```bash
task go-tools:install   # install the tools at their pinned versions
task go:test            # run the unit tests
task go:lint            # run the linters
task go:fmt             # format the code
task go:default         # run the complete validation sequence
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- [GoLobby Config](https://github.com/golobby/config) for providing the configuration management framework.
