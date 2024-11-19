package main

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup

	wg sync.WaitGroup
}

func NewConsumer() (*Consumer, error) {
	config := NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(address, ","), group, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumerGroup: consumerGroup}, nil
}

func (c *Consumer) Run(ctx context.Context) {
	c.wg.Add(2)

	go func() {
		defer c.wg.Done()
		for err := range c.consumerGroup.Errors() {
			log.WithError(err).Error("error on consumer group")
		}
	}()

	go func() {
		defer c.wg.Done()

		for {
			err := c.consumerGroup.Consume(ctx, strings.Split(topics, ","), &ConsumerHandler{})
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}

				log.WithError(err).Error("consume from group")
			}
		}
	}()
}

func (c *Consumer) Close() error {
	err := c.consumerGroup.Close()
	c.wg.Wait()
	return err
}

type ConsumerHandler struct{}

func (c ConsumerHandler) Setup(s sarama.ConsumerGroupSession) error {
	log.Info("setup")
	for topic, partitions := range s.Claims() {
		log.
			WithFields(log.Fields{
				"topic":      topic,
				"partitions": partitions,
			}).
			Info("received claim")
	}
	return nil
}

func (c ConsumerHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.
				WithFields(log.Fields{
					"message": message,
					"offset":  message.Offset,
				}).
				Info("received message")
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			if session.Context().Err() == context.Canceled {
				return nil
			}
			return session.Context().Err()
		}
	}
}

func NewConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = 10 * 1000
	config.Consumer.Group.Heartbeat.Interval = 3 * 1000
	config.Consumer.Group.Rebalance.Timeout = 60 * 1000
	config.Consumer.Group.Rebalance.Retry.Max = 3
	config.Consumer.Group.Rebalance.Retry.Backoff = 3 * 1000

	switch assigner {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	default:
		log.Fatalf("unkown consumer group strategy: %s", assigner)
	}

	if err := config.Validate(); err != nil {
		log.WithError(err).Fatal("validate config")
	}

	return config
}
