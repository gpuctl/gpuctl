package database

import (
	"log/slog"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
)

func DownsampleOverTime(interval config.Duration, database Database) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for range downsampleTicker.C {
		err := database.Downsample(time.Now())

		if err != nil {
			slog.Error("Got error whilst downsampling", "err", err)
		}
	}

	return nil
}

func downsampleDatabase(database Database, t time.Time) error {
	return database.Downsample(t)
}
