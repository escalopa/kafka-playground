package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
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

	if topics == "" {
		log.Fatal("topics name not passed, use --topics")
	}
	if address == "" {
		log.Fatal("brokers addresses not passed, use --address")
	}
	if assigner == "" {
		log.Fatal("group assigner must be set, use --assigner")
	}
	if group == "" {
		log.Fatal("group name must be set, use --group")
	}
}

func main() {
	defer log.Info("application shutdown successfully")

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Create client
	client, err := sarama.NewClient(strings.Split(address, ","), config)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to create client")
	}
	defer func() {
		err = client.Close()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("failed to close kafka client")
		} else {
			log.Info("successfully closed kafka client")
		}
	}()
	log.Info("created kafka client")

	// Create main consumerGroup
	consumerGroup, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to create consumerGroup")
	}
	defer func() {
		err = consumerGroup.Close()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("failed to close consumerGroup")
		} else {
			log.Info("successfully closed consumerGroup")
		}
	}()
	log.Info("created consumerGroup")

	// Process consumerGroup errors
	go func() {
		for err := range consumerGroup.Errors() {
			log.WithFields(log.Fields{"error": err}).Warn("received error from consumerGroup")
		}
	}()

	consumer := &ConsumerHandler{
		ready: make(chan bool),
	}

	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			err = consumerGroup.Consume(appCtx, strings.Split(topics, ","), consumer)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("failed to consumer from group")
			}
			if appCtx.Err() != nil {
				log.WithFields(log.Fields{"error": err}).Info("stop processing messages")
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
		case sig = <-sigTerm:
			log.Info("received termination signal")
			keepAlive = false
		}
	}

	cancel()
	log.Info("shutdown has started")
}

// toggleConsumerState toggle consumer state pause/resume
func toggleConsumerState(client sarama.ConsumerGroup, sig os.Signal) {
	switch sig {
	case syscall.SIGSTOP:
		client.PauseAll()
		log.Info("pause consuming messages")
	case syscall.SIGCONT:
		client.ResumeAll()
		log.Info("resume consuming message")
	default:
		log.Errorf("received unknow signal on toggleConsumerStaet: %s", sig)
	}
}

// ConsumerHandler is an implementation for sarama.ConsumerGroupHandler
type ConsumerHandler struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.WithFields(log.Fields{"offset": message.Offset, "message": message}).Info("successfully received message")
			session.Commit()
		case <-session.Context().Done():
			if session.Context().Err() == context.Canceled {
				return nil
			}
			return session.Context().Err()
		}
	}
}
