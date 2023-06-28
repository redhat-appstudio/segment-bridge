#!/bin/bash

logs_path="$SPLUNK_HOME/test_logs"
endpoint="https://localhost:8089/services/data/inputs/oneshot"
SPLUNK_PASSWORD="Password"

echo "Indexing logs..."
for file in "$logs_path"/*; do
  echo "Indexing $file"
  curl  --insecure --user admin:"$SPLUNK_PASSWORD" "$endpoint" --data name="$file" --data index="test_index"
done

timeout_start=$(date +%s)
while true; do
    if ! curl --insecure --user admin:"$SPLUNK_PASSWORD" "$endpoint" --get --data output_mode=json \
    | jq -e '.paging.total == 0'
    then
        echo "Waiting for log indexing to complete." >&2
        break
    fi
    if [[ "$(($(date +%s) - timeout_start))" -ge 60 ]]
    then
        echo "Log indexing did not complete within the expected time. Exiting." >&2
        exit 1
    fi
    sleep 5
done


echo "Logs indexing finished successfully"
