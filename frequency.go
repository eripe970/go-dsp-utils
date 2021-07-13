package dsp

import (
	"fmt"
	fft2 "github.com/mjibson/go-dsp/fft"
	"math/cmplx"
)

type FrequencySpectrum struct {
	// Store the frequency spectrum values
	Spectrum []float64

	// Store the frequencies used
	Frequencies []float64
}

func (s *FrequencySpectrum) String() string {
	if s.Spectrum != nil && s.Frequencies != nil {
		return fmt.Sprintf("Length: %v, Spectrum: 0Hz - %.1fHz", len(s.Spectrum), s.Frequencies[len(s.Frequencies)-1])
	}
	return "No spectrum or frequencies"
}

// FrequencySpectrum calculates the frequency spectrum of the signal.
func (s *Signal) FrequencySpectrum() (*FrequencySpectrum, error) {
	// Apply the FFT on the signal to get the frequency components
	fft := fft2.FFTReal(s.Signal)

	// Compute the two sided spectrum
	var spectrum2 []float64

	length := float64(len(s.Signal))

	for i := range fft {
		// Get the absolute value since the fft is a complex value with both real and imaginary parts
		spectrum2 = append(spectrum2, cmplx.Abs(fft[i])/length)
	}

	// Look at the one sided spectrum
	spectrum1 := spectrum2[0 : len(spectrum2)/2]

	for i := range spectrum1 {
		spectrum1[i] = spectrum1[i] * 2
	}

	var frequencies []float64

	spectrumLength := float64(len(spectrum2))

	for i := range spectrum1 {
		// Calculate which frequencies the spectrum contains
		frequencies = append(frequencies, float64(i)*s.SampleRate/spectrumLength)
	}

	return &FrequencySpectrum{
		Spectrum:    spectrum1,
		Frequencies: frequencies,
	}, nil
}

func (s *FrequencySpectrum) Length() int {
	return len(s.Spectrum)
}

func (s *FrequencySpectrum) Min() float64 {
	if len(s.Spectrum) == 0 {
		return 0
	}

	min := s.Spectrum[0]

	for i := range s.Spectrum {
		if s.Spectrum[i] < min {
			min = s.Spectrum[i]
		}
	}

	return min
}

func (s *FrequencySpectrum) Max() float64 {
	if len(s.Spectrum) == 0 {
		return 0
	}

	max := s.Spectrum[0]

	for i := range s.Spectrum {
		if s.Spectrum[i] > max {
			max = s.Spectrum[i]
		}
	}

	return max
}
