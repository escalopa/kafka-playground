package main

import (
	"context"
	"flag"

	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
)

var (
	topics   string
	address  string
	assigner string
	group    string
)

func init() {
	flag.StringVar(&address, "address", "", "kafka brokers address")
	flag.StringVar(&topics, "topics", "", "kafka produce topic name")
	flag.StringVar(&assigner, "assigner", "", "assigner protocol")
	flag.StringVar(&group, "group", "", "consumer group name")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	appCtx, cancel := kafka_playground.NewContext(context.Background())
	defer cancel()

	consumer, err := NewConsumer()
	if err != nil {
		log.WithError(err).Fatal("create consumer")
	}

	consumer.Run(appCtx)

	log.Info("consumer is running")

	<-appCtx.Done() // wait for shutdown signal

	log.Info("shutdown has started")
	cancel()

	err = consumer.Close()
	if err != nil {
		log.WithError(err).Error("close consumer group")
	}

	log.Info("shutdown has completed")
}
