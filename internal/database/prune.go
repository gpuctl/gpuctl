package database

import (
	"time"
)

func DownsampleOverTime(interval int, database Database) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for t := range downsampleTicker.C {
		err := downsampleDatabase(database, t)

		if err != nil {
			return err
		}
	}

	return nil
}

func downsampleDatabase(database Database, t time.Time) error {
	return database.Downsample(t.Unix())
}
