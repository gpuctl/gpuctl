package config

import "time"

type Duration = time.Duration

const (
	Nanosecond  Duration = Duration(time.Nanosecond)
	Microsecond Duration = Duration(time.Microsecond)
	Millisecond Duration = Duration(time.Millisecond)
	Second      Duration = Duration(time.Second)
	Minute      Duration = Duration(time.Minute)
	Hour        Duration = Duration(time.Hour)
)
