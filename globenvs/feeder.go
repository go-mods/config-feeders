// Package globenvs_feeder provides a golobby/config feeder that fills a structure
// with the content of every dot env (.env) file matching a list of patterns.
package globenvs_feeder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golobby/dotenv"
)

// GlobEnvs is a feeder to be used with golobby/config.
//
// It feeds using multiple sources of .env files.
type GlobEnvs struct {
	Patterns []string
}

// Feed feeds the structure with the .env files matching the patterns, in the order
// of the patterns: a later file overrides an earlier one.
//
// The files are optional: a pattern that matches nothing is not an error, and the
// structure is left untouched when no file matches at all.
func (f GlobEnvs) Feed(structure interface{}) error {
	paths, err := f.glob()
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err := feedFile(path, structure); err != nil {
			return err
		}
	}

	return nil
}

// glob returns the paths of the files matching the patterns.
func (f GlobEnvs) glob() ([]string, error) {
	var paths []string

	for _, pattern := range f.Patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", pattern, err)
		}
		paths = append(paths, matches...)
	}

	return paths, nil
}

// feedFile feeds the structure with a single .env file.
func feedFile(path string, structure interface{}) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	if err = dotenv.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	return nil
}
