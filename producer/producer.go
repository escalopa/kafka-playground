package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	kafka_playground "github.com/escalopa/kafka-playground"
	log "github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer

	once sync.Once
	wg   sync.WaitGroup
}

func NewProducer() (*Producer, error) {
	config := NewConfig()

	producer, err := sarama.NewSyncProducer(strings.Split(address, ","), config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) Run(ctx context.Context) {
	p.once.Do(func() {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(freq):
					m := &sarama.ProducerMessage{
						Topic: topic,
						//Key:   kafka_playground.Key(fmt.Sprintf("group%d", i%2)),
						Value: kafka_playground.NewUser(),
						//Partition: i % 2,
						Timestamp: time.Now(),
					}
					p.Send(m)
				}
			}
		}()
	})
}

func (p *Producer) Send(m *sarama.ProducerMessage) {
	partition, _, err := p.producer.SendMessage(m)
	if err == nil {
		log.
			WithFields(log.Fields{
				"message":   m,
				"partition": partition,
				"offset":    m.Offset,
			}).
			Info("sent message")
		return
	}

	if err == sarama.ErrClosedClient {
		return
	}

	log.WithError(err).Error("send message")
}

func (p *Producer) Close() error {
	p.wg.Wait()
	return p.producer.Close()
}

func NewConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Retry.Backoff = 1 * time.Second  // Wait 1 second between retries
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	if err := config.Validate(); err != nil {
		log.WithError(err).Fatal("validate config")
	}

	return config
}
