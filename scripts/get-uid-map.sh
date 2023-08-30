#!/bin/bash
# get-uid-map.sh
#   Read UserSignup resources and generate a JSON object mapping from cluster
#   usernames to SSO UDIs.
#   This script assumes `oc` is preconfigured to connect to the right cluster
#   and could be found in $PATH
#
set -o pipefail -o errexit -o nounset

oc get \
  -n toolchain-host-operator \
  UserSignup \
  -o=go-template \
  --template=$'{ 
    {{- $comma := false}}
    {{- range .items}}
      {{- if $uid := (index .metadata.annotations "toolchain.dev.openshift.com/sso-user-id")}}
        {{- if $uname := (.status.compliantUsername)}}
          {{- if $comma}},{{end}}
          {{- $comma = true -}}
          "{{$uname}}":{{$uid}}
        {{- end}}
      {{- end}}
    {{- end}}}\n' 
