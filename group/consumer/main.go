package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
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
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()
}

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	// Create consumer config
	config := sarama.NewConfig()
	switch assigner {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Fatalf("unkown consumerGroup strategy: %s", assigner)
	}

	if err := config.Validate(); err != nil {
		log.WithError(err).Fatal("failed to validate config")
	}

	// Create consumerGroup
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(address, ","), group, config)
	if err != nil {
		log.WithError(err).Fatal("failed to create consumerGroup")
	}
	log.Info("created consumerGroup")

	// Process consumerGroup errors
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range consumerGroup.Errors() {
			log.WithError(err).Error("received error from consumerGroup")
		}
	}()

	consumer := &ConsumerHandler{
		ready: make(chan bool),
	}

	// Consume messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			err = consumerGroup.Consume(appCtx, strings.Split(topics, ","), consumer)
			if err != nil && err != sarama.ErrClosedConsumerGroup {
				log.WithError(err).Fatal("failed to consumer from group")
			}
			if appCtx.Err() != nil {
				log.Info("stop processing messages: context done")
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Wait for consumerGroup to be created
	log.Info("started consuming messages")

	sigState := make(chan os.Signal, 1)
	signal.Notify(sigState, syscall.SIGSTOP, syscall.SIGCONT)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM)

	keepAlive := true
	for keepAlive {
		var sig os.Signal
		select {
		case <-appCtx.Done():
			log.Info("main context is done!")
			keepAlive = false
		case sig = <-sigState:
			toggleConsumerState(consumerGroup, sig)
		case <-sigTerm:
			log.Info("received termination signal")
			keepAlive = false
		}
	}

	log.Info("shutdown has started")
	cancel()

	// Close consumerGroup
	err = consumerGroup.Close()
	if err != nil {
		log.WithError(err).Error("failed to close consumerGroup")
	} else {
		log.Info("successfully closed consumerGroup")
	}

	log.Info("Waiting for all goroutines to finish")
	wg.Wait()
	log.Info("application shutdown successfully")
}

// toggleConsumerState toggle consumer state pause/resume
func toggleConsumerState(group sarama.ConsumerGroup, sig os.Signal) {
	switch sig {
	case syscall.SIGSTOP:
		group.PauseAll()
		log.Info("pause consuming messages")
	case syscall.SIGCONT:
		group.ResumeAll()
		log.Info("resume consuming message")
	default:
		log.Errorf("received unknow signal on toggleConsumerStart: %s", sig)
	}
}

// ConsumerHandler is an implementation for sarama.ConsumerGroupHandler
type ConsumerHandler struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c ConsumerHandler) Setup(s sarama.ConsumerGroupSession) error {
	close(c.ready)
	log.Info("setup...")
	for topic, partitions := range s.Claims() {
		log.WithFields(log.Fields{"topic": topic, "partitions": partitions}).Info("received claim")
	}
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c ConsumerHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.WithFields(log.Fields{"offset": message.Offset, "message": message}).Info("successfully received message")
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			if session.Context().Err() == context.Canceled {
				return nil
			}
			return session.Context().Err()
		}
	}
}
