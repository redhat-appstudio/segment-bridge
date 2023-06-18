#!/bin/bash
# splunk-to-segment.sh
#   Adapt user journey recordes loaded from Splunk for uploading into Segmment:
#   - Map cluster usernames to SSO user IDs
#   - Convert nested JSON objects from strings to actual objects.
#   - Combine the event_* fields into a single UI-flavoured event string.
#
set -o pipefail -o errexit -o nounset

# Add script file directory to PATH so we can use other scripts in the same
# directory
SELFDIR="$(dirname "$0")"
PATH="$SELFDIR:${PATH#$SELFDIR:}"

function event_verb_map() {
  # Print a JSON map for converting preset tense K8s API server verbs to past
  # tense
  echo '{
    "create": "created",
    "delete": "deleted",
    "deletecollection": "collection deleted",
    "get": "fetched",
    "head": "headers fetched",
    "list": "listed",
    "patch": "patched",
    "update": "updated",
    "watch": "watch started"
  }'
}

function event_subject_map() {
  # Print a JSON map for converting the plural resource names found in the
  # audit log to the singular, capitalized names users are used to seeing
  echo '{
    "applications": "Application",
    "bannedusers": "BannedUser",
    "buildpipelineselectors": "BuildPipelineSelector",
    "componentdetectionqueries": "ComponentDetectionQuery",
    "components": "Component",
    "customruns": "CustomRun",
    "deploymenttargetclaims": "DeploymentTargetClaim",
    "deploymenttargets": "DeploymentTarget",
    "enterprisecontractpolicies": "EnterpriseContractPolicy",
    "environments": "Environment",
    "integrationtestscenarios": "IntegrationTestScenario",
    "internalrequests": "InternalRequest",
    "masteruserrecords": "MasterUserRecord",
    "memberoperatorconfigs": "MemberOperatorConfig",
    "memberstatuses": "MemberStatus",
    "notifications": "Notification",
    "nstemplatesets": "NSTemplateSet",
    "nstemplatetiers": "NSTemplateTier",
    "pipelineresources": "PipelineResource",
    "pipelineruns": "PipelineRun",
    "pipelines": "Pipeline",
    "promotionruns": "PromotionRun",
    "proxyplugins": "ProxyPlugin",
    "releaseplanadmissions": "ReleasePlanAdmission",
    "releaseplans": "ReleasePlan",
    "releases": "Release",
    "releasestrategies": "ReleaseStrategy",
    "remotesecrets": "RemoteSecret",
    "runs": "Run",
    "snapshotenvironmentbindings": "SnapshotEnvironmentBinding",
    "snapshots": "Snapshot",
    "socialevents": "SocialEvent",
    "spacebindings": "SpaceBinding",
    "spacerequests": "SpaceRequest",
    "spaces": "Space",
    "spiaccesschecks": "SPIAccessCheck",
    "spiaccesstokenbindings": "SPIAccessTokenBinding",
    "spiaccesstokendataupdates": "SPIAccessTokenDataUpdate",
    "spiaccesstokens": "SPIAccessToken",
    "spifilecontentrequests": "SPIFileContentRequest",
    "taskruns": "TaskRun",
    "tasks": "Task",
    "tiertemplates": "TierTemplate",
    "toolchainclusters": "ToolChainCluster",
    "toolchainconfigs": "ToolChainConfig",
    "toolchainstatuses": "ToolChainStatus",
    "useraccounts": "UserAccount",
    "usersignups": "UserSignup",
    "usertiers": "UserTier",
    "verificationpolicies": "VerificationPolicy"
  }'
}

jq \
  --slurpfile uidm <(get-uid-map.sh) \
  --slurpfile evvm <(event_verb_map) \
  --slurpfile evsm <(event_subject_map) \
  '.result | select($uidm[0][.userId])
  | {
      messageId, 
      timestamp, 
      type, 
      userId: $uidm[0][.userId], 
      event: "\($evsm[0][.event_subject] // .event_subject) \($evvm[0][.event_verb])", 
      properties: (.properties|fromjson)
    }
  '
