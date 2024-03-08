package database

import (
	"log/slog"
	"time"
)

func DownsampleOverTime(interval time.Duration, downsample_type string, database Database) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for range downsampleTicker.C {
		if downsample_type == "DOWNSAMPLE" {
			err := database.Downsample(time.Now())

			if err != nil {
				slog.Error("Got error whilst downsampling", "err", err)
			}
		} else if downsample_type == "DELETE" {
			err := database.Delete(time.Now())

			if err != nil {
				slog.Error("Got error whilst deleting", "err", err)
			}
		} else {
			slog.Error("Unknown downsample type", "type", downsample_type)
		}
	}

	return nil
}

func downsampleDatabase(database Database, t time.Time) error {
	return database.Downsample(t)
}
