package main

import (
	"fmt"
	"github.com/eripe970/go-dsp-utils"
	"time"
)

func main() {
	// Read a signal sampled at 31hz
	signal1, _ := dsp.ReadSignalFile("examples/signals/example_signal_31_hz.txt", 31)
	fmt.Println(signal1)

	// Read a signal sampled at 100Hz
	signal2, _ := dsp.ReadSignalFile("examples/signals/example_signal_100_hz.txt", 100)
	fmt.Println(signal2)

	// Get a 10 second sample of the signal
	signal10s := signal1.Sample(10 * time.Second)

	fmt.Println(signal10s)

	// Normalize the signal between -1 and 1
	normalized, _ := signal1.Normalize()

	// Calculate the frequency spectrum of the signal (FFT + massage of the numbers)
	spectrum, _ := normalized.FrequencySpectrum()

	fmt.Println(spectrum)

	// Run some filters on the signal
	_, _ = signal10s.LowPassFilter(3)
	_, _ = signal10s.HighPassFilter(10)
	_, _ = signal10s.BandPassFilter(3, 10)
}
