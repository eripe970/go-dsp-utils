package dsp

import (
	"github.com/mattetti/audio"
	"github.com/mattetti/audio/transforms/filters"
)

func (s *Signal) LowPassFilter(cutOffFreq float64) (*Signal, error) {
	buffer := audio.NewPCMFloatBuffer(s.Signal, &audio.Format{SampleRate: int(s.SampleRate)})

	err := filters.LowPass(buffer, cutOffFreq)

	if err != nil {
		return nil, err
	}

	return &Signal{
		SampleRate: s.SampleRate,
		Signal:     buffer.Floats,
	}, nil
}

func (s *Signal) HighPassFilter(cutOffFreq float64) (*Signal, error) {
	buffer := audio.NewPCMFloatBuffer(s.Signal, &audio.Format{SampleRate: int(s.SampleRate)})

	err := filters.HighPass(buffer, cutOffFreq)

	if err != nil {
		return nil, err
	}

	return &Signal{
		SampleRate: s.SampleRate,
		Signal:     buffer.Floats,
	}, nil
}

func (s *Signal) BandPassFilter(lower, upper float64) (*Signal, error) {
	signal1, err := s.LowPassFilter(upper)

	if err != nil {
		return nil, err
	}

	signal2, err := signal1.HighPassFilter(lower)

	if err != nil {
		return nil, err
	}

	return signal2, nil
}
