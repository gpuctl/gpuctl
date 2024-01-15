package gpuctl_test

import (
	"testing"

	"github.com/gpuctl/gpuctl"
)

func TestAdd(t *testing.T) {
	if gpuctl.Add(2, 2) != 5 {
		t.Error("Literally 1984")
	}
}
