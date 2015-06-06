package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Requestor struct {
	Rate       float64
	Duration   time.Duration
	Config     *Config
	Results    Results
	NumResults uint64
}

var client = &http.Client{}

func NewRequestor(config *Config) *Requestor {
	return &Requestor{
		Rate:     config.Rate,
		Duration: config.Duration,
		Config:   config,
	}
}

func (r *Requestor) buildRequest() (*http.Request, error) {
	url := fmt.Sprintf("http://%s:%d%s", r.Config.Host, r.Config.Port, r.Config.Endpoint)

	return http.NewRequest("GET", url, nil)
}

// Generate the HTTP request and start concurrent calls
func (r *Requestor) MakeRequest() (Results, error) {
	total := uint64(r.Rate * float64(r.Duration.Seconds()))
	fmt.Println("Total Requests:", total)

	res := make(chan Result, total)
	r.Results = make(Results, total)

	req, err := r.buildRequest()
	if err != nil {
		return nil, err
	}

	go runRequests(r.Rate, req, res, total)

	for i := 0; i < cap(res); i++ {
		r.Results[i] = <-res
		r.NumResults += 1
	}

	close(res)
	return r.Results, nil
}

// Start concurrent HTTP requests
func runRequests(rate float64, req *http.Request, res chan Result, total uint64) {
	throttle := time.Tick(time.Duration(1e9 / rate))
	fmt.Println("Throttle:", time.Duration(1e9/rate))

	for i := 0; uint64(i) < total; i++ {
		go runRequest(req, res)
		if total > 1 {
			<-throttle
		}
	}
}

// Make HTTP Request
func runRequest(req *http.Request, res chan Result) {
	start := time.Now()
	r, err := client.Do(req)
	stop := time.Since(start)

	if r != nil {
		defer r.Body.Close()
	}

	result := Result{
		Timestamp: start,
		Duration:  stop,
	}

	if err != nil {
		result.Error = err.Error()
	}

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		fmt.Println(err)
	} else if r.StatusCode != 200 {
		result.Error = string(body)
	}

	res <- result
}
