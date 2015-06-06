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
	if results, err := requestor.MakeRequest(); err != nil {
		fmt.Println(err)
	} else {
		report(results)
	}
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
	metric := NewMetrics(results)
	metric.Sort()

	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Requests:", metric.Total)
	if metric.Errors > 0 {
		fmt.Println("Errors:", metric.Errors)
	}

	fmt.Println("Latencies:")
	printDuration("Total", metric.TotalRtt)

	if len(results) > 1 {
		printDuration("mean", time.Duration(metric.Mean))
		printDuration("std", time.Duration(metric.StdDev()))
		printDuration("min", metric.Min)
		printDuration("max", metric.Max)
		printPercentiles(metric)
	}

	fmt.Println(strings.Repeat("-", 30))
}

func printDuration(lbl string, duration time.Duration) {
	fmt.Printf("\t%s:\t%s\n", lbl, duration)
}

func printPercentiles(metric *Metrics) {
	_95m := metric.GetPercentile(0.95)
	_99m := metric.GetPercentile(0.99)
	_999m := metric.GetPercentile(0.999)

	fmt.Println("Percentiles:")
	printDuration("0.95", _95m)
	printDuration("0.99", _99m)
	printDuration("0.999", _999m)
}
