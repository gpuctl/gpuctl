package config_test

import (
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeTime(t *testing.T) {
	t.Parallel()

	type Foo struct {
		A time.Duration
		B time.Duration
	}

	var foo Foo

	err := toml.Unmarshal([]byte(`A="3s"
B="2.5h"`), &foo)
	require.NoError(t, err)

	assert.Equal(t, foo.A, 3*time.Second)
	assert.Equal(t, foo.B, 2*time.Hour+30*time.Minute)
}

func TestDecodeTimeNotPresent(t *testing.T) {
	t.Parallel()

	type Bar struct {
		C time.Duration
		D int
	}

	bar := Bar{6 * time.Microsecond, 101}

	err := toml.Unmarshal([]byte(`D = 666`), &bar)
	require.NoError(t, err)

	assert.Equal(t, 6*time.Microsecond, bar.C)
	assert.Equal(t, 666, bar.D)
}
