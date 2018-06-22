FROM dweomer/hashibase as hashibase

WORKDIR /tmp

ARG CONSUL_VERSION=1.1.0

ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS .
ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS.sig .
ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip .

RUN gpg --verify consul_${CONSUL_VERSION}_SHA256SUMS.sig consul_${CONSUL_VERSION}_SHA256SUMS
RUN grep linux_amd64.zip consul_${CONSUL_VERSION}_SHA256SUMS | sha256sum -cs
RUN unzip consul_${CONSUL_VERSION}_linux_amd64.zip -d /usr/local/bin

FROM alpine

ARG CONSUL_GID=8300
ARG CONSUL_UID=8500

RUN set -x \
 && apk add --no-cache \
    coreutils \
    dumb-init \
    jq \
    su-exec \
 && addgroup -g ${CONSUL_GID} consul \
 && adduser -S -G consul -u ${CONSUL_UID} consul
COPY --from=hashibase /usr/local/bin/* /usr/local/bin/

# USER consul
ENTRYPOINT ["dumb-init", "--", "su-exec", "consul:consul", "consul"]
CMD ["help"]
