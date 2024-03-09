package database

import (
	"time"
)

type PruneConfig struct {
	Interval            time.Duration
	DeleteThreshold     time.Duration
	DownsampleThreshold time.Duration
}

func DownsampleOverTime(downsampleConfig PruneConfig, database Database) error {
	interval := downsampleConfig.Interval
	delete_thresh := downsampleConfig.DeleteThreshold
	downsample_thresh := downsampleConfig.DownsampleThreshold

	downsampleTicker := time.NewTicker(time.Duration(interval))

	for range downsampleTicker.C {
		err := database.Downsample(-downsample_thresh)

		if err != nil {
			return err
		}

		err = database.DeleteOldStats(-delete_thresh)

		if err != nil {
			return err
		}
	}

	return nil
}

func downsampleDatabase(database Database, t time.Duration) error {
	return database.Downsample(t)
}
