package toml_feeder

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
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
	buf, err := os.ReadFile(filepath.Clean(f.Path))

	if err != nil {
		return fmt.Errorf("toml: %v", err)
	}

	if TOMLUnmarshal == nil {
		TOMLUnmarshal = toml.Unmarshal
	}

	if err = TOMLUnmarshal(buf, structure); err != nil {
		return fmt.Errorf("toml: %v", err)
	}

	return nil
}
