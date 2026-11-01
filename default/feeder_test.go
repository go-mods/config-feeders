package default_feeder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nested struct {
	Value string `default:"nested"`
}

func TestFeed_InvalidStructure(t *testing.T) {
	var nilPointer *nested
	number := 1

	tests := []struct {
		name      string
		structure interface{}
	}{
		{"nil", nil},
		{"nil pointer to struct", nilPointer},
		{"struct by value", nested{}},
		{"pointer to non-struct", &number},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, Default{}.Feed(tt.structure), "default: invalid structure")
		})
	}
}

func TestFeed_SetsTaggedFields(t *testing.T) {
	var cfg struct {
		Host     string        `default:"localhost"`
		Port     int           `default:"8080"`
		Debug    bool          `default:"true"`
		Timeout  time.Duration `default:"5s"`
		Hosts    []string      `default:"a,b"`
		Nested   nested
		Set      *nested
		Nil      *nested
		Untagged string
		private  string `default:"private"`
	}
	cfg.Set = &nested{}

	require.NoError(t, Default{}.Feed(&cfg))

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, []string{"a", "b"}, cfg.Hosts)
	assert.Equal(t, "nested", cfg.Nested.Value)
	assert.Equal(t, "nested", cfg.Set.Value)
	assert.Nil(t, cfg.Nil)
	assert.Empty(t, cfg.Untagged)
	assert.Equal(t, "private", cfg.private)
}

func TestFeed_OverwritesExistingValues(t *testing.T) {
	cfg := nested{Value: "already set"}

	require.NoError(t, Default{}.Feed(&cfg))

	assert.Equal(t, "nested", cfg.Value)
}

func TestFeed_InvalidDefaultValue(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		var cfg struct {
			Port int `default:"not a number"`
		}
		assert.ErrorContains(t, Default{}.Feed(&cfg), "`Port`")
	})

	t.Run("duration", func(t *testing.T) {
		var cfg struct {
			Timeout time.Duration `default:"not a duration"`
		}
		assert.ErrorContains(t, Default{}.Feed(&cfg), "`Timeout`")
	})
}
