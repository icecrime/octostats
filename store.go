package main

import "github.com/icecrime/octostats/stats"

type Store interface {
	Send(stats.Repository, stats.Metrics) error
}
