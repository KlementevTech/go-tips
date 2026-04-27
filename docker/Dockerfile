ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.21
ARG VERSION=unknown

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags="-X main.version=${VERSION} -s -w" \
    -o /tmp/gotips \
    ./cmd/server

FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S gotips && adduser -S gotips -G gotips

RUN mkdir -p /etc/gotips /data/gotips && \
    chown -R gotips:gotips /etc/gotips /data/gotips

COPY --from=builder --chown=gotips:gotips /tmp/gotips /bin/gotips
COPY --from=builder --chown=gotips:gotips /app/config/config.dev.toml /etc/gotips/config.toml

VOLUME ["/data/gotips"]

USER gotips

ENTRYPOINT ["/bin/gotips"]
CMD ["-c", "/etc/gotips/config.toml"]