package gpustats

import (
	"errors"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

var (
	SampleAdditionError = errors.New("two packets with different contexts cannot be aggregated using Add, consider using UncontextualAdd")
)

// Combine two GPUStatSamples instances into one.
func Add(l, r uplink.GPUStatSample) (uplink.GPUStatSample, error) {
	if l.Uuid != r.Uuid {
		return uplink.GPUStatSample{}, SampleAdditionError
	}

	return uplink.GPUStatSample{
		MemoryUtilisation: l.MemoryUtilisation + r.MemoryUtilisation,
		GPUUtilisation:    l.GPUUtilisation + r.GPUUtilisation,
		FanSpeed:          l.FanSpeed + r.FanSpeed,
		Temp:              l.Temp + r.Temp,
		MemoryUsed:        l.MemoryUsed + r.MemoryUsed,
		MemoryTemp:        l.MemoryTemp + r.MemoryTemp,
		GraphicsVoltage:   l.GraphicsVoltage + r.GraphicsVoltage,
		PowerDraw:         l.PowerDraw + r.PowerDraw,
		GraphicsClock:     l.GraphicsClock + r.GraphicsClock,
		MaxGraphicsClock:  l.MaxGraphicsClock + r.MaxGraphicsClock,
		MemoryClock:       l.MemoryClock + r.MemoryClock,
		MaxMemoryClock:    l.MaxMemoryClock + r.MaxMemoryClock,
	}, nil
}

// Scale each value in s by scalar
func Scale(s uplink.GPUStatSample, scalar float64) uplink.GPUStatSample {
	return uplink.GPUStatSample{
		MemoryUtilisation: s.MemoryUtilisation * scalar,
		GPUUtilisation:    s.GPUUtilisation * scalar,
		FanSpeed:          s.FanSpeed * scalar,
		Temp:              s.Temp * scalar,
		MemoryUsed:        s.MemoryUsed * scalar,
		MemoryTemp:        s.MemoryTemp * scalar,
		GraphicsVoltage:   s.GraphicsVoltage * scalar,
		PowerDraw:         s.PowerDraw * scalar,
		GraphicsClock:     s.GraphicsClock * scalar,
		MaxGraphicsClock:  s.MaxGraphicsClock * scalar,
		MemoryClock:       s.MemoryClock * scalar,
		MaxMemoryClock:    s.MaxMemoryClock * scalar,
	}
}
