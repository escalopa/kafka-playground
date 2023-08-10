# [](https://)Kafka workshop-3

## Основные термины (Шпаргалка)

* Broker - Инстанс кафки (Завод)
* Cluster - Заводской комплекс (Несколько брокеров, объединенных в кластер)
* Topic - Способ организации сообщений определенной категории внутри Kafkа. Это Цех с лентами в заводском комплексе.
* Partition - "Лента внутри цеха" - (Конкретный транзакционный лог)
* Replication Factor - Сколько копий каждой "ленты" (Партиции) надо поддерживать на заводе.
* InSyncReplicas - Сколько копий каждой "ленты" (Партиции) должно поддерживаться в статусе inSync (отставать от Лидера не более чем на **replica.lag.time.max.ms**)
* UnderReplicatedPartition - Партиции - у которых по какой-то причине отсутствует полный набор реплик.
* Partition key - Маркировка на посылке - чтобы уложить ее в нужную ленту
* Алгоритм партиционирования (Partitoneer) - Способ обработки маркировки на "посылках". Определяется продюсером
* Partition Count - Сколько лент в "Цехе" (Партиций в топике)
* Producer - Тот кто кладет посылки на ленту
* Consumer - Тот кто считывает посылки с ленты
* Offset - Указатель у консюмера какую последнюю посылку в ленте он обработал

## Воркшоп - 3

### Предварительные условия

1. Должен быть установлен docker или podman (docker-compose или podman-compose - соответственно)
2. Должны быть свободны локальные порты:
   2181 - для ZooKeeper
   8084 - для Kafka-UI
   9094 - для доступа к kafka-broker
   9095 - для доступа к kafka-broker
   9096 - для доступа к kafka-broker
3. (Опционально) Если у вас podman вы можете заменять docker на podman в каждой команде или сделайте алиасы в `~.zshrc`

```bash
alias docker=podman
alias docker-compose=podman-compose
```

### Воркшоп Kafa-3

---

#### 1. Перезапустим наш кластер c 3 брокерами:

```bash
docker-compose up -d
```

Должно получиться что-то вроде такого:

```bash
CONTAINER ID   IMAGE                             COMMAND                  CREATED        STATUS        PORTS                                                           NAMES
d24aaf58b7f3   provectuslabs/kafka-ui:latest     "/bin/sh -c 'java $J…"   24 hours ago   Up 24 hours   0.0.0.0:8084->8080/tcp, :::8084->8080/tcp                       kafka-ui
5c9a1e1ee2b8   confluentinc/cp-kafka:7.2.0       "sh -c '((sleep 15 &…"   24 hours ago   Up 24 hours   9092/tcp, 0.0.0.0:9096->9094/tcp, :::9096->9094/tcp             kafka_3
e824309ee473   confluentinc/cp-kafka:7.2.0       "sh -c '((sleep 15 &…"   24 hours ago   Up 22 hours   9092/tcp, 0.0.0.0:9095->9094/tcp, :::9095->9094/tcp             kafka_2
cf072c6e59f2   confluentinc/cp-kafka:7.2.0       "sh -c '((sleep 15 &…"   24 hours ago   Up 22 hours   9092/tcp, 0.0.0.0:9094->9094/tcp, :::9094->9094/tcp             kafka_1
5960d3960b08   confluentinc/cp-zookeeper:7.2.0   "/etc/confluent/dock…"   24 hours ago   Up 24 hours   2888/tcp, 0.0.0.0:2181->2181/tcp, :::2181->2181/tcp, 3888/tcp   zookeeper
```

2) Зайдем в kafka-ui http://localhost:8084/
3) И посмотрим на вкладку brokers - должно быть 3 брокера в списке.
4) Теперь мы готовы двигаться дальше.

#### 2. Создадим топик test с Replication Factor 3:

1) Add Topic
   2) Name: Test
   3) Number of partitions: 3
   4) Min In Sync Replicas: 1
   5) Replication Factor: 3
   6) Create Topic!
2) Посмотрим на вкладку Overview - Увидим что у каждой партиции на вкладке replicas появилось 3 цифры. 1 из них зеленая - так помечается лидер.

---

#### 3. Запустим продюсер и попробуем записать сообщение - в топик - при выходе брокера из кластера:

```bash
docker run -it --rm --network kafka-3_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

```bash
kafka-console-producer --topic test --bootstrap-server kafka_1:9092 --request-required-acks 0 --timeout 100 --property parse.key=true --property key.separator=":"
```

Запишем несколько сообщений в виде ключ:значение
3)  Теперь нужно остановить broker_2

```bash
docker stop kafka_2
```

4) Попробуем еще записать несколько сообщений - все должно работать.
5) Зайдем в UI - и посмотрим - Должны светиться красны URP  - и In Sync Replicas.
6) Вернем брокер в рабочее состояние выполнив:

```bash
docker-comppose up -d
```

7) Зайдем в UI - и посмотрим - URP  - и In Sync Replicas - должны прекратить светиться красным.

---

#### 4. Запустим продюсер и попробуем записать сообщение - в топик - при выходе брокера из кластера - но включим ack=1

1) Запустим продюсер с режимом AckLeader

```bash
docker run -it --rm --network kafka-3_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

```bash
kafka-console-producer --topic test --bootstrap-server kafka_1:9092 --timeout 100 --request-required-acks 1 --property parse.key=true --property key.separator=":"
```

2) Остановим 1 из брокеров.

```bash
docker stop kafka_2
```

3) Запишем сообщение в топик. - Успешно записалось
4) Остановим еще 1 брокер.
5) Запишем сообщение в топик. - Успешно записалось
6) Посмотрим на UI - увидим что в лидерах у всех партиций - остался последний живой брокер.

---

#### 5. Запустим продюсер и попробуем записать сообщение - в топик - при выходе 2 брокеров из кластера и

1) Запустим консольный продюсер

```bash
docker run -it --rm --network kafka-3_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

```bash
kafka-console-producer --topic test --bootstrap-server kafka_1:9092 --timeout 100 --request-required-acks -1 --property parse.key=true --property key.separator=":"
```

5) Остановим 1 из брокеров

```bash
docker stop kafka_2
```

6) Попробуем записать сообщения
7) Все еще все работает.

#### 6. Запустим продюсер и попробуем записать сообщение - в топик - при выходе брокера из кластера, но с minISR =3:

1) Зайдем в kafka-ui http://localhost:8084/
2) Перейдем во вкладку kafka_1->topics->test->Settings
3) Выставим параметр minISyncReplicas = 3 и нажмем кнопку **save**
4) Запустим консольный продюсер

```bash
docker run -it --rm --network kafka-3_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

```bash
kafka-console-producer --topic test --bootstrap-server kafka_1:9092 --timeout 100 --property parse.key=true --property key.separator=":"
```

5) Остановим 1 из брокеров

```bash
docker stop kafka_2
```

6) Попробуем записать сообщения - получим сообщения об ошибке.
7) Уменьшим minISR - до 2
8) Запишем сообщение - ошибки пропали

### Статьи по теме


[Про идемпотентность](https://habr.com/ru/company/tinkoff/blog/481784/)

[Оф дока по параметрам](https://docs.confluent.io/platform/current/installation/configuration/producer-configs.html#connections-max-idle-ms)

[Еще про идемпотентность](https://medium.com/@shesh.soft/kafka-idempotent-producer-and-consumer-25c52402ceb9)

[Доклад История одной борбы HL 2019](https://www.youtube.com/watch?v=m5CDfrQLzrs&t=1341s)

Прикольный курс от разработчика по KSQLD - рек посмотреть - для вдохновения и понимания возможностей kafka.
[Бонус - KSQLDB](https://www.youtube.com/playlist?list=PL5T99fPsK7pqg59oKuDIYu3MxcIk7Fh-V)
