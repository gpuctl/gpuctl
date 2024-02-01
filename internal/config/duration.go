package config

import (
	"time"
)

// Duration is a wrapper over time.Duration that can be (un)marshalled from toml.
type Duration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	s := string(text)

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(dur)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (d *Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

/* Forwarding methods
TODO: Implement all of these

alona@Ashtabula:/usr/lib/go-1.21/api$ rg "pkg time, method \\(*?Duration"
go1.13.txt
8023:pkg time, method (Duration) Microseconds() int64
8024:pkg time, method (Duration) Milliseconds() int64

go1.txt
30564:pkg time, method (Duration) Hours() float64
30565:pkg time, method (Duration) Minutes() float64
30566:pkg time, method (Duration) Nanoseconds() int64
30567:pkg time, method (Duration) Seconds() float64
30568:pkg time, method (Duration) String() string

go1.9.txt
195:pkg time, method (Duration) Round(Duration) Duration
196:pkg time, method (Duration) Truncate(Duration) Duration

go1.19.txt
291:pkg time, method (Duration) Abs() Duration #51414
*/

func (d Duration) String() string {
	return time.Duration(d).String()
}

const (
	Nanosecond  Duration = Duration(time.Nanosecond)
	Microsecond Duration = Duration(time.Microsecond)
	Millisecond Duration = Duration(time.Millisecond)
	Second      Duration = Duration(time.Second)
	Minute      Duration = Duration(time.Minute)
	Hour        Duration = Duration(time.Hour)
)
