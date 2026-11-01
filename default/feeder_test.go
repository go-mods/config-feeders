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

type (
	level   string
	count   int
	levels  []string
	timeout time.Duration
)

func TestFeed_SetsNamedTypesAndSlices(t *testing.T) {
	var cfg struct {
		Level     level           `default:"debug"`
		Count     count           `default:"3"`
		Levels    levels          `default:"a, b"`
		Items     []level         `default:"c,d"`
		Durations []time.Duration `default:"1s,2m"`
		Timeout   timeout         `default:"5"`
		Int64     int64           `default:"5000000000"`
	}

	require.NoError(t, Default{}.Feed(&cfg))

	assert.Equal(t, level("debug"), cfg.Level)
	assert.Equal(t, count(3), cfg.Count)
	assert.Equal(t, levels{"a", "b"}, cfg.Levels)
	assert.Equal(t, []level{"c", "d"}, cfg.Items)
	assert.Equal(t, []time.Duration{time.Second, 2 * time.Minute}, cfg.Durations)
	assert.Equal(t, timeout(5), cfg.Timeout)
	assert.Equal(t, int64(5000000000), cfg.Int64)
}

func TestFeed_UnsupportedTypes(t *testing.T) {
	tests := []struct {
		name      string
		structure interface{}
	}{
		{"map", &struct {
			Field map[string]string `default:"a=b"`
		}{}},
		{"pointer", &struct {
			Field *string `default:"a"`
		}{}},
		{"struct", &struct {
			Field nested `default:"a"`
		}{}},
		{"time", &struct {
			Field time.Time `default:"2024-01-01T00:00:00Z"`
		}{}},
		{"slice of slices", &struct {
			Field [][]string `default:"a,b"`
		}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorContains(t, Default{}.Feed(tt.structure), "`Field`")
		})
	}
}

func TestFeed_CyclicStructure(t *testing.T) {
	type node struct {
		Value string `default:"node"`
		Next  *node
	}
	first := &node{}
	second := &node{Next: first}
	first.Next = second

	require.NoError(t, Default{}.Feed(first))

	assert.Equal(t, "node", first.Value)
	assert.Equal(t, "node", second.Value)
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
