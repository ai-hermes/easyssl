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

## Typical workflow
1. Authenticate:
```bash
easyssl login --server https://easyssl.spotty.com.cn --api-key "$EASYSSL_API_KEY" --format json
```
2. List access credentials:
```bash
easyssl access list --openapi --format json
```
3. Apply certificate:
```bash
easyssl certificate apply --provider tencent --access-id <access-id> --domain example.com --domain www.example.com --format json
```
4. Poll run status:
```bash
easyssl certificate status <run-id> --format json
```
5. Inspect run events when needed:
```bash
easyssl certificate events <run-id> --limit 100 --format json
```
6. Download certificate metadata/package:
```bash
easyssl certificate download <certificate-id> --openapi --cert-format PEM --format json
```

## Troubleshooting
- Use `easyssl version --format json` to capture binary metadata.
- Use `--verbose` to log request/response envelope to stderr.
- Use `--trace` only in secure environments (response snippet logging).
- If stderr contains `unexpected response from server` with HTML-like snippets:
  - Verify `--server` points to EasySSL API host root (not a web app URL/path fallback).
  - Verify `X-API-Key` is set and valid.
  - Retry with `--verbose --trace` and capture `request-id` for backend troubleshooting.
