#!/bin/bash
# mk-segment-batch-payload.sh
#   Accept Segment event records on STDIN (One per-line) and wrap them together
#   in a sigle JSON document suitable as payload for the Segment batch API.
#
set -o pipefail -o errexit -o nounset

jq \
  --compact-output \
  --slurp \
  '{batch: .}'
