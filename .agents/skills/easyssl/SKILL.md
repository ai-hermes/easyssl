---
name: easyssl
description: Use EasySSL CLI in agent workflows for certificate apply, run polling, access discovery, and certificate download with strict JSON outputs.
license: MIT
metadata:
  author: easyssl
  version: "1.0.0"
---

# EasySSL CLI Skill

## When to use
- Need to request certificate operations from EasySSL in terminal-driven agents.
- Need machine-readable outputs for chaining steps.
- Need troubleshooting via CLI diagnostics.

## Required command contract
- Always run with `--format json`.
- For diagnostics add `--verbose` (and optionally `--trace`).
- Parse JSON only from stdout; errors appear on stderr.
- Authenticate with API key only (`X-API-Key`); do not rely on username/password login flows.
- Before running any business subcommand, verify login state first (`whoami` preflight).

## Initialization (binary + executable path)
Run once per session before any EasySSL subcommand:

```bash
set -euo pipefail

INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "unsupported arch: ${ARCH}" >&2; exit 1 ;;
esac

VERSION="${EASYSSL_VERSION:-latest}" # e.g. v0.1.1
BASE_URL="${EASYSSL_BASE_URL:-https://gh-proxy.org/https://github.com/ai-hermes/easyssl/releases}"
BIN_NAME="easyssl-${OS}-${ARCH}"
TARGET_BIN="${INSTALL_DIR}/easyssl"

if [ ! -x "${TARGET_BIN}" ]; then
  curl -fL "${BASE_URL}/download/${VERSION}/${BIN_NAME}" -o "${TARGET_BIN}"
  chmod +x "${TARGET_BIN}"
fi

# Option A: append PATH in current shell
export PATH="${INSTALL_DIR}:$PATH"
# Option B: use full path directly
EASYSSL_BIN="${EASYSSL_BIN:-${TARGET_BIN}}"
```

## Login preflight (must run before every subcommand)
Before `access/certificate/workflow/...` commands, check login state:

```bash
set -euo pipefail

EASYSSL_BIN="${EASYSSL_BIN:-easyssl}"
SERVER="${EASYSSL_SERVER:-https://easyssl.spotty.com.cn}"

if ! "${EASYSSL_BIN}" whoami --server "${SERVER}" --format json >/dev/null 2>&1; then
  : "${EASYSSL_API_KEY:?EASYSSL_API_KEY is required when not logged in}"
  "${EASYSSL_BIN}" login --server "${SERVER}" --api-key "${EASYSSL_API_KEY}" --format json >/dev/null
  "${EASYSSL_BIN}" whoami --server "${SERVER}" --format json >/dev/null
fi
```

## Typical workflow
1. Initialize binary and command path (run once):
```bash
# Run the initialization block above first, then:
EASYSSL_BIN="${EASYSSL_BIN:-easyssl}"
```
2. Login preflight (run before every business subcommand):
```bash
# Run the login preflight block above first, then:
"${EASYSSL_BIN}" whoami --server "${SERVER}" --format json
```
3. List access credentials:
```bash
"${EASYSSL_BIN}" access list --openapi --server "${SERVER}" --format json
```
4. Apply certificate:
```bash
"${EASYSSL_BIN}" certificate apply --provider tencent --access-id <access-id> --domain example.com --domain www.example.com --server "${SERVER}" --format json
```
5. Poll run status:
```bash
"${EASYSSL_BIN}" certificate status <run-id> --server "${SERVER}" --format json
```
6. Inspect run events when needed:
```bash
"${EASYSSL_BIN}" certificate events <run-id> --limit 100 --server "${SERVER}" --format json
```
7. Download certificate metadata/package:
```bash
"${EASYSSL_BIN}" certificate download <certificate-id> --openapi --cert-format PEM --server "${SERVER}" --format json
```

## Troubleshooting
- Use `easyssl version --format json` to capture binary metadata.
- Use `--verbose` to log request/response envelope to stderr.
- Use `--trace` only in secure environments (response snippet logging).
- If stderr contains `unexpected response from server` with HTML-like snippets:
  - Verify `--server` points to EasySSL API host root (not a web app URL/path fallback).
  - Verify `X-API-Key` is set and valid.
  - Retry with `--verbose --trace` and capture `request-id` for backend troubleshooting.
