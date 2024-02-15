// package passwd provides a parser for unix /etc/passwd files.
//
// For details of the file format, see "man 5 passwd" or:
//   - https://man7.org/linux/man-pages/man5/passwd.5.html
//   - https://manpages.ubuntu.com/manpages/noble/en/man5/passwd.5.html
package passwd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Entry struct {
	Name     string
	Password string // The users **login** name

	Uid uint32 // User ID
	Gid uint32 // Group ID

	// General Electric Comprehensive Operating System field
	//
	// Use the [Entry.ParseGecos] method to get details.
	//
	// See also https://en.wikipedia.org/wiki/Gecos_field
	Gecos string

	HomeDir string // User's ~/$HOME directory
	Shell   string // User's login shell
}

var ErrBadSyntax = errors.New("passwd: invalid syntax")

type Passwd []Entry

func Parse(contents io.Reader) (Passwd, error) {

	var ents Passwd

	b := bufio.NewScanner(contents)

	for b.Scan() {
		ent, err := parseEntry(b.Text())

		if err != nil {
			return nil, err
		}

		ents = append(ents, ent)
	}

	if err := b.Err(); err != nil {
		return nil, err
	}

	return ents, nil
}

func parseEntry(line string) (Entry, error) {
	parts := strings.Split(line, ":")
	var e Entry

	if len(parts) != 7 {
		return e, fmt.Errorf("%w: line `%s` doesn't have 7 parts", ErrBadSyntax, line)
	}

	uid, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return e, fmt.Errorf("%w: failed to parse uid `%s`: %w", ErrBadSyntax, parts[2], err)
	}

	gid, err := strconv.ParseUint(parts[3], 10, 32)
	if err != nil {
		return e, fmt.Errorf("%w: failed to parse gid `%s`: %w", ErrBadSyntax, parts[3], err)
	}

	e.Name = parts[0]
	e.Password = parts[1]
	e.Uid = uint32(uid) // 2
	e.Gid = uint32(gid) // 3
	e.Gecos = parts[4]
	e.HomeDir = parts[5]
	e.Shell = parts[6]

	return e, nil
}

func (e Entry) ParseGecos() Gecos {
	var r Gecos

	fields := strings.Split(e.Gecos, ",")

	// finger does some processing, TODO: do we need to, or is that only for
	// old insane systems.
	// https://github.com/Distrotech/bsd-finger/blob/dae7f2836b01d7812a32f9dc189c0862f60edff7/finger/util.c#L97-L156

	/*
	 * fields[0] -> real name
	 * fields[1] -> office
	 * fields[2] -> officephone
	 * fields[3] -> homephone
	 */
	switch len(fields) {
	default:
		r.HomePhone = fields[3]
		fallthrough
	case 3:
		r.OfficePhone = fields[2]
		fallthrough
	case 2:
		r.Office = fields[1]
		fallthrough
	case 1:
		r.FullName = fields[0]
	case 0:
	}

	return r
}

type Gecos struct {
	FullName    string
	Office      string
	OfficePhone string
	HomePhone   string
}
