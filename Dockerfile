ARG GOLANG_VERSION=1.21
ARG ALPINE_VERSION=3.20

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV PROJECT=tf-summarize

WORKDIR ${PROJECT}

COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY . .
RUN go build -a -o /${PROJECT} .

### Base image with shell
FROM alpine:${ALPINE_VERSION} AS base-release
RUN apk --update --no-cache add ca-certificates && update-ca-certificates

### Build with goreleaser
FROM base-release AS goreleaser
COPY tf-summarize /bin/
ENTRYPOINT ["/bin/tf-summarize"]

### Build in docker
FROM base-release AS release
COPY --from=builder /tf-summarize /bin/
ENTRYPOINT ["/bin/tf-summarize"]

### Scratch with build in docker
FROM scratch AS scratch-release
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tf-summarize /bin/
ENTRYPOINT ["/bin/tf-summarize"]
USER 65534

### Scratch with goreleaser
FROM scratch AS scratch-goreleaser
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY tf-summarize /bin/
ENTRYPOINT ["/bin/tf-summarize"]
USER 65534
