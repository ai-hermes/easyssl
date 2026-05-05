# EasySSL CLI

A command-line interface for managing EasySSL certificates and workflows.

## Installation

```bash
go build -o easyssl ./cmd/easyssl
```

## Usage

### Login

```bash
easyssl login --server http://localhost:8090 --email admin@easyssl.local --password "your-password"
```

### Config

```bash
# View current config
easyssl config get

# Set API key for OpenAPI endpoints
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
