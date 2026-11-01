// Command example shows each feeder, then several feeders combined.
//
// The paths are relative to the root of the module: run it from there with
// go run ./example
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golobby/config/v3"

	default_feeder "github.com/go-mods/config-feeders/default"
	dotenv_feeder "github.com/go-mods/config-feeders/dotenv"
	env_feeder "github.com/go-mods/config-feeders/env"
	globenvs_feeder "github.com/go-mods/config-feeders/globenvs"
	json_feeder "github.com/go-mods/config-feeders/json"
	toml_feeder "github.com/go-mods/config-feeders/toml"
	yaml_feeder "github.com/go-mods/config-feeders/yaml"
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := os.Setenv("HOST", "example.com"); err != nil {
		return err
	}
	if err := os.Setenv("PORT", "9090"); err != nil {
		return err
	}

	defaultFeeder := default_feeder.Default{}
	envFeeder := env_feeder.Env{}
	jsonFeeder := json_feeder.Json{Path: "./example/config.json"}

	examples := []struct {
		title   string
		feeders []config.Feeder
	}{
		{"Default", []config.Feeder{defaultFeeder}},
		{"DotEnv", []config.Feeder{dotenv_feeder.DotEnv{Path: "./example/.env"}}},
		{"GlobEnvs", []config.Feeder{globenvs_feeder.GlobEnvs{Patterns: []string{"./example/.env.*"}}}},
		{"Env", []config.Feeder{envFeeder}},
		{"JSON", []config.Feeder{jsonFeeder}},
		{"YAML", []config.Feeder{yaml_feeder.Yaml{Path: "./example/config.yaml"}}},
		{"TOML", []config.Feeder{toml_feeder.Toml{Path: "./example/config.toml"}}},
		// Each feeder overrides the previous ones
		{"Multiple", []config.Feeder{defaultFeeder, jsonFeeder, envFeeder}},
	}

	for _, e := range examples {
		fmt.Printf("%s Feeder Example:\n", e.title)

		var cfg Config
		if err := config.New().AddFeeder(e.feeders...).AddStruct(&cfg).Feed(); err != nil {
			return fmt.Errorf("%s feeder: %w", e.title, err)
		}
		printConfig(cfg)
		fmt.Println()
	}

	return nil
}

func printConfig(cfg Config) {
	fmt.Printf("Host: %s\n", cfg.Host)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug: %v\n", cfg.Debug)
	fmt.Printf("Database Name: %s\n", cfg.Database.Name)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("Database Password: %s\n", cfg.Database.Password)
}
