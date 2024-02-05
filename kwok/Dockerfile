FROM registry.k8s.io/kwok/cluster:v0.5.0-k8s.v1.27.3

ENV HOST_CLUSTER=host
ENV MEMBER_CLUSTER=m01
ENV RH_CLUSTER=rh01

COPY --chmod=755 \
    user-signups.sh \
    toolchain.dev.openshift.com_usersignups.yaml \
    test_member_namespace_map.yml \
    test_rh_namespace_map.yml \
    ./

COPY --chmod=644 inittab /etc/inittab

RUN KWOK_KUBE_APISERVER_PORT=0 kwokctl create cluster --name "$HOST_CLUSTER" || exit 1 && \
    KWOK_KUBE_APISERVER_PORT=0 kwokctl create cluster --name "$MEMBER_CLUSTER" || exit 1 && \
    KWOK_KUBE_APISERVER_PORT=0 kwokctl create cluster --name "$RH_CLUSTER" || exit 1 && \
    kwokctl --name="$HOST_CLUSTER" kubectl apply -f toolchain.dev.openshift.com_usersignups.yaml && \
    CLUSTER_NAME="$HOST_CLUSTER" sh ./user-signups.sh && \
    kwokctl --name="$MEMBER_CLUSTER" kubectl apply -f toolchain.dev.openshift.com_usersignups.yaml && \
    kwokctl --name="$MEMBER_CLUSTER" kubectl apply -f test_member_namespace_map.yml && \
    kwokctl --name="$RH_CLUSTER" kubectl apply -f toolchain.dev.openshift.com_usersignups.yaml && \
    kwokctl --name="$RH_CLUSTER" kubectl apply -f test_rh_namespace_map.yml && \
    kwokctl stop cluster --name="$HOST_CLUSTER" && \
    kwokctl stop cluster --name="$MEMBER_CLUSTER" && \
    kwokctl stop cluster --name="$RH_CLUSTER" && \
    rm user-signups.sh toolchain.dev.openshift.com_usersignups.yaml

ENTRYPOINT ["/sbin/init"]
