package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/icecrime/octostats/graphite"
	"github.com/icecrime/octostats/influx"
	"github.com/icecrime/octostats/stats"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	gh "github.com/crosbymichael/octokat"
)

var repository stats.Repository

func newStore(cli *cli.Context) Store {
	target := cli.String("target")
	switch cli.String("output") {
	case "graphite":
		return graphite.New(target)
	case "influx":
		return influx.New(target, cli.String("influx-database"), cli.String("influx-username"), cli.String("influx-password"))
	default:
		log.Fatal("Invalid output '%s'", cli.String("output"))
		return nil
	}
}

func parseRepository(repo string) (*gh.Repo, error) {
	if splitRepos := strings.Split(repo, "/"); len(splitRepos) == 2 {
		return &gh.Repo{Name: splitRepos[1], UserName: splitRepos[0]}, nil
	}
	return nil, fmt.Errorf("bad repo format %s (expected username/repo)", repo)
}

func before(cli *cli.Context) error {
	if len(cli.Args()) > 0 {
		log.Fatalf("too many arguments")
	}

	repoId, err := parseRepository(cli.String("repository"))
	if err != nil {
		return err
	}

	repository, err = NewGitHubRepository(cli, repoId)
	return err
}

func mainCommand(cli *cli.Context) {
	metrics := stats.Metrics{}
	metrics.Compute(repository)

	store := newStore(cli)
	if err := store.Send(repository, metrics); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Action = mainCommand
	app.Before = before
	app.Name = "octostats"
	app.Usage = "Extract metrics from a github repository"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "output", Value: "influx", Usage: "output (influx|graphite)"},
		cli.StringFlag{Name: "repository", Value: "docker/docker", Usage: "target repository (e.g: icecrime/docker"},
		cli.StringFlag{Name: "target", Value: "", Usage: "endpoint to send the output to"},
		cli.StringFlag{Name: "token", Value: "", Usage: "authentication token"},
		cli.StringFlag{Name: "token-file", Value: "", Usage: "authentication token file"},
		cli.StringFlag{Name: "influx-database", Value: "", Usage: "InfluxDB database to write to"},
		cli.StringFlag{Name: "influx-password", Value: "", Usage: "password for InfluxDB"},
		cli.StringFlag{Name: "influx-username", Value: "", Usage: "username for InfluxDB"},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf(err.Error())
	}
}
