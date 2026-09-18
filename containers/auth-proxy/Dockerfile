FROM envoyproxy/envoy:v1.35.3

USER root
RUN apt-get update \
    && apt-get install --no-install-recommends -y gettext-base \
    && rm -rf /var/lib/apt/lists/*

COPY envoy.yaml.template /etc/envoy/envoy.yaml.template
COPY --chmod=755 entrypoint.sh /usr/local/bin/auth-proxy-entrypoint.sh

USER envoy
ENTRYPOINT ["/usr/local/bin/auth-proxy-entrypoint.sh"]
