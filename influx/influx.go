package influx

import (
	"fmt"

	"github.com/icecrime/octostats/metrics"

	influxClient "github.com/influxdb/influxdb/client"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(config *Config) *store {
	return &store{
		config: config,
	}
}

type store struct {
	config *Config
}

func (*store) format(metrics *metrics.Metrics) []*influxClient.Series {
	series := []*influxClient.Series{}
	metricsPrefix := metrics.Origin.Nwo()

	for _, m := range metrics.Items {
		var columns []string
		var points [][]interface{}

		for k, v := range m.Data {
			columns = append(columns, k)
			points = append(points, []interface{}{v})
		}

		series = append(series, &influxClient.Series{
			Name:    fmt.Sprintf("%s.%s", metricsPrefix, m.Path),
			Columns: columns,
			Points:  points,
		})
	}

	return series
}

func (s *store) Send(metrics *metrics.Metrics) error {
	client, err := influxClient.NewClient(&influxClient.ClientConfig{
		Host:     s.config.Endpoint,
		Database: s.config.Database,
		Username: s.config.Username,
		Password: s.config.Password,
	})
	if err != nil {
		return err
	}
	return client.WriteSeries(s.format(metrics))
}
