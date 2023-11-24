#!/bin/bash

# Check if Podman container ID is provided as an argument
CONTAINER_ID="$1"
if [ -z "$CONTAINER_ID" ]; then
    echo "Error: Please provide a container ID as an argument." >&2
    exit 1
fi

CLUSTERS=("host" "m01" "rh01")
KUBECONFIG_DIR="/tmp/kube/"
KUBECONFIG=""

mkdir -p "$KUBECONFIG_DIR" || { echo "Error: Failed to create directory $KUBECONFIG_DIR." >&2; exit 1; }

for CLUSTER_NAME in "${CLUSTERS[@]}"; do
    LOCAL_KUBECONFIG="${KUBECONFIG_DIR}/${CLUSTER_NAME}"
    if podman cp "${CONTAINER_ID}:/${CLUSTER_NAME}" "$LOCAL_KUBECONFIG"; then
        KUBECONFIG="${KUBECONFIG:+$KUBECONFIG:}$LOCAL_KUBECONFIG"
    else
        echo "Warning: Failed to copy kubeconfig for $CLUSTER_NAME from $CONTAINER_ID" >&2
    fi
done

[ -n "$KUBECONFIG" ] || { echo "Error: No kubeconfig files were found." >&2; exit 1; }

MERGED_KUBECONFIG="${KUBECONFIG_DIR}kubeconfig"
kubectl config view --merge --flatten > "$MERGED_KUBECONFIG" || { echo "Error: Failed to merge kubeconfig files." >&2; exit 1; }

echo "Kubeconfig merged and can be found in $MERGED_KUBECONFIG"
