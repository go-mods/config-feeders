package globenvs_feeder

import (
	"fmt"
	"github.com/golobby/dotenv"
	"os"
	"path/filepath"
)

// GlobEnvs is a feeder to be used with golobby/config.
//
// It feeds using multiple sources of .env files.
type GlobEnvs struct {
	Patterns []string
}

// Feed feeds the structure with the .env files data.

// Feed permet de remplir la structure avec les valeurs des fichiers .env
func (f GlobEnvs) Feed(structure interface{}) error {
	// glob files
	paths, err := f.glob()
	if err != nil {
		return err
	}

	// feed files
	for _, path := range paths {
		file, err := os.Open(filepath.Clean(path))

		if err != nil {
			return fmt.Errorf("config: cannot open env file %s; err: %v", path, err)
		}

		if err = dotenv.NewDecoder(file).Decode(structure); err != nil {
			return fmt.Errorf("dotenv: %v", err)
		}

		if err = file.Close(); err != nil {
			return fmt.Errorf("config: cannot close env file %s; err: %v", path, err)
		}
	}

	return nil
}

// glob permet de rechercher les fichiers correspondant au motif dans le répertoire dir
func (f GlobEnvs) glob() ([]string, error) {
	var paths []string

	for _, pattern := range f.Patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("config: cannot glob env files; err: %v", err)
		}
		paths = append(paths, matches...)
	}

	return paths, nil
}
