# First stage: Build the Go binaries
FROM registry.redhat.io/rhel8/go-toolset:1.19.10-3 AS builder
WORKDIR /opt/app-root/src
COPY --chown=default:root . .
RUN go build -o /opt/app-root/build/ ./cmd/...

# Second stage: Create the final container image
FROM registry.redhat.io/openshift4/ose-tools-rhel8:v4.13.0-202306070816.p0.g05d83ef.assembly.stream

COPY --chown=root:root --chmod=644 data/ca-trust/* /etc/pki/ca-trust/source/anchors
RUN /usr/bin/update-ca-trust
COPY --chown=root:root --chmod=755 scripts/* /usr/local/bin
COPY --chown=root:root --chmod=755 --from=builder /opt/app-root/build/* /usr/local/bin/

# While the scripts already have defaults for the following, specifying them
# here too for sake of documenting in the Dockerfile which variables affect the
# image
#
ENV SPLUNK_API_URL="https://splunk-api.corp.redhat.com:8089"
ENV SPLUNK_APP_NAME="rh_rhtap"
ENV SPLUNK_INDEX="federated:rh_rhtap_stage_audit"
ENV QUERY_EARLIEST_TIME="-4hours"
ENV QUERY_LATEST_TIME="-0hours"
ENV SEGMENT_BATCH_API="https://api.segment.io/v1/batch"
ENV SEGMENT_RETRIES="3"
ENV CURL_NETRC="/usr/local/etc/netrc"
ENV KUBECONFIG="/usr/local/etc/kube_config"
