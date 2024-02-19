package procinfo

import (
	"errors"
	"github.com/gpuctl/gpuctl/internal/passwd"
	"os"
	"strconv"
	"syscall"
)

var (
	ErrorNoUserFound = errors.New("Could not find user that owns the given process from passwd file")
)

type UidLookup map[uint32]string

func (lookup UidLookup) Get(pid uint64) (string, error) {
	var zero string

	// get uid
	filename := "/proc/" + strconv.FormatUint(pid, 10)
	statcall, err := os.Stat(filename)
	if err != nil {
		return zero, err
	}
	stat, okay := statcall.Sys().(*syscall.Stat_t)
	if okay {
		return "", err
	}
	uid := uint32(stat.Uid)

	name := lookup[uid]
	if name == zero {
		return zero, ErrorNoUserFound
	}
	return name, nil
}

func PasswdToLookup(entries []passwd.Entry) UidLookup {
	var lookup UidLookup
	for _, entry := range entries {
		lookup[entry.Uid] = entry.ParseGecos().FullName
	}
	return lookup
}
