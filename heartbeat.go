package dsp

import (
	"fmt"
	"github.com/goccmack/godsp/peaks"
)

type RPeaks struct {
	Indices             []int
	RPeakInterval       []float64
	HeartBeatsPerMinute []float64
}

func GetRPeaks(signal *Signal) RPeaks {
	// Assume that the highest heart rate is 220 BPM
	separator := signal.SampleRate / (220.0 / 60.0)

	indices := peaks.Get(signal.Signal, int(separator))

	rPeakInterval := make([]float64, 0)
	bpm := make([]float64, 0)

	previousTime := 0.0

	// Calculate the r-peak interval
	for i := range indices {
		currentTime := float64(indices[i]) * (1.0 / signal.SampleRate)

		t := currentTime - previousTime
		rPeakInterval = append(rPeakInterval, t)
		if t == 0 {
			bpm = append(bpm, 0)
		} else {
			bpm = append(bpm, 60/t)
		}

		previousTime = currentTime
	}

	return RPeaks{
		Indices:             indices,
		RPeakInterval:       rPeakInterval,
		HeartBeatsPerMinute: bpm,
	}
}

func (r *RPeaks) IsRPeak(index int) bool {
	for i := range r.Indices {
		if r.Indices[i] == index {
			return true
		}
	}

	return false
}

func (r *RPeaks) Avg() int {
	sum := 0.0
	for i := range r.HeartBeatsPerMinute {
		sum += r.HeartBeatsPerMinute[i]
	}

	return int(sum / float64(len(r.HeartBeatsPerMinute)))
}

func (r *RPeaks) Count() int {
	return len(r.HeartBeatsPerMinute)
}

func (r *RPeaks) String() string {
	return fmt.Sprintf("Total heart beats detected: %v, avg heart rate: %v BPM", r.Count(), r.Avg())
}
