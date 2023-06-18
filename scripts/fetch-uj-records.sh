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

SPLUNK_API_URL=https://splunk-api.corp.redhat.com:8089
SPLUNK_APP_NAME=rh_rhtap
SPLUNK_APP_API_URL="$SPLUNK_API_URL/servicesNS/nobody/$SPLUNK_APP_NAME"
SPLUNK_APP_SEARCH_URL="$SPLUNK_APP_API_URL/search/v2/jobs/export"

FIELDS=(
  auditID impersonatedUser.username objectRef.resource objectRef.namespace
  objectRef.apiGroup objectRef.apiVersion objectRef.name verb
  requestReceivedTimestamp
)
INDEX="federated:rh_rhtap_stage_audit" 
SSO_LOOKUP="rhtap_staging_lookup"

SQ_EV_SELECTOR='search 
  index="'"$INDEX"'"
  log_type=audit 
  NOT verb IN (get, watch, list, deletecollection) 
  "objectRef.apiGroup" IN ("toolchain.dev.openshift.com", "appstudio.redhat.com", "tekton.dev")
  "impersonatedUser.username"="*"
  '
SQ_DEDUP_FIELDS="eval dummy=0$(for F in "${FIELDS[@]}"; do echo -n ",$F=mvindex('$F', 0)"; done)"
SQ_GEN_JSON="eval
  messageId=auditID,
  timestamp=requestReceivedTimestamp,
  type=\"track\",
  userId='impersonatedUser.username',
  event_verb=verb,
  event_subject='objectRef.resource',
  properties=json_object(
    \"apiGroup\", 'objectRef.apiGroup',
    \"apiVersion\", 'objectRef.apiVersion',
    \"kind\", 'objectRef.resource',
    \"namespace\", 'objectRef.namespace',
    \"name\", 'objectRef.name'
  )
  | fields messageId, timestamp, type, userId, event_verb, event_subject, properties | fields - _*
"

QUERY="$SQ_EV_SELECTOR | $SQ_DEDUP_FIELDS | $SQ_GEN_JSON"

set -o xtrace
curl -k -n \
    "$SPLUNK_APP_SEARCH_URL" \
    -d output_mode=json \
    -d search="$QUERY"
