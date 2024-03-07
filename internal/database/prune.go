package database

import (
	"log/slog"
	"time"
)

func DownsampleOverTime(interval time.Duration, database Database) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for range downsampleTicker.C {
		// XXX: emergency fix to get some type of downsampling if we dont get this working
		//err := database.Downsample(time.Now())
		err := database.Downsample(time.Now().AddDate(0, 0, -7))

		if err != nil {
			slog.Error("Got error whilst downsampling", "err", err)
		}
	}

	return nil
}

func downsampleDatabase(database Database, t time.Time) error {
	return database.Downsample(t)
}
