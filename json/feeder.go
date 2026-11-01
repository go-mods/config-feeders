// Package json_feeder provides a golobby/config feeder that fills a structure
// with the content of a JSON file.
package json_feeder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Json is a feeder to be used with golobby/config.
//
// It feeds using a JSON file.
type Json struct {
	Path string
}

// JSONUnmarshal is a function that unmarshals JSON data.
// If you want to use a custom JSON unmarshal function, you can set this variable.
// By default, it uses json.Unmarshal.
var JSONUnmarshal func(data []byte, v interface{}) error

// Feed feeds the structure with the JSON file data.
func (f Json) Feed(structure interface{}) error {
	if f.Path == "" {
		return errors.New("json: empty path")
	}

	buf, err := os.ReadFile(filepath.Clean(f.Path))
	if err != nil {
		return err
	}

	unmarshal := JSONUnmarshal
	if unmarshal == nil {
		unmarshal = json.Unmarshal
	}

	if err = unmarshal(buf, structure); err != nil {
		return fmt.Errorf("%s: %w", f.Path, err)
	}

	return nil
}
