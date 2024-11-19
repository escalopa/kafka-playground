# kafka-playground 💪

This is a simple kafka playground where you can run commands and experience how kafka works.

## Run 🚀

Clone the repo.

```shell
git clone https://github.com/escalopa/kafka-playground.git
```

### KRaft 🐼

To run the kafka cluster you have to use the following command.

```shell
docker compose up -d
```

### Topic 📝

In the docker compose config, creating a new topic on usage is forbidden, so you have to create the topic in advance.

To create a topic you have to specify 2 fields.

* `topic`: the topic name.
* `part`: number of partitions in the topic.

Use the make command to create a topic.

```shell
make create-topic TOPIC=topic PART=3
```

### Produce ⏩

To run producer you have to specify 1 value.

* `topic`: the topic name to which the producer will connect.
* `freq`: the frequency of the message production in seconds. (default is 1s)

If you want to use `key` | `partition` values, change the values in source code 🙂. By default the message are spread using `round-robin`

Use the make command to run a producer.

```shell
make produce TOPIC=topic 
```

### Consumer-Group ⏪

ConsumerGroup as the name suggest consumes from more than one  partions/topic at the same time

In a single group can participate up to `N` consumers where `N` is the total numebr of partition of all topics(You can use more than `N` but is useless since you will have consumers with no partitions to listen to)

To run a consumer-group you have to specify the following values

* `group`: the name of the group
* `assigner`: the strategy for spreading messages between consumers, must be one out of `sticky,roundrobin,range` otherwise the code panics
* `topics`: topics names to consume from

Use the following command to run consumer-group

```shell
make consume GROUP=g0 ASSIGNER="sticky" TOPICS="topic1,topic2,topic3"
```
