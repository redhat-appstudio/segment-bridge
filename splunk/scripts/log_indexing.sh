#!/bin/bash

logs_path="$SPLUNK_HOME/test_logs"
endpoint="https://localhost:8089/services/data/inputs/oneshot"
SPLUNK_PASSWORD="Password"

echo "Indexing logs..."
max_retries=3

for file in "$logs_path"/*; do
  echo "Indexing $file"
  if ! curl --insecure --user admin:"$SPLUNK_PASSWORD" --retry $max_retries --retry-max-time 60 "$endpoint" --data name="$file" --data index="test_index"; then
    echo "Encountered an error while initiating indexing for $file. Max retry attempts reached. Exiting."
    exit 1
  fi
done

timeout_start=$(date +%s)
while true; do
    if curl --insecure --user admin:"$SPLUNK_PASSWORD" "$endpoint" --get --data output_mode=json \
    | jq -e '.paging.total == 0'
    then
        break
    fi
    if [[ "$(($(date +%s) - timeout_start))" -ge 60 ]]
    then
        echo "Log indexing did not complete within the expected time. Exiting." >&2
        exit 1
    fi
    echo "Waiting for log indexing to complete." >&2
    sleep 5
done


echo "Logs indexing finished successfully"
