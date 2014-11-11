package graphite

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/icecrime/octostats/stats"
)

func New(target string) *store {
	return &store{target: target}
}

type store struct {
	target string
}

func (*store) format(repository stats.Repository, metrics stats.Metrics) []byte {
	timestamp := time.Now().Unix()
	metricsPrefix := fmt.Sprintf("github.%s.%s", repository.Id().UserName, repository.Id().Name)

	var buffer bytes.Buffer
	for key, value := range metrics {
		buffer.WriteString(fmt.Sprintf("%s.%s %d %d\n", metricsPrefix, key, value, timestamp))
	}
	return buffer.Bytes()
}

func (s *store) Send(repository stats.Repository, metrics stats.Metrics) error {
	conn, err := net.Dial("tcp", s.target)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := s.format(repository, metrics)
	_, err = conn.Write(payload)
	return err
}
