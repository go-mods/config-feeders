// Package yaml_feeder provides a golobby/config feeder that fills a structure
// with the content of a YAML file.
package yaml_feeder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Yaml is a feeder to be used with golobby/config.
//
// It feeds using a YAML file.
type Yaml struct {
	Path string
}

// YAMLUnmarshal is a function that unmarshals YAML data.
// If you want to use a custom YAML unmarshal function, you can set this variable.
// By default, it uses yaml.Unmarshal.
var YAMLUnmarshal func(data []byte, v interface{}) error

// Feed feeds the structure with the YAML file data.
func (f Yaml) Feed(structure interface{}) error {
	if f.Path == "" {
		return errors.New("yaml: empty path")
	}

	buf, err := os.ReadFile(filepath.Clean(f.Path))
	if err != nil {
		return err
	}

	unmarshal := YAMLUnmarshal
	if unmarshal == nil {
		unmarshal = yaml.Unmarshal
	}

	if err = unmarshal(buf, structure); err != nil {
		return fmt.Errorf("%s: %w", f.Path, err)
	}

	return nil
}
