package uplink_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

// test summarising an empty processes list
func TestSummariseEmpty(t *testing.T) {
	p := uplink.Processes{}

	inUse, user := p.Summarise()

	if inUse {
		t.Errorf("inUse value incorrect. was '%v', wanted false", inUse)
	}

	if user != "" {
		t.Errorf("user value incorrect. was '%v', wanted empty string", user)
	}
}

// test summarising a single processes list
func TestSummariseOne(t *testing.T) {
	fakeUser := "dominic"

	p := uplink.Processes{
		{
			Pid:     456,
			Name:    "python",
			MemUsed: 456.7,
			Owner:   fakeUser,
		},
	}

	inUse, user := p.Summarise()

	if !inUse {
		t.Errorf("inUse value incorrect. was '%v', wanted true", inUse)
	}

	if user != fakeUser {
		t.Errorf("user value incorrect. was '%v', wanted '%v'", user, fakeUser)
	}
}

// test summarising a longer processes list
// not necessarily the behaviour we want long term, I'm just documenting this behaviour in tests
func TestSummariseMulti(t *testing.T) {
	fakeUser1 := "clive"
	fakeUser2 := "brenda"
	fakeUser3 := "steve"

	p := uplink.Processes{
		{
			Pid:     456,
			Name:    "python",
			MemUsed: 456.7,
			Owner:   fakeUser1,
		},
		{
			Pid:     12345,
			Name:    "python3",
			MemUsed: 15.3,
			Owner:   fakeUser2,
		},
		{
			Pid:     543,
			Name:    "python",
			MemUsed: 987.0,
			Owner:   fakeUser3,
		},
	}

	inUse, user := p.Summarise()

	if !inUse {
		t.Errorf("inUse value incorrect. was '%v', wanted true", inUse)
	}

	if user != fakeUser1 {
		t.Errorf("user value incorrect. was '%v', wanted '%v'", user, fakeUser1)
	}
}
