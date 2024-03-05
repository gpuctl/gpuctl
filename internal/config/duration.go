package config

import "time"

type Duration time.Duration

const (
	Nanosecond  Duration = Duration(time.Nanosecond)
	Microsecond Duration = Duration(time.Microsecond)
	Millisecond Duration = Duration(time.Millisecond)
	Second      Duration = Duration(time.Second)
	Minute      Duration = Duration(time.Minute)
	Hour        Duration = Duration(time.Hour)
)

func (d *Duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))

	if err != nil {
		return err
	}

	*d = Duration(dur)
	return nil
}

// NewTicker returns a new Ticker containing a channel that will send
// the current time on the channel after each tick. The period of the
// ticks is specified by the duration argument. The ticker will adjust
// the time interval or drop ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will
// panic. Stop the ticker to release associated resources.
func (d Duration) NewTicker() *time.Ticker {
	return time.NewTicker(time.Duration(d))
}
