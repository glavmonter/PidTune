package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func loadCSV(path string) ([]float64, []float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var x, y []float64
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		xv, _ := strconv.ParseFloat(line[0], 64)
		yv, _ := strconv.ParseFloat(line[1], 64)
		x = append(x, xv)
		y = append(y, yv)
	}

	return x, y, nil
}

func main() {
	var frequecy float64
	flag.Float64Var(&frequecy, "frequency", 50.0, "Sampling frequency")
	flag.Float64Var(&frequecy, "f", 50, "Sampling frequency (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file> [begin] [end]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Файл должен быть задан")
		flag.Usage()
	}
}
