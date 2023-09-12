FROM registry.k8s.io/kwok/cluster:v0.3.0-k8s.v1.27.3

ARG CLUSTER_NAME=kwok
ENV CLUSTER_NAME=${CLUSTER_NAME}

COPY --chmod=755 user-signups.sh entrypoint.sh toolchain.dev.openshift.com_usersignups.yaml ./

RUN KWOK_KUBE_APISERVER_PORT=0 kwokctl create cluster --name "$CLUSTER_NAME" || exit 1 && \
    kwokctl --name="$CLUSTER_NAME" kubectl apply -f toolchain.dev.openshift.com_usersignups.yaml && \
    sh ./user-signups.sh && \
    rm user-signups.sh toolchain.dev.openshift.com_usersignups.yaml && \
    kwokctl stop cluster --name="$CLUSTER_NAME"

ENTRYPOINT ["/entrypoint.sh"]
