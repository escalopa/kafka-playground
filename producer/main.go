package main

import (
	"context"
	"flag"
	"time"

	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
)

var (
	topic   string
	address string
	freq    time.Duration
)

func init() {
	flag.StringVar(&address, "address", "", "kafka brokers address")
	flag.StringVar(&topic, "topic", "", "kafka produce topic name")
	flag.DurationVar(&freq, "freq", 1*time.Second, "frequency to produce message")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	appCtx, cancel := kafka_playground.NewContext(context.Background())
	defer cancel()

	producer, err := NewProducer()
	if err != nil {
		log.WithError(err).Fatal("create producer")
	}

	producer.Run(appCtx)

	log.Info("producer is running")

	<-appCtx.Done() // Wait for termination signal

	log.Info("shutdown has started")
	cancel()

	err = producer.Close()
	if err != nil {
		log.WithError(err).Error("close producer")
	}

	log.Info("shutdown has completed")
}
