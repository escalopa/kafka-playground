## Variables

BROKERS_LIST="localhost:9094,localhost:9095,localhost:9096,localhost:9097,localhost:9098"
TOPIC="quickstart-events-0"
PART=3
REP=3

ACKS=-1

ASSIGNER=sticky # sticky, range, roundrobin
GROUP=group1
TOPICS="quickstart-events-1,quickstart-events-2"

## Docker commands (Use this if you don't have the kafka-cli installed on your machine)

run:
	docker run -it --rm --name kafka_env --network kafka-playground_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash

## Golang application commands

produce:
	go run ./producer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC)

consume:
	go run ./consumer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC) --part $(PART)

consume-group:
	go run ./group/consumer/main.go --address  $(BROKERS_LIST) --group $(GROUP) --assigner $(ASSIGNER) --topics $(TOPICS)

## Topics Commands

list-topic:
	kafka-topics.sh --list --bootstrap-server $(BROKERS_LIST) --exclude-internal

describe-topic:
	kafka-topics.sh --describe --bootstrap-server $(BROKERS_LIST) --topic $(TOPIC)

describe-topic-all:
	./scripts/describe-topic-all.sh $(BROKERS_LIST) 

create-topic:
	kafka-topics.sh --create \
				--if-not-exists \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--replication-factor $(REP) \
				--config min.insync.replicas=$(MIN_SYNC) \
				--topic $(TOPIC)

alter-topic:
	kafka-topics.sh --alter \
				--if-exists \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--topic $(TOPIC)

delete-topic:
	kafka-topics.sh --delete \
				--if-exists \
				--bootstrap-server $(BROKERS_LIST) \
				--topic $(TOPIC)

## Group Commands

list-group:
	kafka-consumer-groups.sh --bootstrap-server $(BROKERS_LIST) --list

describe-group:
	kafka-consumer-groups.sh --bootstrap-server $(BROKERS_LIST) --describe --group $(GROUP)	

describe-group-all:
	./scripts/describe-group-all.sh $(BROKERS_LIST) 

## Cli Commands

consume-cli	:
	kafka-console-consumer.sh \
				--bootstrap-server $(BROKERS_LIST) \
				--property print.key=true \
				--property key.separator=":" 	\
				--group $(GROUP) \
				--topic $(TOPIC)

produce-cli:
	kafka-console-producer.sh \
				--bootstrap-server $(BROKERS_LIST) \
				--property parse.key=true \
				--request-required-acks $(ACKS)	\
				--timeout 100	\
				--property key.separator=":" \
				--topic $(TOPIC) 
	