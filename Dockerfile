############################
# STEP 1 build executable plugin binary
############################
FROM golang:1.16-buster AS builder

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	ca-certificates upx-ucl

WORKDIR /plugin

ENV GO111MODULE=on
COPY go.mod go.sum ./
COPY LICENSE ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -v -o quorum-hashicorp-vault-plugin
RUN upx quorum-hashicorp-vault-plugin
RUN sha256sum -b quorum-hashicorp-vault-plugin | cut -d' ' -f1 > SHA256SUM

############################
# STEP 2 build new vault image
############################
FROM library/vault:1.8.4

RUN apk add --no-cache \
    jq \
    curl

# Expose the plugin directory as a volume
VOLUME /vault/plugins

COPY --from=builder /plugin/LICENSE /
COPY --from=builder /plugin/quorum-hashicorp-vault-plugin /vault/plugins/quorum-hashicorp-vault-plugin
COPY --from=builder /plugin/scripts/vault-init.sh /usr/local/bin/vault-init.sh
COPY --from=builder /plugin/scripts/vault-init-dev.sh /usr/local/bin/vault-init-dev.sh

EXPOSE 8200
