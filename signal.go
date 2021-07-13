package dsp

import (
	"errors"
	"fmt"
	"time"
)

type Signal struct {
	// SampleRate is the sampling rate of the signal in Hz
	SampleRate float64
	// Store the signal values
	Signal []float64
}

func (s *Signal) Sample(duration time.Duration) *Signal {

	sample := make([]float64, 0)

	for i := 0; i < int(s.SampleRate*duration.Seconds()); i++ {
		sample = append(sample, s.Signal[i])
	}

	return &Signal{
		SampleRate: s.SampleRate,
		Signal:     sample,
	}
}

func (s *Signal) Split(split time.Duration) []*Signal {
	sample := make([]float64, 0)

	result := make([]*Signal, 0)

	size := int(s.SampleRate * split.Seconds())

	for i := 0; i < len(s.Signal); i++ {
		sample = append(sample, s.Signal[i])

		if i > 0 && i%size == 0 {
			result = append(result, &Signal{
				SampleRate: s.SampleRate,
				Signal:     sample,
			})
			sample = make([]float64, 0)
		}
	}

	if len(sample) > 0 {
		result = append(result, &Signal{
			SampleRate: s.SampleRate,
			Signal:     sample,
		})
	}

	return result
}

func (s *Signal) String() string {
	return fmt.Sprintf("SampleRate: %vHz, Length: %v, Duration: %.1fs", s.SampleRate, len(s.Signal), s.Duration())
}

func (s *Signal) Length() int {
	return len(s.Signal)
}

func (s *Signal) Duration() float64 {
	duration := float64(len(s.Signal)) / s.SampleRate
	return duration
}

func (s *Signal) Min() float64 {
	if len(s.Signal) == 0 {
		return 0
	}

	min := s.Signal[0]

	for i := range s.Signal {
		if s.Signal[i] < min {
			min = s.Signal[i]
		}
	}

	return min
}

func (s *Signal) Max() float64 {
	if len(s.Signal) == 0 {
		return 0
	}

	max := s.Signal[0]

	for i := range s.Signal {
		if s.Signal[i] > max {
			max = s.Signal[i]
		}
	}

	return max
}

// Normalize the signal between -1 and 1 and return a new signal
func (s *Signal) Normalize() (*Signal, error) {
	input := s.Signal

	if len(input) == 0 {
		return &Signal{
			SampleRate: s.SampleRate,
			Signal:     s.Signal,
		}, nil
	}

	min := s.Min()
	max := s.Max()

	// We can't normalize a flat signal where min == max
	if min == max {
		return nil, errors.New("cannot normalize signal")
	}

	normalized := make([]float64, len(input))

	for i := range input {
		val := input[i]

		normalized[i] = 2*((val-min)/(max-min)) - 1
	}

	return &Signal{
		SampleRate: s.SampleRate,
		Signal:     normalized,
	}, nil
}
