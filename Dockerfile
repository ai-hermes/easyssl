# Production Dockerfile: builds both web frontend and Go backend into a single image.
ARG REGISTRY=docker.io

# ------------------------------
# Stage 1: Build web frontend
# ------------------------------
FROM ${REGISTRY}/node:22-alpine AS web-builder
WORKDIR /src/web

COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY web/ ./
RUN pnpm run build

# ------------------------------
# Stage 2: Build Go backend
# ------------------------------
FROM ${REGISTRY}/golang:1.25-alpine AS go-builder
WORKDIR /src/server

RUN apk add --no-cache git

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

ARG GIT_BRANCH=unknown
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X easyssl/server/internal/version.Branch=${GIT_BRANCH} -X easyssl/server/internal/version.Commit=${GIT_COMMIT}" \
    -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/migrate ./cmd/migrate

# ------------------------------
# Stage 3: Production image
# ------------------------------
FROM ${REGISTRY}/alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates wget

COPY --from=go-builder /out/api /usr/local/bin/api
COPY --from=go-builder /out/migrate /usr/local/bin/migrate
COPY --from=go-builder /src/server/migrations ./migrations
COPY --from=go-builder /src/server/docs ./docs
COPY server/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=web-builder /src/web/dist /web/dist

RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8090

HEALTHCHECK --interval=10s --timeout=5s --start-period=15s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
