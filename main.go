package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

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
		return
	}

	fmt.Println("SSS", toFloat("-0,199"))
	fmt.Println("SSS", toFloat("-0.199"))
	fmt.Println("SSS", toFloat("  -0,199"))

	file := args[0]
	begin, end := 0.0, math.MaxFloat64
	if len(args) > 1 {
		begin, _ = strconv.ParseFloat(args[1], 64)
	}
	if len(args) > 2 {
		end, _ = strconv.ParseFloat(args[2], 64)
	}

	fmt.Println("Используем файл", file)
	_, _, x, y, err := loadCSV(file)
	if err != nil {
		log.Fatalf("Ошибка чтения CSV: %v\n", err)
	}
	fmt.Printf("Количество записей: %d, %.2f секунд\n", len(x), x[len(x)-1])
	_, startIdx, endIdx := Slice(x, begin, end)
	y = y[startIdx:endIdx]

	var sum, sumSq float64
	for _, v := range y {
		sum += v
		sumSq += v * v
	}
	mean := sum / float64(len(y))
	variance := sumSq/float64(len(y)) - mean*mean
	stdDev := math.Sqrt(variance)
	fmt.Printf("Temperature:\n  Mean: %.3f C\n  Std: %.3f K\n  Var: %.3f K^2\n", mean, stdDev, variance)
}
