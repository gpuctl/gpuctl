package database

import (
	"time"
)

type PruneConfig struct {
	interval           time.Duration
	deleteThreshold    time.Duration
	downsampeThreshold time.Duration
}

func DownsampleOverTime(downsampleConfig PruneConfig, database Database) error {
	interval := downsampleConfig.interval
	delete_thresh := downsampleConfig.deleteThreshold
	downsample_thresh := downsampleConfig.downsampeThreshold

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
