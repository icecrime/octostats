package main

import (
	"github.com/icecrime/octostats/log"
	"github.com/icecrime/octostats/metrics"
)

type Store interface {
	Send(*metrics.Metrics) error
}

type debugStore struct {
}

func (*debugStore) Send(m *metrics.Metrics) error {
	log.Logger.WithField("origin", m.Origin.Nwo()).Info("Sending metrics")
	for k, v := range m.Items {
		log.Logger.Info("  %s = %d", k, v)
	}
	return nil
}
