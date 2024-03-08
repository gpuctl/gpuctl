package database

import (
	"time"
)

func DownsampleOverTime(interval time.Duration, downsample_thresh time.Duration, delete_thresh time.Duration, database Database) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for range downsampleTicker.C {
		database.Downsample(downsample_thresh)
		database.Delete(delete_thresh)
	}

	return nil
}

func downsampleDatabase(database Database, t time.Duration) error {
	return database.Downsample(t)
}
