package main

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/brianvoe/gofakeit/v6"
	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
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

type errorSend struct {
	err       error
	partition int32
}

func main() {
	defer log.Info("application shutdown successfully")

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	log.Infof("using brokers with address: %s ", strings.Split(address, ","))
	producer, err := sarama.NewSyncProducer(strings.Split(address, ","), config)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "message": "failed to create producer"}).Fatal()
	}
	defer func() {
		err = producer.Close()
		if err != nil {
			log.WithFields(log.Fields{"error": err, "message": "failed to close producer"}).Error()
		} else {
			log.WithFields(log.Fields{"message": "successfully closed kafka producer"}).Info()
		}
	}()
	log.Info("created kafka producer")

	messages := make(chan *sarama.ProducerMessage)
	errors := make(chan errorSend, 1000)

	// Process messages
	go func() {
		for m := range messages {
			partition, _, err := producer.SendMessage(m)
			if err != nil {
				go func() {
					select {
					case <-appCtx.Done():
					case errors <- errorSend{err: err, partition: partition}:
					}
				}()
			} else {
				log.WithFields(log.Fields{"message": m}).Info("message produced")
			}
		}
	}()

	// Process errors
	go func() {
		for e := range errors {
			log.WithFields(log.Fields{"error": e}).Info("failed to send message")
		}
	}()

	// Produce random message
	go func() {
		var i int32
		for {
			duration := time.Duration(rand.Intn(2))
			select {
			case <-appCtx.Done():
				return
			case <-time.After(duration * time.Second):
				select {
				case <-appCtx.Done():
					return
				case messages <- &sarama.ProducerMessage{
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
				}:
				}
			}
			i++
		}
	}()

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Info("shutdown has started")
	cancel()

	close(messages)
	log.Info("closed message channel")

	close(errors)
	log.Info("close errors channel")
}
