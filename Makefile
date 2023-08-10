BROKERS_LIST="localhost:9094,localhost:9095,localhost:9096"

produce:
	go run ./producer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC)

consume:
	go run ./consumer/main.go --address  $(BROKERS_LIST) --topic $(TOPIC) --part $(PART)

consume-group:
	go run ./group/consumer/main.go --address  $(BROKERS_LIST) --group $(GROUP) --assigner $(ASSIGNER) --topics $(TOPICS)

list:
	kafka-topics.sh --list --bootstrap-server $(BROKERS_LIST) --topic $(TOPIC)

create-topic:
	kafka-topics.sh --create \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--replication-factor $(REP) \
				--topic $(TOPIC)

alter-topic:
	kafka-topic.sh --alter \
				--bootstrap-server $(BROKERS_LIST) \
				--partitions $(PART) \
				--topic $(TOPIC)

delete-topic:
	kafka-topic.sh --delete \
				--bootstrap-server $(BROKERS_LIST) \
				--topic $(TOPIC)