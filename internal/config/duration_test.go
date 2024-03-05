package config_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gpuctl/gpuctl/internal/config"
)

func TestDecodeTime(t *testing.T) {
	t.Parallel()

	type Foo struct {
		A config.Duration
		B config.Duration
	}

	var foo Foo

	err := toml.Unmarshal([]byte(`A="3s"
B="2.5h"`), &foo)
	require.NoError(t, err)

	assert.Equal(t, foo.A, 3*config.Second)
	assert.Equal(t, foo.B, 2*config.Hour+30*config.Minute)
}

func TestDecodeTimeNotPresent(t *testing.T) {
	t.Parallel()

	type Bar struct {
		C config.Duration
		D int
	}

	bar := Bar{6 * config.Microsecond, 101}

	err := toml.Unmarshal([]byte(`D = 666`), &bar)
	require.NoError(t, err)

	assert.Equal(t, 6*config.Microsecond, bar.C)
	assert.Equal(t, 666, bar.D)
}
