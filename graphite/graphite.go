package graphite

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/icecrime/octostats/metrics"
)

type Config struct {
	Endpoint string `json:"endpoint"`
}

func New(config *Config) *store {
	return &store{endpoint: config.Endpoint}
}

type store struct {
	endpoint string
}

func (*store) format(metrics *metrics.Metrics) []byte {
	timestamp := time.Now().Unix()
	metricsPrefix := metrics.Origin.Nwo()

	var buffer bytes.Buffer
	for key, value := range metrics.Items {
		buffer.WriteString(fmt.Sprintf("%s.%s %d %d\n", metricsPrefix, key, value, timestamp))
	}
	return buffer.Bytes()
}

func (s *store) Send(metrics *metrics.Metrics) error {
	conn, err := net.Dial("tcp", s.endpoint)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := s.format(metrics)
	_, err = conn.Write(payload)
	return err
}
