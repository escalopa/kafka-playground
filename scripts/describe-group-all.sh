#!/usr/bin/bash

BROKERS_LIST=$1
group_names=$(kafka-consumer-groups.sh --bootstrap-server "$BROKERS_LIST" --list)

# Iterate through each group name
for group_name in $group_names; do
	kafka-consumer-groups.sh --bootstrap-server "$BROKERS_LIST" --describe --group "$group_name"
done