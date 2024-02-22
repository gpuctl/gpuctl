package procinfo

import (
	"errors"
	"strconv"

	"golang.org/x/sys/unix"

	"github.com/gpuctl/gpuctl/internal/passwd"
)

var (
	ErrorNoUserFound = errors.New("Could not find user that owns the given process from passwd file")
)

type UidLookup map[uint32]string

func (lookup UidLookup) UserForPid(pid uint64) (string, error) {
	// get uid for running process, from procFS
	filename := "/proc/" + strconv.FormatUint(pid, 10)
	var stat unix.Stat_t
	err := unix.Stat(filename, &stat)
	if err != nil {
		return "", err
	}

	name := lookup[stat.Uid]
	if name == "" {
		return "", ErrorNoUserFound
	}
	return name, nil
}

func PasswdToLookup(entries passwd.Passwd) UidLookup {
	lookup := make(UidLookup)
	for _, entry := range entries {
		// Try to assign full name to lookup, otherwise use username
		fullname := entry.ParseGecos().FullName
		if fullname == "" {
			lookup[entry.Uid] = entry.Name
		} else {
			lookup[entry.Uid] = fullname
		}
	}
	return lookup
}
