#!/bin/bash

namespace="toolchain-host-operator"
kwokctl --name "$CLUSTER_NAME" kubectl create namespace "$namespace"

for i in $(seq 1 10)
do
email="user$i@gmail.com"
email_hash=$(printf "%s" "$email" | md5sum | awk '{ print $1 }')

cat <<EOF | kwokctl --name "$CLUSTER_NAME" kubectl apply -f -
apiVersion: toolchain.dev.openshift.com/v1alpha1
kind: UserSignup
metadata:
  annotations:
    toolchain.dev.openshift.com/user-email: "$email"
    toolchain.dev.openshift.com/last-target-cluster: localcluster.openshiftapps.com
    toolchain.dev.openshift.com/sso-account-id: "5254247$i"
    toolchain.dev.openshift.com/sso-user-id: "5254247$i"
    toolchain.dev.openshift.com/user-email: "user$i@gmail.com"
  labels:
    toolchain.dev.openshift.com/email-hash: "$email_hash"
  name: user$i
  namespace: toolchain-host-operator
spec:
  userid: "f:528d74ff-f703-47ed-9cd5-f0ct$i:user$i"
  username: "user$i"
EOF
done
