package incrementer

import (
	"testing"
)

// make sure that incrementing a number increments it
func TestIncrement(t *testing.T) {
	value := 1
	inc := Inc(value)
	if inc != 2 {
		t.Fatalf(`Inc(%d) = %d, wanted 2`, value, inc)
	}
}

func TestIncNegative(t *testing.T) {
	value := -45
	inc := Inc(value)
	if inc != -44 {
		t.Fatalf(`Inc(%d) = %d, wanted -44`, value, inc)
	}
}
