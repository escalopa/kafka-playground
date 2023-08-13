#!/usr/bin/bash

BROKERS_LIST=$1
topic_names=$(kafka-topics.sh --bootstrap-server "$BROKERS_LIST" --list --exclude-internal)

# Iterate through each topic name
for topic_name in $topic_names; do
	kafka-topics.sh --bootstrap-server "$BROKERS_LIST" --describe --topic "$topic_name"
done