package main

import (
	"github.com/icecrime/octostats/metrics"
)

type Store interface {
	Send(*metrics.Metrics) error
}

type debugStore struct {
}

func (*debugStore) Send(m *metrics.Metrics) error {
	logger.WithField("origin", m.Origin.Nwo()).Info("Sending metrics")
	for k, v := range m.Items {
		logger.Info("  %s = %d", k, v)
	}
	return nil
}
