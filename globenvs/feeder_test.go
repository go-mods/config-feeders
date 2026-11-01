package globenvs_feeder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type settings struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFeed_LaterFilesOverrideEarlierOnes(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.env", "HOST=first\nPORT=8080\n")
	writeFile(t, dir, "b.env", "HOST=second\n")
	writeFile(t, dir, "override", "PORT=9090\n")

	var cfg settings
	patterns := []string{filepath.Join(dir, "*.env"), filepath.Join(dir, "override")}

	require.NoError(t, GlobEnvs{Patterns: patterns}.Feed(&cfg))

	assert.Equal(t, settings{Host: "second", Port: 9090}, cfg)
}

func TestFeed_NoMatchIsNotAnError(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{"no pattern", nil},
		{"pattern without match", []string{filepath.Join(t.TempDir(), "*.env")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := settings{Host: "unchanged"}

			require.NoError(t, GlobEnvs{Patterns: tt.patterns}.Feed(&cfg))

			assert.Equal(t, settings{Host: "unchanged"}, cfg)
		})
	}
}

func TestFeed_BadPatternNamesThePattern(t *testing.T) {
	err := GlobEnvs{Patterns: []string{"["}}.Feed(&settings{})

	require.ErrorIs(t, err, filepath.ErrBadPattern)
	assert.ErrorContains(t, err, "[")
}

func TestFeed_InvalidFileNamesThePathAndClosesIt(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "invalid.env", "INVALID\n")

	err := GlobEnvs{Patterns: []string{path}}.Feed(&settings{})

	require.ErrorContains(t, err, path)
	// Windows refuses to delete a file that is still open.
	assert.NoError(t, os.Remove(path))
}
