# Builder
FROM golang:1.25.6-alpine AS builder
RUN apk add --no-cache make nodejs npm gcc musl-dev g++ linux-headers binutils-dev
ADD . /app
WORKDIR /app
COPY config.env .
RUN CGO_ENABLED=1 make all
RUN GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o /build/wallet ./cmd/wallet

# Service
FROM alpine:latest
COPY --from=builder /build/wallet /usr/bin
COPY --from=builder /app/config.env /usr/bin
WORKDIR /usr/bin
ENTRYPOINT ["./wallet"]