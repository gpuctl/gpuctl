package procinfo

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/gpuctl/gpuctl/internal/passwd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: could not import same file in a different module, dont want to hardcode
//
//go:embed testdata/ashtabula
var ashtabula string

func loadAshtabula(t *testing.T) passwd.Passwd {
	t.Helper()
	users, err := passwd.Parse(strings.NewReader(ashtabula))
	require.NoError(t, err)
	return users
}

func TestPasswdToLookup(t *testing.T) {
	t.Parallel()
	entries := loadAshtabula(t)
	lookup := PasswdToLookup(entries)
	assert.Equal(t, "Mailing List Manager", lookup[38])
	assert.Equal(t, "alona", lookup[1000])
}

func TestPidLookupOnSelf(t *testing.T) {
	t.Parallel()
	pid := uint64(os.Getpid())
	uid := uint32(os.Getuid())
	fakemap := UidLookup{uid: "Name"}
	res, err := fakemap.UserForPid(pid)
	require.NoError(t, err)
	assert.Equal(t, "Name", res)
}
