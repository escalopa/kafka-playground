# [](https://)Kafka workshop-1

## Основные термины (Шпаргалка)

* Broker - Инстанс кафки (Завод)
* Cluster - Заводской комплекс (Несколько брокеров, объединенных в кластер)
* Topic - Способ организации сообщений определенной категории внутри Kafkа. Это Цех с лентами в заводском комплексе.
* Partition - "Лента внутри цеха" - (Конкретный транзакционный лог)
* Replication Factor - Сколько копий каждой "ленты" (Партиции) надо поддерживать на заводе.
* Partition key - Маркировка на посылке - чтобы уложить ее в нужную ленту
* Алгоритм партиционирования (Partitoneer) - Способ обработки маркировки на "посылках". Определяется продюсером
* Partition Count - Сколько лент в "Цехе" (Партиций в топике)
* Producer - Тот кто кладет посылки на ленту
* Consumer - Тот кто считывает посылки с ленты
* Offset - Указатель у консюмера какую последнюю посылку в ленте он обработал

## Воркшоп - 1

### Предварительные условия

1. Должен быть установлен docker или podman (docker-compose или podman-compose - соответственно)
2. Должны быть свободны локальные порты:
   2181 - для ZooKeeper
   8084 - для Kafka-UI
   9094 - для доступа к kafka-broker
3. (Опционально) Если у вас podman вы можете заменять docker на podman в каждой команде или сделайте алиасы в `~.zshrc`

```bash
alias docker=podman
alias docker-compose=podman-compose
```

Не забываем сделать:

```bash
docker login gitlab-registry.ozon.ru
```

Для тех у кого по какой-то причине не получилось залогиниться в gitlab:

```bash
docker-compose -f docker-compose.yaml up -d
```

### Воркшоп Kafa-1

#### 1. Запустим локальный Kafka-кластер:

##### 1.1 Выполним clone репозитория и запустим kafka-cluster с 1 брокером:

```
git clone git@gitlab.ozon.ru:route-mentors/kafka-workshops.git
cd kafka-workshops/kafka-1
docker-compose up -d
```

Если все запустилось успешно, то команда `docker ps` вернет вам что-то типа этого:

```
CONTAINER ID  IMAGE                                      COMMAND               CREATED       STATUS           PORTS                   NAMES
00deb322df8b  docker.io/confluentinc/cp-zookeeper:7.1.0  /etc/confluent/do...  16 hours ago  Up 1 minutes ago  0.0.0.0:2181->2181/tcp  zookeeper
dba22b7f7400  docker.io/confluentinc/cp-kafka:7.1.0      sh -c ((sleep 15 ...  16 hours ago  Up 1 minutes ago  0.0.0.0:9094->9094/tcp  kafka_1
cfbc12e5f1a0  docker.io/provectuslabs/kafka-ui:latest    /bin/sh -c java $...  16 hours ago  Up 1 minutes ago  0.0.0.0:8084->8080/tcp  kafka-ui
```

##### 1.2 Перейдем в kafka-ui

1. Откроем в браузере http://localhost:8084/
2. Перейдем во вкладку Brokers
3. Перейдем во владку Topics
   1. Add Topic
      1. Name: test
      2. Number of Partitions: 1
   2. Produce Message
      1. Value: 123
      2. Produce Message
   3. Messages
      1. Кликнуть на сообщение - посмотреть что там за поля
   4. ...-> Remove Topic -> Confirm

Про дополнительные параметры топика можно прочитать [тут](https://docs.confluent.io/platform/current/installation/configuration/topic-configs.html).

### 2. Создадим топик order через консоль

Запустим контейнер с башем для OS: Mac, Linux:

```bash
docker run -it --rm -v $(pwd):/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Запустим контейнер с башем для OS: Windows:

```bash
docker run -it --rm -v ${pwd}:/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```


Внутри контейнера запустите:

```bash
kafka-topics --bootstrap-server kafka_1:9092 --create --topic order
```

#### 3. Зальем сообщения


Запустим контейнер с башем для OS: Mac, Linux:

```bash
docker run -it --rm -v $(pwd):/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Запустим контейнер с башем для OS: Windows:

```bash
docker run -it --rm -v ${pwd}:/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```


Внутри контейнера запустим:


```bash
kafka-console-producer --topic order --bootstrap-server kafka_1:9092 --timeout 100 < /data/example_messages_1.txt
```


#### 4. Отредактируем Топик через kafka-ui (можем сразу его зачистить)

1. Зайдем на http://localhost:8084/ - попадем в kafka-ui
2. Выберем топик order
3. ... -> Edit setting 
4. Number of Partitions: **3**
5. Нажимаем кнопку **Submit**
6. Увидим что у топика order - изменилось количество партиций
7. Зачистим топик - нажав кнопку ...-> clear messages

#### 5. Зальем сообщения с Partition Key


Запустим контейнер с башем для OS: Mac, Linux:

```bash
docker run -it --rm -v $(pwd):/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Запустим контейнер с башем для OS: Windows:

```bash
docker run -it --rm -v ${pwd}:/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Внутри контейнера: 

```bash
kafka-console-producer --topic order --bootstrap-server kafka_1:9092 --property parse.key=true --property key.separator=":" --timeout 100 < /data/example_messages_2.txt
```



#### (Дополнительно) Чтобы подключиться к топику и послушать сообщения через консольного консюмера можно выполнить следующий код:


Запустим контейнер с башем для OS: Mac, Linux:

```bash
docker run -it --rm -v $(pwd):/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Запустим контейнер с башем для OS: Windows:

```bash
docker run -it --rm -v ${pwd}:/data --network kafka-1_kafka-network gitlab-registry.ozon.ru/route-mentors/kafka-workshops/confluentinc/cp-kafka:7.2.0 /bin/bash
```

Запустим консольный продюсер:
```bash
kafka-console-consumer --bootstrap-server kafka_1:9092 --topic order --property print.key=true --property key.separator=":" --from-beginning --group g1
```
