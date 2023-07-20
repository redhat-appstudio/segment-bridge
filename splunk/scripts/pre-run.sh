#!/bin/bash

export SPLUNK_PASSWORD="Password"

# Run an instance of Splunk during the build, to prepare the indexed logs for the tests
/sbin/entrypoint.sh start-service >> $SPLUNK_HOME/output.log 2>&1 &

# Make sure the Splunk service is up
timeout_start=$(date +%s)
while true; do
    if grep -q "Ansible playbook complete" $SPLUNK_HOME/output.log; then
        break
    fi
    if [ "$(($(date +%s) - timeout_start))" -ge 120 ]; then
        echo "The Splunk instance is not up, cannot verify indexing for tests." >&2
        exit 1
    fi
    sleep 5
done

# Index the logs required for the tests
bash /opt/splunk/log_indexing.sh

# Copy the whole Splunk DB to an accessible directory
sudo cp -r /opt/splunk/var/lib/splunk/ /opt/splunk/splunk

# Shut down the Splunk service
sudo -u splunk /opt/splunk/bin/splunk stop