package nsq

import "github.com/bitly/go-nsq"

type Config struct {
	Topic      string `json:"topic"`
	Channel    string `json:"channel"`
	LookupAddr string `json:"lookup_address"`
}

func New(config *Config, handler nsq.Handler) (*Queue, error) {
	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	consumer.AddHandler(handler)
	if err := consumer.ConnectToNSQLookupd(config.LookupAddr); err != nil {
		return nil, err
	}

	return &Queue{Consumer: consumer}, nil
}

type Queue struct {
	Consumer *nsq.Consumer
}
