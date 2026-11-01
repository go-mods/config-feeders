package json_feeder

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type settings struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFeed(t *testing.T) {
	var cfg settings

	require.NoError(t, Json{Path: writeFile(t, `{"host":"localhost","port":8080}`)}.Feed(&cfg))

	assert.Equal(t, settings{Host: "localhost", Port: 8080}, cfg)
}

func TestFeed_DoesNotChangeTheUnmarshalVariable(t *testing.T) {
	require.NoError(t, Json{Path: writeFile(t, `{}`)}.Feed(&settings{}))

	assert.Nil(t, JSONUnmarshal)
}

func TestFeed_CustomUnmarshal(t *testing.T) {
	errCustom := errors.New("custom")
	JSONUnmarshal = func([]byte, interface{}) error { return errCustom }
	t.Cleanup(func() { JSONUnmarshal = nil })

	err := Json{Path: writeFile(t, `{}`)}.Feed(&settings{})

	assert.ErrorIs(t, err, errCustom)
}

func TestFeed_Errors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.EqualError(t, Json{}.Feed(&settings{}), "json: empty path")
	})

	t.Run("missing file", func(t *testing.T) {
		err := Json{Path: filepath.Join(t.TempDir(), "missing.json")}.Feed(&settings{})
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("invalid file names the path", func(t *testing.T) {
		path := writeFile(t, `{`)
		assert.ErrorContains(t, Json{Path: path}.Feed(&settings{}), path)
	})
}
