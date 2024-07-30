FROM docker.io/splunk/splunk@sha256:ec4479b588a6c6edf62bdc812fac92a563a1f0384dc63f9dbad55c35340fa763 AS builder

ENV SPLUNK_START_ARGS=--accept-license

USER root

COPY --chmod=777 tests/test_logs /opt/splunk/test_logs
COPY --chmod=777 tests/indexes.conf /opt/splunk/etc/apps/search/local/indexes.conf
COPY --chmod=777 scripts/log_indexing.sh /opt/splunk/log_indexing.sh
COPY --chmod=777 scripts/pre-run.sh /opt/splunk/pre-run.sh
COPY --chmod=777 props.conf /opt/splunk/etc/system/local/
COPY --chmod=777 limits.conf /opt/splunk/etc/system/local/


RUN microdnf install -y jq --disablerepo="*" --enablerepo="ubi-8-appstream-rpms" \
    && microdnf clean all \
    && umask 000 \
    && touch "/opt/splunk/output.log"

USER ansible
RUN bash -o errexit /opt/splunk/pre-run.sh

FROM docker.io/splunk/splunk@sha256:ec4479b588a6c6edf62bdc812fac92a563a1f0384dc63f9dbad55c35340fa763

ENV SPLUNK_START_ARGS=--accept-license

EXPOSE 8000
EXPOSE 8089

USER root

RUN mkdir -p /opt/splunk/var/lib/splunk
COPY --chmod=777 --from=builder /opt/splunk/etc/apps/search/local/indexes.conf /opt/splunk/etc/apps/search/local/indexes.conf
COPY --chmod=777 --from=builder /var/splunk_buildtime_db/splunk/ /opt/splunk/var/lib/splunk/

ENV SPLUNKD_SSL_ENABLE=false


USER ansible
ENTRYPOINT ["/sbin/entrypoint.sh"]
CMD ["start-service"]
