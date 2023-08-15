#!/bin/bash
# main-job.sh
#   Combine the scripts into a single pipeline flowing data from Splunk to
#   Segment.
#   This script is meant for when running things in some background job, it may
#   add things like monitoring and log formatting.
#
set -o pipefail -o errexit -o nounset -o xtrace

fetch-uj-records.sh | splunk-to-segment.sh  | segment-mass-uploader.sh
