#!/bin/bash
# fetch-uj-records.sh
#   Fetch user journey records from Splunk
#   Records look almost like Segment track records but some further adjustments
#   are required such as:
#   - Converting nested JSON objects from strings to objects
#   - Mapping cluster usernames to SSO user IDs
#   This script assumes that credentials are preconfigured for curl for
#   connecting to Splunk in a .netrc file
#
set -o pipefail -o errexit -o nounset
#
# ======= Parameters ======
# The following variables can be set from outside the script by setting
# similarly named environment variables.
#
# The API URL to use for connecting to Splunk
SPLUNK_API_URL="${SPLUNK_API_URL:-https://splunk-api.corp.redhat.com:8089}"
# The Splunk app name that may store custom data objects we may use for the
# query
SPLUNK_APP_NAME="${SPLUNK_APP_NAME:-rh_rhtap}"
# The Splunk index to fetch data from
SPLUNK_INDEX="${SPLUNK_INDEX:-federated:rh_rhtap_stage_audit}"
# Specify the earliest time to retrieve records from
# Value is a Splunk time string, defaults to 4 hours ago
QUERY_EARLIEST_TIME="${QUERY_EARLIEST_TIME:-"-4hours"}"
# Specify the latest time to retrieve records from
# Value is a Splunk time string, defaults to now
QUERY_LATEST_TIME="${QUERY_LATEST_TIME:-"-0hours"}"
# How many times to retry calling Splunk
SPLUNK_RETRIES="${SPLUNK_RETRIES:-3}"
#
# A .netrc file to load credentials from
CURL_NETRC="${CURL_NETRC:-$HOME/.netrc}"
#
# === End of parameters ===

SPLUNK_APP_API_URL="$SPLUNK_API_URL/servicesNS/nobody/$SPLUNK_APP_NAME"
SPLUNK_APP_SEARCH_URL="$SPLUNK_APP_API_URL/search/v2/jobs/export"

GO_PACKAGE="github.com/redhat-appstudio/segment-bridge.git"

if command -v querygen > /dev/null; then
  QUERYGEN=querygen
elif command -v go > /dev/null; then
  QUERYGEN="go run $GO_PACKAGE/cmd/querygen"
else
  echo "Couldn\`t find the querygen binary or go in $PATH" 1>&2
  exit 127
fi

$QUERYGEN -0 --index="$SPLUNK_INDEX" \
  | xargs -0 --no-run-if-empty -iQ curl --netrc-file "$CURL_NETRC" \
    --fail --fail-early \
    "$SPLUNK_APP_SEARCH_URL" \
    --retry "$SPLUNK_RETRIES" \
    --data output_mode=json \
    --data earliest_time="$QUERY_EARLIEST_TIME" \
    --data latest_time="$QUERY_LATEST_TIME" \
    --data-urlencode search=Q
