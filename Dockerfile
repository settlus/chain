FROM golang:1.21-alpine as builder

ARG GITHUB_TOKEN
ENV CGO_ENABLED=1 GO111MODULE=on GOOS=linux

# install make
RUN apk update && \
    apk upgrade && \
    apk add --update --no-cache alpine-sdk linux-headers make

WORKDIR /app
COPY . .

RUN go mod download
RUN make settlusd

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/build/settlusd /usr/local/bin/

ENTRYPOINT ["settlusd"]
