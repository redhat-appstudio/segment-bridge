#!/bin/bash
#  check-same-go-toolset.sh
#    Check that we use the same Go Toolset image in our container build and
#    while running our tests in CI.
#
set -o pipefail -o errexit -o nounset

BUILDER_FILE="Dockerfile"
ACTION_FILE=".github/workflows/unit_tests.yaml"

BUILDER_IMAGE=$(sed -nre 's/^FROM\s+(.*)\s+AS builder$/\1/p' "$BUILDER_FILE")
ACTION_IMAGE=$(sed -nre 's/^\s*image:\s+(.*)\s*$/\1/p' "$ACTION_FILE")

if [[ -z $BUILDER_IMAGE ]]; then
    echo 1>&2 "Go toolset image not found in '$BUILDER_FILE'"
    exit 1
fi

if [[ -z $ACTION_IMAGE ]]; then
    echo 1>&2 "Go toolset image not found in '$ACTION_FILE'"
    exit 1
fi

if [[ "$ACTION_IMAGE" != "$BUILDER_IMAGE" ]]; then
    echo 1>&2 "Go Toolset images in '$BUILDER_FILE' and '$ACTION_FILE' do not match"
    echo 1>&2 "  got '$BUILDER_IMAGE' in '$BUILDER_FILE'"
    echo 1>&2 "  got '$ACTION_IMAGE' in '$ACTION_FILE'"
    exit 1
fi

echo "The same Go toolset image is used everywhere"
exit 0
