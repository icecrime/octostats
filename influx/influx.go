package influx

import (
	"fmt"

	"github.com/icecrime/octostats/stats"

	influxClient "github.com/influxdb/influxdb/client"
)

func New(target string) *store {
	return &store{target: target}
}

type store struct {
	target string
}

func (*store) format(repository stats.Repository, metrics stats.Metrics) []*influxClient.Series {
	series := []*influxClient.Series{}
	metricsPrefix := fmt.Sprintf("%s.%s", repository.Id().UserName, repository.Id().Name)

	for k, v := range metrics {
		series = append(series, &influxClient.Series{
			Name:    fmt.Sprintf("%s.%s", metricsPrefix, k),
			Columns: []string{"count"},
			Points:  [][]interface{}{{v}},
		})
	}

	return series
}

func (s *store) Send(repository stats.Repository, metrics stats.Metrics) error {
	client, err := influxClient.NewClient(&influxClient.ClientConfig{
		Host:     s.target,
		Database: "github",
	})
	if err != nil {
		return err
	}
	return client.WriteSeries(s.format(repository, metrics))
}
