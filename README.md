# Config feeders

[![Go Reference](https://pkg.go.dev/badge/github.com/go-mods/config-feeders.svg)](https://pkg.go.dev/github.com/go-mods/config-feeders)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mods/config-feeders)](https://goreportcard.com/report/github.com/go-mods/config-feeders)
[![Release](https://img.shields.io/github/release/go-mods/config-feeders.svg)](https://github.com/go-mods/config-feeders/releases/latest)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/go-mods/config-feeders/blob/master/LICENSE.md)


`Config Feeders` is a comprehensive package that provides a collection of feeders for the [GoLobby Config](https://github.com/golobby/config) package. These feeders allow you to easily load configuration data from various sources, including environment variables, .env files, JSON, YAML, and TOML files, as well as default values specified in struct tags.

This package aims to simplify the process of configuration management in Go applications by offering a flexible and extensible set of tools to handle different configuration formats and sources.

## Installation

To install Config Feeders, use the following command:

```bash
  go get github.com/go-mods/config-feeders
```

## Usage

To use Config Feeders, import the desired feeder package and create an instance of the feeder. Then, use it with the GoLobby Config package.
Here's an example using the JSON feeder:

```go
package main

import (
    "fmt"
    "github.com/go-mods/config-feeders/json"
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


## Feeders

The following feeders are available in this package:

| Feeder Name | File               | Description                                                                       |
|-------------|--------------------|-----------------------------------------------------------------------------------|
| Default     | default/feeder.go  | Fills a structure with default values specified in structure tags.                |
| DotEnv      | dotenv/feeder.go   | Feeds a structure with data from a .env file.                                     |
| GlobEnvs    | globenvs/feeder.go | Feeds a structure with data from multiple .env files matching specified patterns. |
| Env         | env/feeder.go      | Feeds a structure with data from environment variables.                           |
| Json        | json/feeder.go     | Feeds a structure with data from a JSON file.                                     |
| Yaml        | yaml/feeder.go     | Feeds a structure with data from a YAML file.                                     |
| Toml        | toml/feeder.go     | Feeds a structure with data from a TOML file.                                     |


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- [GoLobby Config](https://github.com/golobby/config) for providing the configuration management framework.
