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
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func handler(w http.ResponseWriter, _ *http.Request) {
	signal, _ := dsp.ReadSignalFile("examples/signals/example_signal_31_hz.txt", 31)

	signal, err := signal.Normalize()

	if err != nil {
		panic(err)
	}

	// In our example we only care about resting heart rates zones
	// 0.5Hz => 30 BPM
	// 1.2Hz => 70 BPM
	signal, err = signal.BandPassFilter(0.5, 1.2)

	if err != nil {
		panic(err)
	}

	// Detect the r-peaks in the signal
	rPeaks := dsp.GetRPeaks(signal)

	println(rPeaks.String())

	plot1 := plotSignal(signal, rPeaks, opts.Title{
		Title:    "R-Peak detection in signal",
		Subtitle: signal.String(),
	})

	plot2 := plotRInterval(rPeaks, opts.Title{
		Title:    "R interval time",
		Subtitle: signal.String(),
	})

	plot3 := plotBeatsPerMinute(rPeaks, opts.Title{
		Title:    "Heart beats per minute",
		Subtitle: signal.String(),
	})

	page := components.NewPage()
	page.AddCharts(plot1, plot2, plot3)

	page.Render(w)
}

func plotSignal(signal *dsp.Signal, rPeaks dsp.RPeaks, title opts.Title) *charts.Line {
	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	y2 := make([]opts.LineData, 0)
	for i := 0; i < signal.Length(); i++ {
		x = append(x, fmt.Sprintf("%.1f", float64(i)/signal.SampleRate))
		y = append(y, opts.LineData{Value: signal.Signal[i], Symbol: "none"})

		// Check if this is an R-Peak
		if rPeaks.IsRPeak(i) {
			y2 = append(y2, opts.LineData{Value: signal.Signal[i], Symbol: "circle"})
		} else {
			y2 = append(y2, opts.LineData{Value: 0, Symbol: "none"})
		}
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWonderland}),
		charts.WithTitleOpts(title))

	line.SetXAxis(x).AddSeries("data", y).
		AddSeries("rpeaks", y2, charts.WithLineStyleOpts(opts.LineStyle{Color: "red"}))

	return line
}

func plotRInterval(rPeaks dsp.RPeaks, title opts.Title) *charts.Line {
	x := make([]string, 0)
	y := make([]opts.LineData, 0)

	for i := 0; i < len(rPeaks.RPeakInterval); i++ {
		x = append(x, fmt.Sprintf("%v", i))
		y = append(y, opts.LineData{Value: rPeaks.RPeakInterval[i], Symbol: "none"})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWonderland}),
		charts.WithTitleOpts(title))

	line.SetXAxis(x).AddSeries("data", y)

	return line
}

func plotBeatsPerMinute(rPeaks dsp.RPeaks, title opts.Title) *charts.Line {
	x := make([]string, 0)
	y := make([]opts.LineData, 0)

	for i := 0; i < len(rPeaks.HeartBeatsPerMinute); i++ {
		x = append(x, fmt.Sprintf("%v", i))
		y = append(y, opts.LineData{Value: rPeaks.HeartBeatsPerMinute[i], Symbol: "none"})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWonderland}),
		charts.WithTitleOpts(title))

	line.SetXAxis(x).AddSeries("data", y)

	return line
}
