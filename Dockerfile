FROM golang:1.21-alpine as builder

ENV CGO_ENABLED=1 GO111MODULE=on GOOS=linux

# install make
RUN apk update && \
    apk upgrade && \
    apk add --update --no-cache alpine-sdk linux-headers make

WORKDIR /go/src/github.com/settlus/chain

COPY . .

RUN go mod download

# Dockerfile Cross-Compilation Guide
# https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide
ARG TARGETOS TARGETARCH

# install simapp, remove packages
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM alpine:3

EXPOSE 26656 26657 1317 9090 8545

WORKDIR /app

RUN apk add --no-cache ca-certificates curl bash jq

COPY --from=builder /go/src/github.com/settlus/chain/build/settlusd /usr/local/bin/

ENTRYPOINT ["settlusd"]
