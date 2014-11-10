package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	gh "github.com/crosbymichael/octokat"
)

type sourceConstructor func(*cli.Context, *gh.Repo) (Repository, error)

var repository Repository

func createDataSource(cli *cli.Context, target *gh.Repo) (Repository, error) {
	sourceTypes := map[string]sourceConstructor{
		"github":    NewGitHubRepository,
		"rethinkdb": NewRethinkRepository,
	}

	source := cli.String("source")
	if ctor, ok := sourceTypes[source]; ok {
		return ctor(cli, target)
	}
	return nil, fmt.Errorf("invalid datasource %s", source)
}

func before(cli *cli.Context) error {
	if len(cli.Args()) > 0 {
		log.Fatalf("too many arguments")
	}

	repoId, err := parseRepository(cli.String("repository"))
	if err != nil {
		return err
	}
	repository, err = createDataSource(cli, repoId)
	return err
}

func mainCommand(cli *cli.Context) {
	metrics := &Metrics{}
	metrics.Compute()
	metrics.Output()
}

func main() {
	app := cli.NewApp()
	app.Action = mainCommand
	app.Before = before
	app.Name = "octostats"
	app.Usage = "Extract metrics from a github repository"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "source", Value: "github", Usage: "data source"},
		cli.StringFlag{Name: "repository", Value: "docker/docker", Usage: "target repository (e.g: icecrime/docker"},
		cli.StringFlag{Name: "token", Value: "", Usage: "authentication token"},
		cli.StringFlag{Name: "token-file", Value: "", Usage: "authentication token file"},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf(err.Error())
	}
}
