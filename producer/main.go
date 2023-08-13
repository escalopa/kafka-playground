package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/brianvoe/gofakeit/v6"
	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
)

const (
	freq = 1 * time.Nanosecond
)

var (
	topic   string
	address string
)

func init() {
	flag.StringVar(&address, "address", "", "kafka brokers address")
	flag.StringVar(&topic, "topic", "", "kafka produce topic name")
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()
}

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	// Create kafka producer config
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Retry.Backoff = 1 * time.Second  // Wait 1 second between retries
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	if err := config.Validate(); err != nil {
		log.WithError(err).Fatal("failed to validate config")
	}

	// Create kafka producer
	log.Infof("using brokers with address: %s ", strings.Split(address, ","))
	producer, err := sarama.NewSyncProducer(strings.Split(address, ","), config)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "message": "failed to create producer"}).Fatal()
	}
	log.Info("created kafka producer")

	// Send message function that handler errors
	send := func(m *sarama.ProducerMessage) {
		partition, _, err := producer.SendMessage(m)
		if err != nil {
			if err == sarama.ErrClosedClient {
				return
			}
			log.WithFields(log.Fields{"error": err.Error(), "partition": partition}).Info("failed to send message")
		} else {
			log.WithFields(log.Fields{"message": m}).Info("message produced")
		}
	}

	// Produce random message
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-appCtx.Done():
				return
			case <-time.After(freq):
				m := &sarama.ProducerMessage{
					Topic: topic,
					//Key:   kafka_playground.Key(fmt.Sprintf("group%d", i%2)),
					Value: kafka_playground.User{
						ID:        gofakeit.UUID(),
						Name:      gofakeit.Name(),
						Email:     gofakeit.Email(),
						CreatedAt: time.Now(),
					},
					//Partition: i % 2,
					Timestamp: time.Now(),
				}
				send(m)
			}
		}
	}()

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Info("shutdown has started")
	cancel()

	// Close producer
	err = producer.Close()
	if err != nil {
		log.WithError(err).Error("failed to close producer")
	} else {
		log.Info("successfully closed kafka producer")
	}

	log.Info("Waiting for all goroutines to finish")
	wg.Wait()

	log.Info("application shutdown successfully")
}
