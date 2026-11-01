// Package toml_feeder provides a golobby/config feeder that fills a structure
// with the content of a TOML file.
package toml_feeder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Toml is a feeder to be used with golobby/config.
//
// It feeds using a TOML file.
type Toml struct {
	Path string
}

// TOMLUnmarshal is a function that unmarshals TOML data.
// If you want to use a custom TOML unmarshal function, you can set this variable.
// By default, it uses toml.Unmarshal.
var TOMLUnmarshal func(data []byte, v interface{}) error

// Feed feeds the structure with the TOML file data.
func (f Toml) Feed(structure interface{}) error {
	if f.Path == "" {
		return errors.New("toml: empty path")
	}

	buf, err := os.ReadFile(filepath.Clean(f.Path))
	if err != nil {
		return err
	}

	unmarshal := TOMLUnmarshal
	if unmarshal == nil {
		unmarshal = toml.Unmarshal
	}

	if err = unmarshal(buf, structure); err != nil {
		return fmt.Errorf("%s: %w", f.Path, err)
	}

	return nil
}
