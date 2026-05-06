# EasySSL CLI

A command-line interface for managing EasySSL certificates and workflows.

## Build

```bash
cd cli
go build -o easyssl ./cmd/easyssl
```

## Core usage

```bash
# login (api-key)
easyssl login --api-key <key>

# include diagnostics
easyssl --verbose whoami

# version metadata
easyssl version

# list resources
easyssl access list --openapi
easyssl certificate list --openapi
easyssl workflow list

# apply cert via openapi
easyssl certificate apply --provider tencent --access-id <id> --domain example.com

# run/event polling
easyssl certificate status <run-id>
easyssl certificate events <run-id> --limit 100
```

## Output and diagnostics

- Default output format: `json`.
- Supported formats: `--format json|text`.
- Verbose diagnostics: `--verbose`.
- Include response body snippets: `--trace`.
- When server returns HTML/fallback pages, CLI emits a structured error to stderr with status/content-type/request URL hints.

## Config keys

```bash
easyssl config get
easyssl config set server https://easyssl.spotty.com.cn
easyssl config set api_key <key>
easyssl config set output json
easyssl config set timeout 30
easyssl config set trace false
```
