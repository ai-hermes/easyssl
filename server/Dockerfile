ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.25 AS builder

WORKDIR /src/server

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

ARG GIT_BRANCH=unknown
ARG GIT_COMMIT=unknown
ARG GIT_REPO_URL=
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X easyssl/server/internal/version.Branch=${GIT_BRANCH} -X easyssl/server/internal/version.Commit=${GIT_COMMIT} -X easyssl/server/internal/version.RepoURL=${GIT_REPO_URL}" \
    -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/migrate ./cmd/migrate

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
