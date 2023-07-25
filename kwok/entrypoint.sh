#!/bin/sh
kwokctl start cluster --name "$CLUSTER_NAME"
kwokctl --name "$CLUSTER_NAME" kubectl proxy --port="${KWOK_KUBE_APISERVER_PORT}" --accept-hosts='^*$' --address="0.0.0.0"
