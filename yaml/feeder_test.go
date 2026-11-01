package yaml_feeder

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type settings struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFeed(t *testing.T) {
	var cfg settings

	require.NoError(t, Yaml{Path: writeFile(t, "host: localhost\nport: 8080\n")}.Feed(&cfg))

	assert.Equal(t, settings{Host: "localhost", Port: 8080}, cfg)
}

func TestFeed_DoesNotChangeTheUnmarshalVariable(t *testing.T) {
	require.NoError(t, Yaml{Path: writeFile(t, "")}.Feed(&settings{}))

	assert.Nil(t, YAMLUnmarshal)
}

func TestFeed_CustomUnmarshal(t *testing.T) {
	errCustom := errors.New("custom")
	YAMLUnmarshal = func([]byte, interface{}) error { return errCustom }
	t.Cleanup(func() { YAMLUnmarshal = nil })

	err := Yaml{Path: writeFile(t, "")}.Feed(&settings{})

	assert.ErrorIs(t, err, errCustom)
}

func TestFeed_Errors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.EqualError(t, Yaml{}.Feed(&settings{}), "yaml: empty path")
	})

	t.Run("missing file", func(t *testing.T) {
		err := Yaml{Path: filepath.Join(t.TempDir(), "missing.yaml")}.Feed(&settings{})
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("invalid file names the path", func(t *testing.T) {
		path := writeFile(t, "host: [")
		assert.ErrorContains(t, Yaml{Path: path}.Feed(&settings{}), path)
	})
}
