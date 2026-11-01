package toml_feeder

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type settings struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFeed(t *testing.T) {
	var cfg settings

	require.NoError(t, Toml{Path: writeFile(t, "host = \"localhost\"\nport = 8080\n")}.Feed(&cfg))

	assert.Equal(t, settings{Host: "localhost", Port: 8080}, cfg)
}

func TestFeed_DoesNotChangeTheUnmarshalVariable(t *testing.T) {
	require.NoError(t, Toml{Path: writeFile(t, "")}.Feed(&settings{}))

	assert.Nil(t, TOMLUnmarshal)
}

func TestFeed_CustomUnmarshal(t *testing.T) {
	errCustom := errors.New("custom")
	TOMLUnmarshal = func([]byte, interface{}) error { return errCustom }
	t.Cleanup(func() { TOMLUnmarshal = nil })

	err := Toml{Path: writeFile(t, "")}.Feed(&settings{})

	assert.ErrorIs(t, err, errCustom)
}

func TestFeed_Errors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.EqualError(t, Toml{}.Feed(&settings{}), "toml: empty path")
	})

	t.Run("missing file", func(t *testing.T) {
		err := Toml{Path: filepath.Join(t.TempDir(), "missing.toml")}.Feed(&settings{})
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("invalid file names the path", func(t *testing.T) {
		path := writeFile(t, "host = ")
		assert.ErrorContains(t, Toml{Path: path}.Feed(&settings{}), path)
	})
}
