package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/icecrime/octostats/influx"
	"github.com/icecrime/octostats/nsq"

	log "github.com/Sirupsen/logrus"
)

type GitHubConfig struct {
	AuthToken     string `json:"token"`
	AuthTokenFile string `json:"tokenfile"`
	Repository    string `json:"repository"`
}

type Config struct {
	Output          string `json:"output"`
	StoreEndpoint   string `json:"store"`
	UpdateFrequency string `json:"update_frequency"`

	GitHubConfig   GitHubConfig  `json:"github"`
	InfluxDBConfig influx.Config `json:"influxdb"`
	NSQConfig      nsq.Config    `json:"nsq"`
}

func loadConfig(filename string) *Config {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatal(err)
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("failed to unmarshal config")
	}
	return &config
}
