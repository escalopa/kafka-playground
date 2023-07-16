produce:
	go run ./producer/main.go --address "localhost:9092" --topic $(TOPIC_NAME)

consume:
	go run ./consumer/main.go --address "localhost:9092" --topic $(TOPIC_NAME) --part $(KAFKA_PART)

create-topic:
	docker exec kafka0 /usr/bin/kafka-topics --create \
								--bootstrap-server localhost:9092 \
                                --partitions 3 \
                                --topic $(TOPIC_NAME)