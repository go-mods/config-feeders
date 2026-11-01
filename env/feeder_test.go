package env_feeder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeed(t *testing.T) {
	t.Setenv("CONFIG_FEEDERS_HOST", "localhost")
	t.Setenv("CONFIG_FEEDERS_PORT", "8080")

	var cfg struct {
		Host string `env:"CONFIG_FEEDERS_HOST"`
		Port int    `env:"CONFIG_FEEDERS_PORT"`
	}

	require.NoError(t, Env{}.Feed(&cfg))

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
}

func TestFeed_InvalidStructure(t *testing.T) {
	assert.EqualError(t, Env{}.Feed(nil), "env: invalid structure")
}
