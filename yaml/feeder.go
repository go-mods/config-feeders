package yaml_feeder

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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
	buf, err := os.ReadFile(filepath.Clean(f.Path))

	if err != nil {
		return fmt.Errorf("yaml: %v", err)
	}

	if YAMLUnmarshal == nil {
		YAMLUnmarshal = yaml.Unmarshal
	}

	if err = YAMLUnmarshal(buf, structure); err != nil {
		return fmt.Errorf("yaml: %v", err)
	}

	return nil
}
