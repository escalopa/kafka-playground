# [](https://)Kafka workshop-4

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

## Воркшоп - 4

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

### Воркшоп Kafka-4

Предварительные условия

1) Остановим kafka-cluster из прошлого воркшопа
2) Запустим наш новый кафка-кластер - с настроенным более частым лог-роллером

```bash
docker-compose up -d
```

#### 1) Тестируем RetentionBytes

1) Создадим топик test_topic
   2) Cleanup Policy: Delete
   3) Больше ничего не трогаем.
3) Настроим нужный нам параметр через консоль (через UI его настроить нельзя) 

Запустим терминал с утилитами
```bash
docker run -it --rm --network kafka-4_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

Настроим наш топик на максимальный размер в 10kb и добавим настроек - чтобы триггер происходил почаще
```bash
 kafka-configs --bootstrap-server kafka_1:9092 --entity-name test_topic --entity-type topics --alter --add-config retention.bytes=10240,max.compaction.lag.ms=10000,min.cleanable.dirty.ratio=0.0,segment.ms=10000,delete.retention.ms=10000
```
 
4) Откроем логи kafka_1 в отдельном терминале

```bash
docker logs -f kafka_1
```

5) Давайте начнем заливать туда данные

Откроем терминал в папке kafka-4/bin
Тут лежит 4 бинарника с kafka-пушкой

* **kafka_gun64.exe** - для windows
* **kafka_gun_linux** - для линукс
* **kafka_gun_mac** - Mac Intel
* **kafka_gun_mac_m1** - Mac M1
* **binary.ammo** - файл с патроном для заливки
* **binary-config.yaml** - файл наводки для  kafka-пушки

 Для пользователей Mac/Linux - напоминаю - надо сделать  

Для пользователей Linux/Max
```bash
chmod +x ./kafka_gun*
```

```bash
./kafka_gun_mac ./binary-config.yaml
```

Для пользователей Windows:

```bash
./kafka_gun64.exe ./binary-config.yaml
```

Запустится Наливка в топик test_topic данными из файла binary.ammo

5) Открываем [kafka-ui](http://localhost:8084/ui/clusters/kafka_1/all-topics)
6) Смотрим в терминал где-то в течении минуты должно появится сообщение:
```
Deleting segment LogSegment * due to retention size 10240
`````
Смотрим в kafka-ui в наш тестовый [топик](http://localhost:8084/ui/clusters/kafka_1/all-topics/test_topic):  
Видим увеличение в колонке first Offset 

#### 2) Teстируем RetentionMs

1) Запустим консоль с утилитами:

```bash
docker run -it --rm --network kafka-4_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

В консоли запустим:
```bash
kafka-configs --bootstrap-server kafka_1:9092 --entity-name test_topic --entity-type topics --alter --add-config retention.ms=15000,retention.bytes=10240000000
```

2) Запустим еще раз подачу нагрузки

Для Mac: 
```bash
cd ./kafka-4/bin
./kafka_gun_mac ./binary-config.yaml
```

Для пользователей Windows:

```powershell
cd ./kafka-4/bin
kafka_gun64.exe ./binary-config.yaml
```

Для Linux:


```powershell
cd ./kafka-4/bin
kafka_gun_linux ./binary-config.yaml
```

3) Смотрим в терминал где-то в течении минуты должно появится сообщение:

```
Deleting segment LogSegment(*) due to retention time 15000ms breach based on the largest record timestamp in the segment (kafka.log.UnifiedLog)
```

#### 3) Teстируем Log Compact
1) Давайте сменим Retention Policy на Log Compact
   2) Для этого запустим консоль с утилитами 
   

```bash
docker run -it --rm --network kafka-4_kafka-network confluentinc/cp-kafka:7.2.0 /bin/bash
```

И выполним команду смены типа retention policy

```bash
kafka-configs --bootstrap-server kafka_1:9092 --entity-name test_topic --entity-type topics --alter --add-config cleanup.policy=compact
```

2) Зальем пушкой данные:

```bash
cd ./kafka-4/bin
./kafka_gun_mac ./compact_config.yaml
```
3) Зайдем и посмотрим оставшиеся сообщения
Немножко подождем - и увидим - что в рамках 1 ключа - самы старые записи - остались по 1 штуку. И некоторое количество новых записей не было закомпакчено.



Дополнительное домашнее задание для тех кто хочет поиграться с кафкой глубже:
1) Раскомментировать блоки c ksql-db и ksql-cli в docker-compose.yaml
2) Пройти шаги отсюда: https://ksqldb.io/quickstart.html с пункта 3.
3) Создать свой Stream и Materialized View и поиграться с ними через ksqlDB-cli:
   4) Фильтрация стрима по условию
   5) Создание Materialized View по условию
   5) Квери по Materialized View с условиями
