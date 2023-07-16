# kafka-playground 😈

This is a simple kafka playground where you can run commands and experience how kafka
works.

## Run 🚀

Clone the repo.

```shell
git clone https://github.com/escalopa/kafka-playground.git
```
### ZK/Kafka 🐼

To start zk/kafka node run.
```shell
docker compose up -d
```

📝 Notice:  the kafka container sometimes crash on startup so consider auto restart when needed.

You can use the following commands for kafka restart.
```shell
docker stop kafka0
docker start kafka0 
```

### Topic 📝

In the docker compose config, creating a new topic on usage is forbidden so you have to create the topic in advance.

To create a topic you have to specify 2 fields.
* `topic`: the topic name.
* `part`: number of partitions in the topic.

Use the following command to create a topic.
```shell
make create-topic TOPIC=topic PART=3
```

### Consumer ⏪

To run consumer you have to specify 2 values.
* `topic`: the topic name to which the consumer will connect.
* `part`: the partition id to consumer from.

Use the make command to run a consumer.
```shell
make consume TOPIC=topic PART=0
```

### Produce ⏩

To run producer you have to specify 1 value.
* `topic`: the topic name to which the producer will connect.

If you want to use `key` | `partition` values, change the values in source code. 🙂

Use the make command to run a producer.
```shell
make produce TOPIC=topic 
```

