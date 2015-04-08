package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/icecrime/octostats/nsq"
)

type GitHubConfig struct {
	AuthToken     string `json:"token"`
	AuthTokenFile string `json:"tokenfile"`
	Repository    string `json:"repository"`
}

type InfluxConfig struct {
	Endpoint string `json:"endpoint"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Output          string `json:"output"`
	StoreEndpoint   string `json:"store"`
	UpdateFrequency string `json:"update_frequency"`

	GitHubConfig   GitHubConfig `json:"github"`
	InfluxDBConfig InfluxConfig `json:"influxdb"`
	NSQConfig      nsq.Config   `json:"nsq"`
}

func Load(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
