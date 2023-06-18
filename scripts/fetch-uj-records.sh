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
#
# === End of parameters ===

SPLUNK_APP_API_URL="$SPLUNK_API_URL/servicesNS/nobody/$SPLUNK_APP_NAME"
SPLUNK_APP_SEARCH_URL="$SPLUNK_APP_API_URL/search/v2/jobs/export"

FIELDS=(
  auditID impersonatedUser.username user.username objectRef.resource
  objectRef.namespace objectRef.apiGroup objectRef.apiVersion objectRef.name
  verb requestReceivedTimestamp
)

SQ_EV_SELECTOR='search 
  index="'"$SPLUNK_INDEX"'"
  log_type=audit 
  NOT verb IN (get, watch, list, deletecollection) 
  "objectRef.apiGroup" IN ("toolchain.dev.openshift.com", "appstudio.redhat.com", "tekton.dev")
  ("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*"))
  '
SQ_DEDUP_FIELDS="eval dummy=0$(for F in "${FIELDS[@]}"; do echo -n ",$F=mvindex('$F', 0)"; done)"
SQ_GEN_JSON="eval
  messageId=auditID,
  timestamp=requestReceivedTimestamp,
  type=\"track\",
  userId=if(isnull('impersonatedUser.username'), 'user.username', 'impersonatedUser.username'),
  event_verb=verb,
  event_subject='objectRef.resource',
  properties=json_object(
    \"apiGroup\", 'objectRef.apiGroup',
    \"apiVersion\", 'objectRef.apiVersion',
    \"kind\", 'objectRef.resource',
    \"name\", 'objectRef.name'
  )
  | fields messageId, timestamp, type, userId, event_verb, event_subject, properties | fields - _*
"

QUERY="$SQ_EV_SELECTOR | $SQ_DEDUP_FIELDS | $SQ_GEN_JSON"

set -o xtrace
curl --netrc \
  "$SPLUNK_APP_SEARCH_URL" \
  --data output_mode=json \
  --data earliest_time="$QUERY_EARLIEST_TIME" \
  --data latest_time="$QUERY_LATEST_TIME" \
  --data search="$QUERY"
