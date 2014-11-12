package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/icecrime/octostats/graphite"
	"github.com/icecrime/octostats/influx"
	"github.com/icecrime/octostats/repository"
)

var (
	config *Config
	source repository.Repository
	store  Store

	logger = logrus.New()
)

func configureLogger(loglevel string) {
	if level, err := logrus.ParseLevel(loglevel); err != nil {
		logger.Fatal(err)
	} else {
		logger.Level = level
	}
}

func newStore(config *Config) Store {
	switch config.Output {
	case "console":
		return &debugStore{}
	case "graphite":
		return graphite.New(&config.GraphiteConfig)
	case "influxdb":
		return influx.New(&config.InfluxDBConfig)
	default:
		logger.Fatal("Invalid output '%s'", config.Output)
		return nil
	}
}

func before(cli *cli.Context) error {
	configureLogger(cli.String("loglevel"))
	if len(cli.Args()) > 0 {
		logger.Fatalf("too many arguments")
	}

	config = loadConfig(cli.String("config"))

	store = newStore(config)
	source = NewGitHubRepository(&config.GitHubConfig)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Action = mainCommand
	app.Before = before
	app.Name = "octostats"
	app.Usage = "Extract metrics from a github repository"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config", Value: "octostats.json", Usage: "configuration file"},
		cli.StringFlag{Name: "loglevel", Value: "info", Usage: "logging level"},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatalf(err.Error())
	}
}
