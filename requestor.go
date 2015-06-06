package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Requestor struct {
	Rate       float64
	Duration   time.Duration
	Config     *Config
	Results    Results
	NumResults uint64
}

var t = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		c, err := net.Dial(network, addr)
		if c == nil {
			fmt.Println("No Connection")
			return c, err
		}

		remote := c.RemoteAddr()
		if remote != nil {
			fmt.Printf("Remote: %s : %s\n", remote.String(), addr)
		}
		return c, err
	},
}

var client = &http.Client{Transport: t}

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
	done := make(chan string)
	r.Results = make(Results, total)

	req, err := r.buildRequest()
	if err != nil {
		return nil, err
	}

	go runRequests(r.Rate, req, res, total, done)

	for i := 0; i < cap(res); i++ {
		r.Results[i] = <-res
		r.NumResults += 1
	}
	close(res)

	return r.Results, nil
}

// Start concurrent HTTP requests
func runRequests(rate float64, req *http.Request, res chan Result, total uint64, done chan string) {
	throttle := time.Tick(time.Duration(1e9 / rate))
	fmt.Println("Throttle:", time.Duration(1e9/rate))

	for i := 0; uint64(i) < total; i++ {
		<-throttle
		go runRequest(req, res)
	}

	done <- "done"
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
		BytesOut:  uint64(req.ContentLength),
	}

	if err != nil {
		result.Error = err.Error()
	}

	result.Code = uint16(r.StatusCode)

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		fmt.Println(err)
	} else {
		if result.Code == 200 {
			result.BytesIn = uint64(len(body))
		} else {
			fmt.Println(result.Code)
			printError(req, r, string(body), "Invalid Status Code")
		}
	}

	res <- result
}

func getDurations(res *http.Response) []int64 {
	durationsString := res.Header.Get("Durations")
	durationsArray := strings.Split(durationsString, ",")

	var durations = make([]int64, len(durationsArray))
	for i := 0; i < len(durationsArray); i++ {
		durations[i], _ = strconv.ParseInt(durationsArray[i], 10, 64)
	}

	return durations
}

func printError(req *http.Request, res *http.Response, body string, msg string) {
	fmt.Println(strings.Repeat("-", 30))
	if msg != "" {
		fmt.Println("Error: " + msg)
	}
	fmt.Println("Request: " + req.URL.String())
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("Status: " + res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Print(string(body))
	fmt.Println(strings.Repeat("-", 30))
}
