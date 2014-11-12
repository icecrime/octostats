package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icecrime/octostats/metrics"
	"github.com/icecrime/octostats/nsq"

	"github.com/codegangsta/cli"
)

func updateTicker() *time.Ticker {
	duration, err := time.ParseDuration(config.UpdateFrequency)
	if err != nil {
		logger.Fatal(err)
	}
	return time.NewTicker(duration)
}

func onTimerTick() {
	logger.Debug("Tick: fetching statistics")
	stats := metrics.Retrieve(source)
	if err := store.Send(stats); err != nil {
		logger.Error(err)
	}
}

func mainCommand(cli *cli.Context) {
	s := make(chan os.Signal, 64)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	queue, err := nsq.New(&config.NSQConfig, NewNSQHandler())
	if err != nil {
		logger.Fatal(err)
	}

	ticker := updateTicker()
	for {
		select {
		case <-ticker.C:
			onTimerTick()
		case <-queue.Consumer.StopChan:
			logger.Debug("Queue stop channel signaled")
			return
		case sig := <-s:
			logger.WithField("signal", sig).Debug("received signal")
			queue.Consumer.Stop()
		}
	}
}
