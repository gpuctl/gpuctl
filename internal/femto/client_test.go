package femto_test

import (
	"errors"
	"net"
	"testing"

	"github.com/gpuctl/gpuctl/internal/femto"
)

func TestPostToUnresolvable(t *testing.T) {
	t.Parallel()

	err := femto.Post("https://lol.invalid", 101)

	var target *net.DNSError
	if !errors.As(err, &target) {
		t.Fatal("DNS error, but got", err)
	}
	if target.Name != "lol.invalid" {
		t.Fatal("Expected target name `lol.invalid`, but got: ", target.Name)
	}
}
