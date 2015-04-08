package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/icecrime/octostats/metrics"
)

type partialPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		CreatedAt time.Time  `json:"created_at"`
		ClosedAt  *time.Time `json:"closed_at"`
		Merged    bool       `json:"merged"`
	} `json:"pull_request"`
}

func NewNSQHandler() *NSQHandler {
	return &NSQHandler{store: store}
}

type NSQHandler struct {
	store Store
}

func (n *NSQHandler) HandleMessage(m *nsq.Message) error {
	logger.Debug("Queue event received")

	var p partialPayload
	if err := json.Unmarshal(m.Body, &p); err != nil {
		return nil
	}

	stats := metrics.New(source)
	if p.Action == "closed" && p.PullRequest.ClosedAt != nil {
		mergeString := map[bool]string{true: "merged", false: "not_merged"}
		metricsPath := fmt.Sprintf("pull_requests.close_delay.%s", mergeString[p.PullRequest.Merged])
		hours := int(p.PullRequest.ClosedAt.Sub(p.PullRequest.CreatedAt).Hours())
		metric := metrics.NewMetric(metricsPath, map[string]interface{}{"count": hours})
		stats.Items = append(stats.Items, metric)
	}

	if err := n.store.Send(stats); err != nil {
		logger.Error(err)
	}
	return nil
}
