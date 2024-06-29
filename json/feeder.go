package json_feeder

import (
	"encoding/json"
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
	buf, err := os.ReadFile(filepath.Clean(f.Path))

	if err != nil {
		return fmt.Errorf("json: %v", err)
	}

	if JSONUnmarshal == nil {
		JSONUnmarshal = json.Unmarshal
	}

	if err = JSONUnmarshal(buf, structure); err != nil {
		return fmt.Errorf("json: %v", err)
	}

	return nil
}
