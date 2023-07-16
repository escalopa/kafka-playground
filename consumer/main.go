package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
)

var (
	topic     string
	address   string
	partition int
)

func init() {
	flag.StringVar(&address, "address", "", "kafka brokers address")
	flag.StringVar(&topic, "topic", "", "kafka produce topic name")
	flag.IntVar(&partition, "part", -1, "partition to consume from")
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()

	if topic == "" {
		log.Errorf("topic name not passed, use --topic")
	}

	if address == "" {
		log.Errorf("brokers addresses not passed, use --address")
	}

	if partition < 0 {
		log.Errorf("consumer partion value must be greater than 0, use --part")
	}
}

func main() {
	defer log.Info("application shutdown successfully")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create kafka producer
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Metadata.Full = true

	// Create main kafka consumer
	consumer, err := sarama.NewConsumer(strings.Split(address, ","), config)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to create consumer")
	}
	defer func() {
		err = consumer.Close()
		if err != nil {
			log.WithFields(log.Fields{"error": err, "message": "failed to close consumer"}).Error()
		} else {
			log.WithFields(log.Fields{"message": "successfully closed kafka consumer"}).Info()
		}
	}()
	log.Info("create kafka consumer")

	// Create partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to create partition consumer")
	}
	defer func() {
		err = partitionConsumer.Close()
		if err != nil {
			log.WithFields(log.Fields{"error": err, "message": "failed to close partition consumer"}).Error()
		} else {
			log.WithFields(log.Fields{"message": "successfully closed kafka partition consumer"}).Info()
		}
	}()
	log.Info("create kafka partition consumer")

	// Process errors
	go func() {
		for err := range partitionConsumer.Errors() {
			log.WithFields(log.Fields{"error": err}).Error("failed to consumer message")
		}
	}()

	// Process message
	go func() {
		for m := range partitionConsumer.Messages() {
			log.WithFields(log.Fields{"message": m}).Info("received message")
			var u kafka_playground.User
			err := json.Unmarshal(m.Value, &u)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("failed to unmarshal user")
			}
		}
	}()

	log.Infof("higher partition offSet: %d", partitionConsumer.HighWaterMarkOffset())

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	s := <-sig

	log.Infof("received signal %s", s)
	log.Info("shutdown has started")
	cancel()
}
