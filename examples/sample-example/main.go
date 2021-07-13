package main

import (
	"fmt"
	"github.com/eripe970/go-dsp-utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

// Example that detects how much noisy there is on a signal
func handler(w http.ResponseWriter, _ *http.Request) {
	signal, _ := dsp.ReadSignalFile("examples/signals/example-noisy-signal_31_hz.txt", 31)

	plot1 := plotSignal(signal, opts.Title{
		Title:    "Original signal",
		Subtitle: signal.String(),
	})

	page := components.NewPage()
	page.AddCharts(plot1)

	split := 10
	parts := signal.Split(time.Duration(split) * time.Second)

	duration := 0.0

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		from := i * split
		to := from + int(part.Duration())

		page.AddCharts(plotSignal(part, opts.Title{
			Title:    fmt.Sprintf("Signal %vs to %vs", i*split, (i+1)*split),
			Subtitle: part.String(),
		}))

		// Normalize the signal around -1 and 1
		signalNormalized, err := part.Normalize()

		if err != nil {
			panic(err)
		}

		page.AddCharts(plotSignal(signalNormalized, opts.Title{
			Title:    fmt.Sprintf("Normalized signal %vs to %vs", from, to),
			Subtitle: signalNormalized.String(),
		}))

		// 4. Frequency spectrum of the signal
		frequencySpectrum, err := signalNormalized.FrequencySpectrum()

		if err != nil {
			panic(err)
		}

		page.AddCharts(plotSpectrum(frequencySpectrum, opts.Title{
			Title:    fmt.Sprintf("Frequency spectrum %vs to %vs", from, to),
			Subtitle: frequencySpectrum.String(),
		}))

		// Some example calculations
		for j, frequency := range frequencySpectrum.Frequencies {
			// Only care frequencies above 2hz
			if frequency > 2 {
				if frequencySpectrum.Spectrum[j] > 0.3 {
					// More than 30%, it's considered noisy
					duration += part.Duration()
					break
				}
			}
		}
	}

	fmt.Printf("noisy %vs, percentage %.1f%%\n", duration, 100*duration/signal.Duration())

	page.Render(w)
}

func plotSignal(signal *dsp.Signal, title opts.Title) *charts.Line {
	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	for i := 0; i < signal.Length(); i++ {
		x = append(x, fmt.Sprintf("%.1f", float64(i)/signal.SampleRate))
		y = append(y, opts.LineData{Value: signal.Signal[i], Symbol: "none"})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWonderland}),
		charts.WithYAxisOpts(opts.YAxis{Max: signal.Max(), Min: signal.Min(), SplitNumber: 10}),
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 10}),
		charts.WithTitleOpts(title))

	line.SetXAxis(x).AddSeries("data", y).SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
	)

	return line
}

func plotSpectrum(spectrum *dsp.FrequencySpectrum, title opts.Title) *charts.Line {
	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	for i := 0; i < spectrum.Length(); i++ {
		x = append(x, fmt.Sprintf("%.1f", spectrum.Frequencies[i]))
		y = append(y, opts.LineData{Value: spectrum.Spectrum[i], Symbol: "none"})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWonderland}),
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 10}),
		charts.WithTitleOpts(title))

	line.SetXAxis(x).AddSeries("data", y).SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
	)

	return line
}
