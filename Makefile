## Variables
BROKERS_LIST="localhost:9001,localhost:9002,localhost:9003,localhost:9004,localhost:9005"
TOPIC="events"
PART=3
REP=3

ACKS=-1

ASSIGNER=sticky # sticky, range, roundrobin
GROUP=random
TOPICS="events"

## Docker commands (Use this if you don't have the kafka-cli installed on your machine)

run:
	docker run -it --rm --name kafka_local --network kafka-playground_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash

## Go app commands (`2>&1` is used to redirect the stderr to stdout)

produce:
	@go run ./producer/main.go --address $(BROKERS_LIST) --topic $(TOPIC)

consume:
	@go run ./group/consumer/main.go --address $(BROKERS_LIST) --group $(GROUP) --assigner $(ASSIGNER) --topics $(TOPIC)

## Topics Commands

list-topic:
	kafka-topics.sh --list --bootstrap-server $(BROKERS_LIST) --exclude-internal

describe-topic:
	kafka-topics.sh --describe --bootstrap-server $(BROKERS_LIST) --topic $(TOPIC) 

describe-topic-all:
	kafka-topics.sh --describe --bootstrap-server $(BROKERS_LIST) --exclude-internal

create-topic:
	kafka-topics.sh --create \
				--if-not-exists \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--replication-factor $(REP) \
				--topic $(TOPIC) \
				--config min.insync.replicas=$(MIN_SYNC) \
				--config cleanup.policy=compact \
				--config compression.type=gzip \
				--config delete.retention.ms=86400000  \
				--config max.message.bytes=104857600 \
				--config retention.bytes=1073741824 \
				--config retention.ms=86400000

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

## Configs

## Partition Commands

partition-verification:
	kafka-replica-verification.sh --broker-list $(BROKERS_LIST) --topics-include $(TOPICS_INCLUDE)

## CLI Commands

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
				--topic $(TOPIC) \
				--timeout 100	\
				--property key.separator=":" \
				--sync
	