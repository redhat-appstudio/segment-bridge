#!/bin/bash
# segment-uploader.sh
#   Get Segemnt event records from STDIN (One per-line) and POST them to the
#   Segment batch API. The API is limited to max request size of 500KB, so no
#   more then that should be sent as input.
#   This script assumes we have a .netrc file with suitable credentials for
#   Segment
#
set -o pipefail -o errexit -o nounset

# Add script file directory to PATH so we can use other scripts in the same
# directory
SELFDIR="$(dirname "$0")"
PATH="$SELFDIR:${PATH#$SELFDIR:}"

# ======= Parameters ======
# The following variables can be set from outside the script by setting
# similarly named environment variables.
#
# The Segment API URL to use:
SEGMENT_BATCH_API="${SEGMENT_BATCH_API:-https://api.segment.io/v1/batch}"
# How many times to rety calling into Segment
SEGMENT_RETRIES="${SEGMENT_RETRIES:-3}"
#
# === End of parameters ===

set -o xtrace
mk-segment-batch-payload.sh | \
  curl --netrc \
    "$SEGMENT_BATCH_API" \
    --header "Content-Type: application/json" \
    --fail \
    --retry "$SEGMENT_RETRIES" \
    --data @- \
