FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /nexusforge ./cmd/nexusforge

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && addgroup -S nexus && adduser -S -G nexus nexus
COPY --from=builder /nexusforge /usr/local/bin/nexusforge
USER nexus
EXPOSE 8080
ENV NEXUSFORGE_ADDR=0.0.0.0:8080
ENTRYPOINT ["/usr/local/bin/nexusforge"]
