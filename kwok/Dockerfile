FROM registry.k8s.io/kwok/cluster:v0.3.0-k8s.v1.27.3

ARG CLUSTER_NAME=kwok
ENV CLUSTER_NAME=${CLUSTER_NAME}

COPY --chmod=755 kwok/user-signups.sh kwok/entrypoint.sh kwok/toolchain.dev.openshift.com_usersignups.yaml ./
COPY --chown=root:root --chmod=755 scripts/* /usr/local/bin

RUN KWOK_KUBE_APISERVER_PORT=0 kwokctl create cluster --name "$CLUSTER_NAME" || exit 1 && \
    kwokctl --name="$CLUSTER_NAME" kubectl apply -f toolchain.dev.openshift.com_usersignups.yaml && \
    sh ./user-signups.sh && \
    rm user-signups.sh toolchain.dev.openshift.com_usersignups.yaml && \
    kwokctl stop cluster --name="$CLUSTER_NAME"

ENTRYPOINT ["/entrypoint.sh"]
