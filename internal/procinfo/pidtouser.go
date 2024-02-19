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

func (lookup UidLookup) UserForPid(pid uint64) (string, error) {

	// get uid
	filename := "/proc/" + strconv.FormatUint(pid, 10)
	statcall, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	stat, ok := statcall.Sys().(*syscall.Stat_t)
	if !ok {
		return "", err
	}
	uid := uint32(stat.Uid)

	name := lookup[uid]
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
