package dotenv_feeder

import (
	"fmt"
	"github.com/golobby/dotenv"
	"os"
	"path/filepath"
)

// DotEnv is a feeder to be used with golobby/config.
//
// It feeds using dot env (.env) files.
type DotEnv struct {
	Path string
}

// Feed feeds the structure with the dot env file data.
func (f DotEnv) Feed(structure interface{}) error {
	file, err := os.Open(filepath.Clean(f.Path))
	defer func(file *os.File) { _ = file.Close() }(file)

	if err != nil {
		return fmt.Errorf("config: cannot open env file; err: %v", err)
	}

	if err = dotenv.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("dotenv: %v", err)
	}

	return nil
}
