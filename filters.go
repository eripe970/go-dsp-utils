package dsp

import (
	"github.com/mattetti/audio"
	"github.com/mattetti/audio/transforms/filters"
)

// fakedFreq returns faked cutoff frequency and sample rate if the given
// sample rate is less than 0.
// it is a workaround for https://github.com/eripe970/go-dsp-utils/issues/4
func fakedFreq(cutOffFreq, sampleRate float64) (float64, int) {
	sr := int(sampleRate)
	m := 1.0
	for sr == 0 {
		m = m * 1000
		cutOffFreq = cutOffFreq * 1000
		sr = int(sampleRate * m)
	}
	return cutOffFreq, sr
}

func (s *Signal) LowPassFilter(cutOffFreq float64) (*Signal, error) {
	cutOffFreq, sampleRate := fakedFreq(cutOffFreq, s.SampleRate)
	buffer := audio.NewPCMFloatBuffer(s.Signal, &audio.Format{SampleRate: sampleRate})

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
	cutOffFreq, sampleRate := fakedFreq(cutOffFreq, s.SampleRate)
	buffer := audio.NewPCMFloatBuffer(s.Signal, &audio.Format{SampleRate: sampleRate})

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
