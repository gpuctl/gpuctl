package gpuctl_test

import (
	"testing"

	"github.com/gpuctl/gpuctl"
)

func TestDemo(t *testing.T) {
	if gpuctl.GetANumber() < 0 {
		t.Error("this is bad")
	}
}
