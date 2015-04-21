package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/icecrime/octostats/config"
	"github.com/icecrime/octostats/github"
	"github.com/icecrime/octostats/influx"
	"github.com/icecrime/octostats/log"
	"github.com/icecrime/octostats/repository"
)

var (
	source       repository.Repository
	store        Store
	globalConfig *config.Config
)

func newStore(c *config.Config) Store {
	switch c.Output {
	case "console":
		return &debugStore{}
	case "influxdb":
		return influx.New(&c.InfluxDBConfig)
	default:
		log.Logger.Fatalf("Invalid output '%s'", c.Output)
		return nil
	}
}

func before(cli *cli.Context) error {
	log.Configure(cli.String("loglevel"))
	if len(cli.Args()) > 0 {
		log.Logger.Fatal("too many arguments")
	}

	var err error
	globalConfig, err = config.Load(cli.String("config"))
	if err != nil {
		log.Logger.Fatal(err)
	}

	store = newStore(globalConfig)
	source, err = github.NewGitHubRepository(&globalConfig.GitHubConfig)

	return err
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
		log.Logger.Fatalf(err.Error())
	}
}
