FROM dweomer/hashibase as hashibase

WORKDIR /tmp

ARG VAULT_VERSION=0.10.2

ADD https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS .
ADD https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS.sig .
ADD https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip .

RUN gpg --verify vault_${VAULT_VERSION}_SHA256SUMS.sig vault_${VAULT_VERSION}_SHA256SUMS
RUN grep linux_amd64 vault_${VAULT_VERSION}_SHA256SUMS | sha256sum -cs
RUN unzip vault_${VAULT_VERSION}_linux_amd64.zip -d /usr/local/bin
RUN vault version

FROM alpine

ARG VAULT_GID=8201
ARG VAULT_UID=8200

RUN set -x \
 && apk add --no-cache \
    dumb-init \
    libcap \
    su-exec \
 && addgroup -g ${VAULT_GID} vault \
 && adduser -S -G vault -u ${VAULT_UID} vault

COPY --from=hashibase /usr/local/bin/* /usr/local/bin/

# USER vault
ENTRYPOINT ["dumb-init", "--", "su-exec", "vault:vault", "vault"]
CMD ["help"]
