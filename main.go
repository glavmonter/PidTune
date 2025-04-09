package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/mjibson/go-dsp/fft"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
	fmt.Printf("Temperature:\n  Mean: %.3f C\n  Std: %.4f K\n  Var: %.4f K^2\n", mean, stdDev, variance)

	// FFT
	fftResult := fft.FFTReal(y)
	n := len(y)
	amplitudes := make([]float64, n/2)
	frequencies := make([]float64, n/2)
	for i := 1; i < n/2; i++ {
		amplitudes[i] = cmplAbs(fftResult[i]) / float64(n)
		frequencies[i] = float64(i) * frequecy / float64(n)
	}

	maxIndex := 1
	for i := 2; i < len(amplitudes); i++ {
		if amplitudes[i] > amplitudes[maxIndex] {
			maxIndex = i
		}
	}

	oscFreq := frequencies[maxIndex]
	period := 1 / oscFreq
	fmt.Printf("Oscillation frequency %.4f Hz, %.4f seconds\n", oscFreq, period)

	endAmpIndex := int(math.Min(float64(len(y)), frequecy*period*2.5))
	ampSlice := y[:endAmpIndex]
	min, max := ampSlice[0], ampSlice[0]
	for _, v := range ampSlice {
		if v < min {
			min = v
		}

		if v > max {
			max = v
		}
	}

	amp := max - min
	fmt.Printf("Maximum: %.3f\nMinimum: %.3f\nAmplitude: %.3f\n", max, min, amp)

	err = plotSignal(x, y, "original_signal.png", "Original Signal", "Time (s)", "Amplitude")
	if err != nil {
		log.Fatalf("Ошибка построения временного графика: %v\n", err)
	}

	// Частотная область (без нулевой частоты)
	freqData := make(plotter.XYs, len(frequencies)-1)
	for i := 1; i < len(frequencies); i++ {
		freqData[i-1].X = frequencies[i]
		freqData[i-1].Y = amplitudes[i]
	}
	err = plotXY(freqData, "fft_plot.png", "FFT of Signal", "Frequency (Hz)", "Amplitude")
	if err != nil {
		log.Fatalf("Ошибка построения спектра: %v\n", err)
	}
}

func cmplAbs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

func plotSignal(x, y []float64, filename, title, xlabel, ylabel string) error {
	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return plotXY(pts, filename, title, xlabel, ylabel)
}

func plotXY(pts plotter.XYs, filename, title, xlabel, ylabel string) error {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}

	p.Add(line)
	p.Add(plotter.NewGrid())

	return p.Save(10*vg.Inch, 4*vg.Inch, filename)
}
