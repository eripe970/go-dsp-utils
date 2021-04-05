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

func handler(w http.ResponseWriter, _ *http.Request) {
	// 1. Read the initial signal
	signal, _ := dsp.ReadSignalFile("examples/signals/example_signal_31_hz.txt", 31)

	plot1 := plotSignal(signal, opts.Title{
		Title:    "Original signal",
		Subtitle: signal.String(),
	})

	// 2. Sample the signal to 10s
	signal10s := signal.Sample(10 * time.Second)

	plot2 := plotSignal(signal10s, opts.Title{
		Title:    "Original signal - 10 seconds sample",
		Subtitle: signal10s.String(),
	})

	// 3. Normalize the signal around -1 and 1
	signalNormalized, err := signal.Normalize()

	if err != nil {
		panic(err)
	}

	plot3 := plotSignal(signalNormalized, opts.Title{
		Title:    "Normalized signal around -1 and 1",
		Subtitle: signalNormalized.String(),
	})

	// 4. Frequency spectrum of the signal
	frequencySpectrum, err := signalNormalized.FrequencySpectrum()

	if err != nil {
		panic(err)
	}

	plot4 := plotSpectrum(frequencySpectrum, opts.Title{
		Title:    "Frequency spectrum",
		Subtitle: frequencySpectrum.String(),
	})

	// 5. Filter the signal
	// 0.5Hz => 30 BPM
	// 1.2Hz => 70 BPM
	// In our example we only care about resting heart rates zones
	signalFiltered, err := signalNormalized.BandPassFilter(0.5, 1.2)

	if err != nil {
		panic(err)
	}

	signalFiltered, err = signalFiltered.Normalize()

	if err != nil {
		panic(err)
	}

	plot5 := plotSignal(signalFiltered, opts.Title{
		Title:    "Signal after low and high pass filters",
		Subtitle: signalFiltered.String(),
	})

	page := components.NewPage()
	page.AddCharts(
		plot1,
		plot2,
		plot3,
		plot4,
		plot5,
	)

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
