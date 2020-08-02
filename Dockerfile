# Build code
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o dns "./cmd/server"
WORKDIR /dist

RUN cp /build/dns .

# Build a small image
FROM scratch

ENV DNS_PORT=8081 \
    DNS_ENVIRONMENT=PRODUCTION \
    DNS_LOG_LEVEL=INFO \
    DNS_SECTOR_ID=1

COPY --from=builder /dist/dns /

ENTRYPOINT ["/dns"]