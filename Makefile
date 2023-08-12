BROKERS_LIST="localhost:9094,localhost:9095,localhost:9096"
TOPIC="test"
PART=3
REP=3

ASSIGNER=roundrobin
GROUP=group1

run:
	docker run -it --rm --name kafka_env --network kafka-playground_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash

produce:
	go run ./producer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC)

consume:
	go run ./consumer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC) --part $(PART)

consume-group:
	go run ./group/consumer/main.go --address  $(BROKERS_LIST) --group $(GROUP) --assigner $(ASSIGNER) --topics $(TOPICS)

list:
	kafka-topics.sh --list --bootstrap-server $(BROKERS_LIST) --topic $(TOPIC)

describe-group:
	kafka-consumer-groups.sh --bootstrap-server $(BROKERS_LIST) --describe --group $(GROUP)	

create-topic:
	kafka-topics.sh --create \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--replication-factor $(REP) \
				--topic $(TOPIC)

alter-topic:
	kafka-topics.sh --alter \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--topic $(TOPIC)

delete-topic:
	kafka-topics.sh --delete \
				--bootstrap-server $(BROKERS_LIST) \
				--topic $(TOPIC)

consume-topic:
	kafka-console-consumer.sh \
				--bootstrap-server $(BROKERS_LIST) \
				--property print.key=true \
				--property key.separator=":" 	\
				--group $(GROUP) \
				--topic $(TOPIC)

produce-topic:
	kafka-console-producer.sh \
				--bootstrap-server $(BROKERS_LIST) \
				--property parse.key=true \
				--request-required-acks 0	\
				--timeout 100	\
				--property key.separator=":" \
				--topic $(TOPIC) 
	