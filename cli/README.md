# EasySSL CLI

A command-line interface for managing EasySSL certificates and workflows.

## Installation

```bash
go build -o easyssl ./cmd/easyssl
```

## Usage

### Login

```bash
easyssl login --api-key <your-api-key>
```

Default backend: `https://easyssl.spotty.com.cn/`

Use a custom backend for debugging:

```bash
easyssl login --server http://localhost:8090 --api-key <your-api-key>
easyssl config set server http://localhost:8090
```

### Config

```bash
# View current config
easyssl config get

# Set API key manually
easyssl config set api_key <your-api-key>
```

### Workflows

```bash
# List workflows
easyssl workflow list
```

### Certificates

```bash
# List certificates
easyssl certificate list

# Apply for a certificate via OpenAPI
easyssl certificate apply --workflow <workflow-id> --api-key <key>

# Check certificate run status
easyssl certificate status <run-id> --api-key <key>

# Download a certificate
easyssl certificate download <certificate-id> --api-key <key> --output cert.zip
```

### Logout

```bash
easyssl logout
```
