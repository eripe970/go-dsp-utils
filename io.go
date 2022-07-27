package dsp

import (
	"bufio"
	"os"
	"strconv"
)

func ReadSignalFile(path string, sampleRate float64) (*Signal, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	signal := Signal{
		SampleRate: sampleRate,
		Signal:     make([]float64, 0),
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		v, err := strconv.ParseFloat(scanner.Text(), 64)

		if err != nil {
			return nil, err
		}

		signal.Signal = append(signal.Signal, v)
	}
	return &signal, scanner.Err()
}
// Function allows for reading
func ReadArray(dataArray []float64, sampleRate float64) (*Signal, error){

	signal := Signal{
		SampleRate: sampleRate,
		Signal:     make([]float64, 0),
	}

	for i:=0; i < len(dataArray); i++ {
		v:=  dataArray[i]
		signal.Signal = append(signal.Signal, v)
	}
	return &signal, nil
}