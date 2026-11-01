// Package env_feeder provides a golobby/config feeder that fills a structure
// with the environment variables named by its `env` tags.
package env_feeder

import (
	"github.com/golobby/env/v2"
)

// Env is a feeder to be used with golobby/config.
//
// It feeds using environment variables and struct tags.
type Env struct{}

// Feed feeds the structure with the environment variables.
func (f Env) Feed(structure interface{}) error {
	return env.Feed(structure)
}
