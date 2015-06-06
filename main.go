package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Rate     float64
	Duration time.Duration
	Host     string
	Port     int64
	Endpoint string
}

func main() {
	config := getConfig()
	requestor := NewRequestor(config)
	results, err := requestor.MakeRequest()
	if err != nil {
		fmt.Println(err)
		return
	}

	report(results)
}

func getConfig() *Config {
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.Int64("port", 8080, "The port of the host server")
	endpoint := flag.String("endpoint", "/", "Then endpoint on the server")
	rate := flag.Float64("rate", 1.0, "Requests per second")
	duration := flag.Duration("duration", 1*time.Second, "Duration of the test")

	flag.Parse()

	return &Config{
		Rate:     *rate,
		Duration: *duration,
		Host:     *host,
		Port:     *port,
		Endpoint: *endpoint,
	}
}

func report(results []Result) {
	calc := NewMetrics(results)
	calc.Sort()

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Requests:", calc.Total)

	fmt.Println("Latencies:")
	printDuration("Total", calc.TotalRtt)

	printDuration("mean", time.Duration(calc.Mean))
	printDuration("std", time.Duration(calc.StdDev()))
	printDuration("min", calc.Min)
	printDuration("max", calc.Max)

	printPercentiles(calc)

	fmt.Println(strings.Repeat("=", 30))
}

func printDuration(lbl string, duration time.Duration) {
	fmt.Printf("\t%s:\t%s\n", lbl, duration)
}

func printPercentiles(calc Metrics) {
	_25m := calc.GetPercentile(0.25)
	_75m := calc.GetPercentile(0.75)
	_95m := calc.GetPercentile(0.95)
	_99m := calc.GetPercentile(0.99)
	_999m := calc.GetPercentile(0.999)

	fmt.Println("Percentiles:")
	printDuration("0.25", _25m)
	printDuration("0.75", _75m)
	printDuration("0.95", _95m)
	printDuration("0.99", _99m)
	printDuration("0.999", _999m)
}
