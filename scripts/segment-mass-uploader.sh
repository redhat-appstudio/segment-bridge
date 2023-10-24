#!/bin/bash
# segment-mass-uploader.sh
#   Accept a stream of Segment event JSON records from STDIN and pipes them to
#   segment-uploader.sh in chunks suitable for the Segment API (No more than
#   500KB).
#
set -o pipefail -o errexit -o nounset

# Add script file directory to PATH so we can use other scripts in the same
# directory
SELFDIR="$(dirname "$0")"
PATH="$SELFDIR:${PATH#"$SELFDIR":}"

# ======= Parameters ======
# The following variables can be set from outside the script by setting
# similarly named environment variables.
#
# The size of data chunk to send in bytes.
# While the Segment API can accept calls up to 500KB, we send 490KB to leave
# some room for JSON overhead and HTTP headers.
SEGMENT_BATCH_DATA_SIZE="${SEGMENT_BATCH_DATA_SIZE:-$((490 * 1024))}"
#
# === End of parameters ===

split \
  --line-bytes "$SEGMENT_BATCH_DATA_SIZE" \
  --filter="segment-uploader.sh"
