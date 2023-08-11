# kafka-playground 💪

This is a simple kafka playground where you can run commands and experience how kafka
works.

## Run 🚀

Clone the repo.

```shell
git clone https://github.com/escalopa/kafka-playground.git
```
### ZK/Kafka 🐼

To start zk/kafka nodes run.
```shell
docker compose up -d
```

📝 Notice:  the kafka container sometimes hangs on startup so consider auto restart when needed.

You can use the following commands for restart the container.
```shell
docker stop kafka0
docker start kafka0 
```

### Topic 📝

In the docker compose config, creating a new topic on usage is forbidden, so you have to create the topic in advance.

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
* `part`: the partition id to consume messages from.

Use the make command to run a consumer.
```shell
make consume TOPIC=topic PART=0
```

### Produce ⏩

To run producer you have to specify 1 value.
* `topic`: the topic name to which the producer will connect.

If you want to use `key` | `partition` values, change the values in source code 🙂. By default the message are spread using `round-robin`

Use the make command to run a producer.
```shell
make produce TOPIC=topic 
```

### Consumer-Group ⏪⏪⏪

ConsumerGroup as the name suggest consumes from more than one  partions/topic at the same time

In a single group can participate up to `N` consumers where `N` is the total numebr of partition of all topics(You can use more than `N` but is useless since you will have consumers with no partitions to listen to)

To run a consumer-group you have to specify the following values
* `group`: the name of the group
* `assigner`: the strategy for spreading messages between consumers, must be one out of `sticky,roundrobin,range` otherwise the code panics
* `topics`: topics names to consume from

Use the following command to run consumer-group
```shell
make consume-group GROUP=g0 ASSIGNER="sticky" TOPICS="topic1,topic2,topic3"
```
