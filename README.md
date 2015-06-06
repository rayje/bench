# bench
A Go client for latency testing

## Overview
This client will make HTTP requests to a remote host, capturing the round trip time and calculating measurements to determine variablility.

The measurements captured are:

* Mean
* Standard Deviation
* Min
* Max
* Inter-Quartile range
* 25 Percentile
* 75 Percentile
* 95 Percentile
* 99 Percentile
* 999 Percentile

## Getting Started

    $ bench -host=192.168.1.2 -port=8080 -endpoint=/
    Total Requests: 1
	Throttle: 2s
	Remote: 192.168.1.2:8080 : 192.168.1.2:8080
	==============================
	Requests: 1
	Latencies:
		Total:	65.505638ms
		mean:	65.505638ms
		std:	0
		min:	65.505638ms
		max:	65.505638ms
	Percentiles:
		0.25:	65.505638ms
		0.75:	65.505638ms
		0.95:	65.505638ms
		0.99:	65.505638ms
		0.999:	65.505638ms
	==============================

## Installation

    $ go get github.com/rayje/bench



