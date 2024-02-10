package assets_test

import (
	"bytes"
	"debug/elf"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gpuctl/gpuctl/internal/assets"
)

func TestSatelliteAmd64LinuxIsElf(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader(assets.SatelliteAmd64Linux)
	elfFile, err := elf.NewFile(reader)
	require.NoError(t, err)

	imports, err := elfFile.ImportedLibraries()
	assert.Empty(t, imports,
		"satellite binary has imports, but should be static binary")
	assert.NoError(t, err)

	assert.Equal(t, elf.EM_X86_64, elfFile.Machine)
	assert.Equal(t, elf.ET_EXEC, elfFile.Type)
}
