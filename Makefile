produce:
	go run ./producer/main.go --address "localhost:9092" --topic $(TOPIC)

consume:
	go run ./consumer/main.go --address "localhost:9092" --topic $(TOPIC) --part $(PART)

create-topic:
	docker exec kafka0 /usr/bin/kafka-topics --create \
								--bootstrap-server localhost:9092 \
                                --partitions $(PART) \
                                --topic $(TOPIC)