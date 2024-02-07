package gpustats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// IMMENSELY STUPID WE HAVE TO WRITE TESTS FOR THESE

func TestFakeGPU(t *testing.T) {
	var f FakeGPU
	a, _ := f.GetGPUInformation()
	b, _ := f.GetGPUStatus()
	assert.Equal(t, a[0].Uuid, "some_id")
	assert.Equal(t, b[0].Uuid, "some_id")
}
