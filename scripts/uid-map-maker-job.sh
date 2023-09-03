#!/bin/bash
# uid-map-maker-job.sh
#   Generate the UID map and write it to a ConfigMap object. This script assumes
#   it is running in a POD with permissions to write ConfigMaps. KUBECONFIG_SRC
#   should point to a file with data about how to connect to the RHTAP clusters.
#   We also expect to have configuration in place for connecting to the local
#   cluster.
#
#   This script is meant for when running things in some background job, it may
#   add things like monitoring and log formatting.
#
set -o pipefail -o errexit -o nounset -o xtrace

oc create configmap --dry-run=client -o yaml uid-map \
  --from-file=uid-map.json=<(KUBECONFIG="$KUBECONFIG_SRC" get-uid-map.sh) \
| oc apply -f -
