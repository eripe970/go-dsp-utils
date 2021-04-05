# GO-DSP-UTILS

go-dsp-utils is a digital signal processing package for the [Go programming language](http://golang.org). It's a wrapper
around some of the most popular go-lang packages for digital signal processing. The purpose is to make the life easier
when working with digital signal processing in golang.

#### The core package

- Read signal files
- Sample signal files
- Normalize signals
- Calculate the frequency spectrum (FFT + some spectrum logic)
- Low and high pass filtering on a signal

#### The heart beat package

The heart beat package is a package that can be used to detect heart beats (time between r-peaks) in a signal.

- Calculate R-peaks of a heart beat signal

## DSP Packages used

* **[go-dsp](https://github.com/mjibson/go-dsp)** - DSP package for go
* **[audio](https://github.com/mattetti/audio)** - DSP package as well as low and high pass filters
* **[godsp](https://github.com/goccmack/godsp)** - DSP packages for go

## Graph packages used

* **[go-echarts](https://github.com/go-echarts/go-echarts)** - A nice charts library for Golang.

## Details and background

There is a code walk through with a lot of examples at medium <to add>.

Example signal from http://www.paulvangent.com/

## Installation and Usage

```$ go get github.com/eripe970/go-dsp-utils```

### Example

Example program for working with signals (see examples/basic).

```
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
	
	// Calculate the frequency spectrum of the signal
	spectrum, _ := normalized.FrequencySpectrum()

	fmt.Println(spectrum)
	
	// Run some filters on the signal
	_, _ = signal10s.LowPassFilter(3)
	_, _ = signal10s.HighPassFilter(10)
	_, _ = signal10s.BandPassFilter(3, 10)
}
```

Output:

```
SampleRate: 31Hz, Length: 1577, Duration: 50.9s
SampleRate: 100Hz, Length: 2483, Duration: 24.8s
SampleRate: 31Hz, Length: 310, Duration: 10.0s
Length: 788, Spectrum: 0Hz - 15.5Hz
```

### R-peak example

Example program for detecting heart rate (r-peaks).

```
package main

import (
	"fmt"
	"github.com/eripe970/go-dsp-utils"
)

func main() {
    signal, _ := dsp.ReadSignalFile("examples/signals/example_signal_31_hz.txt", 31)
    
    // Detect the r-peaks in the signal
    rPeaks := dsp.GetRPeaks(signal)
    
    println(rPeaks.String())	
}
```

Output:

```
Total heart beats detected: 46, avg heart rate: 55 BPM
```