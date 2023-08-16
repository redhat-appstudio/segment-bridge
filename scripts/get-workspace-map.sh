#!/bin/bash
# get-workspace-map.sh
#   Read Namespace resources and generate a JSON object mapping from namespaces
#   to workspace based on the value of label toolchain.dev.openshift.com/space.
#   This script assumes `oc` is preconfigured with all the required clusters
#   and could be found in $PATH.
#
#   if the CONTEXTS environment variable is set, the script will query
#   locally-defined contexts with those names. Otherwise, the script will query
#   contexts associated with a hard-coded list of staging clusters.
#
set -o pipefail -o errexit -o nounset

# URLs for staging member clusters to be queried if contexts were not provided
# explicitly. A local context with the same URL must exist for each element below (can be
# named arbitrarily).
# I.e. for each url u, kubeconfig must have an entry c on the list named contexts for
# which c.context.cluster == u
stage_clusters=$(cat <<- EOM
  [
    "api-stone-stg-m01-7ayg-p1-openshiftapps-com:6443",
    "api-stone-stg-rh01-l2vh-p1-openshiftapps-com:6443"
  ]
EOM
)

# ======= Parameters ======
# The following variables can be set from outside the script by setting
# similarly named environment variables.
#
# Locally-defined context names to be queried (space-separated)
read -r -a CONTEXTS <<< "${CONTEXTS:-""}"
#
# === End of parameters ===

if [[ ${#CONTEXTS[@]} -eq 0 ]]; then
  mapfile -t CONTEXTS < <(
    oc config view -o=json \
    | jq --argjson clusters "$stage_clusters" \
    '.contexts | .[] | select(.context.cluster as $c | $clusters | index($c)) | .name'
  )
fi

printf "%s\n" "${CONTEXTS[@]}" | xargs -r --replace=C \
  oc --context=C get namespaces \
    --selector=toolchain.dev.openshift.com/type=tenant \
    -o=go-template \
    --template=$'{
      {{- $comma := false}}
      {{- range .items}}
          {{- $namespace := .metadata.name }}
          {{- $workspace := index .metadata.labels "toolchain.dev.openshift.com/space" }}
          {{- if $comma}},{{end}}
          {{- $comma = true -}}
          "{{$namespace}}":"{{ $workspace }}"
      {{- end -}}}\n' \
  | jq --slurp 'add'
