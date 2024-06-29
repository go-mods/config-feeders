
# Config feeders

This package contains feeders that provide configuration data to the GoLobby Config package.

## Feeders

| Feeder Name | File                 | Description                                                                       |
|-------------|----------------------|-----------------------------------------------------------------------------------|
| Default     | `default/feeder.go`  | Fills a structure with default values specified in structure tags.                |
| DotEnv      | `dotenv/feeder.go`   | Feeds a structure with data from a .env file.                                     |
| GlobEnvs    | `globenvs/feeder.go` | Feeds a structure with data from multiple .env files matching specified patterns. |
| Env         | `env/feeder.go`      | Feeds a structure with data from environment variables.                           |
| Json        | `json/feeder.go`     | Feeds a structure with data from a JSON file.                                     |
| Yaml        | `yaml/feeder.go`     | Feeds a structure with data from a YAML file.                                     |
| Toml        | `toml/feeder.go`     | Feeds a structure with data from a TOML file.                                     |
