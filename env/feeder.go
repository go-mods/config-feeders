package env_feeder

import (
	"fmt"
	"github.com/golobby/env/v2"
)

// Env is a feeder to be used with golobby/config.
//
// It feeds using environment variables and struct tags.
type Env struct{}

// Feed feeds the structure with the environment variables.
func (f Env) Feed(structure interface{}) error {
	if err := env.Feed(structure); err != nil {
		return fmt.Errorf("env: %v", err)
	}

	return nil
}
