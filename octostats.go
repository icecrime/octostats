package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icecrime/octostats/log"
	"github.com/icecrime/octostats/metrics"
	"github.com/icecrime/octostats/nsq"

	"github.com/codegangsta/cli"
)

func updateTicker() *time.Ticker {
	duration, err := time.ParseDuration(globalConfig.UpdateFrequency)
	if err != nil {
		log.Logger.Fatal(err)
	}
	return time.NewTicker(duration)
}

func onTimerTick() {
	log.Logger.Debug("Tick: fetching statistics")
	stats := metrics.Retrieve(source)
	if err := store.Send(stats); err != nil {
		log.Logger.Error(err)
	}
}

func mainCommand(cli *cli.Context) {
	s := make(chan os.Signal, 64)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	var queue *nsq.Queue
	var queueStopChan chan int

	if globalConfig.NSQConfig != nil {
		var err error
		queue, err = nsq.New(globalConfig.NSQConfig, NewNSQHandler())
		if err != nil {
			log.Logger.Fatal(err)
		}
		queueStopChan = queue.Consumer.StopChan
	} else {
		queueStopChan = make(chan int, 1)
	}

	ticker := updateTicker()
	for {
		select {
		case <-ticker.C:
			onTimerTick()
		case <-queueStopChan:
			log.Logger.Debug("Queue stop channel signaled")
			return
		case sig := <-s:
			log.Logger.WithField("signal", sig).Debug("received signal")
			if queue != nil {
				queue.Consumer.Stop()
			}
		}
	}
}
