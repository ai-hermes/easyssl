ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.22 AS builder

WORKDIR /src/server

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrate ./cmd/migrate

FROM ${REGISTRY}/alpine:3.20

WORKDIR /app/server

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/api /usr/local/bin/api
COPY --from=builder /out/migrate /usr/local/bin/migrate
COPY --from=builder /src/server/migrations ./migrations
COPY server/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
