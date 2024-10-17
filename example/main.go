package main

import (
	"fmt"
	"github.com/go-mods/config-feeders/default"
	"github.com/go-mods/config-feeders/dotenv"
	"github.com/go-mods/config-feeders/env"
	"github.com/go-mods/config-feeders/globenvs"
	"github.com/go-mods/config-feeders/json"
	"github.com/go-mods/config-feeders/toml"
	"github.com/go-mods/config-feeders/yaml"
	"github.com/golobby/config/v3"
	"os"
)

type Config struct {
	Host     string `json:"host" yaml:"host" toml:"host" env:"HOST" default:"localhost"`
	Port     int    `json:"port" yaml:"port" toml:"port" env:"PORT" default:"8080"`
	Debug    bool   `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG" default:"false"`
	Database struct {
		Name     string `json:"name" yaml:"name" toml:"name" env:"DB_NAME" default:"mydb"`
		User     string `json:"user" yaml:"user" toml:"user" env:"DB_USER" default:"root"`
		Password string `json:"password" yaml:"password" toml:"password" env:"DB_PASSWORD" default:""`
	} `json:"database" yaml:"database" toml:"database"`
}

func main() {
	// Create a new config instance
	cfg := Config{}

	// Example for Default Feeder
	fmt.Println("Default Feeder Example:")
	defaultFeeder := default_feeder.Default{}
	c := config.New()
	_ = c.AddFeeder(defaultFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for DotEnv Feeder
	fmt.Println("\nDotEnv Feeder Example:")
	dotenvFeeder := dotenv_feeder.DotEnv{Path: "./example/.env"}
	c = config.New()
	_ = c.AddFeeder(dotenvFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for GlobEnvs Feeder
	fmt.Println("\nGlobEnvs Feeder Example:")
	globenvsFeeder := globenvs_feeder.GlobEnvs{Patterns: []string{"./example/.env.*"}}
	c = config.New()
	_ = c.AddFeeder(globenvsFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for Env Feeder
	fmt.Println("\nEnv Feeder Example:")
	os.Setenv("HOST", "example.com") // #nosec G104
	os.Setenv("PORT", "9090")        // #nosec G104
	envFeeder := env_feeder.Env{}
	c = config.New()
	_ = c.AddFeeder(envFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for JSON Feeder
	fmt.Println("\nJSON Feeder Example:")
	jsonFeeder := json_feeder.Json{Path: "./example/config.json"}
	c = config.New()
	_ = c.AddFeeder(jsonFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for YAML Feeder
	fmt.Println("\nYAML Feeder Example:")
	yamlFeeder := yaml_feeder.Yaml{Path: "./example/config.yaml"}
	c = config.New()
	_ = c.AddFeeder(yamlFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example for TOML Feeder
	fmt.Println("\nTOML Feeder Example:")
	tomlFeeder := toml_feeder.Toml{Path: "./example/config.toml"}
	c = config.New()
	_ = c.AddFeeder(tomlFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)

	// Example of using multiple feeders
	fmt.Println("\nMultiple Feeders Example:")
	c = config.New()
	_ = c.AddFeeder(defaultFeeder)
	_ = c.AddFeeder(jsonFeeder)
	_ = c.AddFeeder(envFeeder)
	_ = c.AddStruct(&cfg)
	_ = c.Feed()
	printConfig(cfg)
}

func printConfig(cfg Config) {
	fmt.Printf("Host: %s\n", cfg.Host)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug: %v\n", cfg.Debug)
	fmt.Printf("Database Name: %s\n", cfg.Database.Name)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("Database Password: %s\n", cfg.Database.Password)
}
