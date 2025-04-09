package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

func toFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func Slice(x []float64, begin, end float64) ([]float64, int, int) {
	startIdx := 0
	endIdx := len(x)

	if begin > 0 {
		for i, v := range x {
			if v >= begin {
				startIdx = i
				break
			}
		}
	}

	if end < 1<<32 {
		for i := startIdx; i < len(x); i++ {
			if x[i] > end {
				endIdx = i
				break
			}
		}
	}

	return x[startIdx:endIdx], startIdx, endIdx
}

func loadCSV(path string) (axisX, axisY, asixXTemp, axitYTemp []float64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	_, err = reader.Read() // Пропустим заголовок
	if err != nil {
		return
	}

	for {
		record, errRead := reader.Read()
		if errRead == io.EOF {
			break
		}
		// if errRead != nil {
		// 	err = errRead
		// 	continue
		// }

		if len(record) < 3 {
			continue
		}

		axisX = append(axisX, toFloat(record[1]))
		axisY = append(axisY, toFloat(record[2]))

		if len(record) == 4 {
			asixXTemp = append(asixXTemp, toFloat(record[1]))
			axitYTemp = append(axitYTemp, toFloat(record[3]))
		}
	}

	return
}
