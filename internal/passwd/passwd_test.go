package passwd_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/gpuctl/gpuctl/internal/passwd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/csg_complex
var csgComplex string

//go:embed testdata/ashtabula
var ashtabula string

func TestGecos4(t *testing.T) {
	t.Parallel()

	users, err := passwd.Parse(strings.NewReader(csgComplex))
	require.NoError(t, err)
	assert.Len(t, users, 27)

	var duncan passwd.Entry
	for _, user := range users {
		if user.Name == "dcw" {
			duncan = user
		}
	}

	assert.Equal(t, "*", duncan.Password)
	assert.Equal(t, uint32(343), duncan.Uid)
	assert.Equal(t, uint32(100), duncan.Gid)
	assert.Equal(t, "/homes/dcw", duncan.HomeDir)
	assert.Equal(t, "/bin/tcsh", duncan.Shell)

	gecos := duncan.ParseGecos()
	assert.Equal(t, "Duncan C White", gecos.FullName)
	assert.Equal(t, "huxley 305", gecos.Office)
	assert.Equal(t, "48254", gecos.OfficePhone)
	assert.Equal(t, "48222", gecos.HomePhone)
}

func loadGecos(t *testing.T, user string) passwd.Gecos {
	t.Helper()

	users, err := passwd.Parse(strings.NewReader(csgComplex))
	require.NoError(t, err)

	for _, ent := range users {
		if ent.Name == user {
			return ent.ParseGecos()
		}
	}

	t.Errorf("No user called %s", user)
	return passwd.Gecos{} // No bottom type :'(
}

func TestGecos3(t *testing.T) {
	t.Parallel()

	gecos := loadGecos(t, "svb")
	assert.Equal(t, gecos.FullName, "Steffen van Bakel")
	assert.Equal(t, gecos.Office, "huxley 425")
	assert.Equal(t, gecos.OfficePhone, "48263")
	assert.Equal(t, gecos.HomePhone, "")
}

func TestGecos2(t *testing.T) {
	t.Parallel()

	gecos := loadGecos(t, "bglocker")
	assert.Equal(t, gecos.FullName, "Ben Glocker")
	assert.Equal(t, gecos.Office, "huxley 377")
	assert.Equal(t, gecos.OfficePhone, "")
	assert.Equal(t, gecos.HomePhone, "")
}

func TestGecos1(t *testing.T) {
	t.Parallel()

	gecos := loadGecos(t, "ae321")
	assert.Equal(t, gecos.FullName, "Alona Enraght-Moony")
	assert.Equal(t, gecos.Office, "")
	assert.Equal(t, gecos.OfficePhone, "")
	assert.Equal(t, gecos.HomePhone, "")
}

func TestGecos0(t *testing.T) {
	t.Parallel()

	gecos := loadGecos(t, "sshd")
	assert.Zero(t, gecos)
}

func TestCsgNoPanic(t *testing.T) {
	t.Parallel()

	users, err := passwd.Parse(strings.NewReader(csgComplex))
	require.NoError(t, err)

	for _, u := range users {
		u.ParseGecos()
	}
}

func TestAshtabulaNoPanic(t *testing.T) {
	t.Parallel()

	users, err := passwd.Parse(strings.NewReader(ashtabula))
	require.NoError(t, err)

	for _, u := range users {
		u.ParseGecos()
	}
}
