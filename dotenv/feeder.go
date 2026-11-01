// Package dotenv_feeder provides a golobby/config feeder that fills a structure
// with the content of a dot env (.env) file, using its `env` tags.
package dotenv_feeder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golobby/dotenv"
)

// DotEnv is a feeder to be used with golobby/config.
//
// It feeds using dot env (.env) files.
type DotEnv struct {
	Path string
}

// Feed feeds the structure with the dot env file data.
func (f DotEnv) Feed(structure interface{}) error {
	if f.Path == "" {
		return errors.New("dotenv: empty path")
	}

	file, err := os.Open(filepath.Clean(f.Path))
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	if err = dotenv.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("%s: %w", f.Path, err)
	}

	return nil
}
