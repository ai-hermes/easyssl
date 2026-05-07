# EasySSL

Certimate-style SSL automation system rebuilt with:

- Backend: Go + Gin + PostgreSQL
- Frontend: React + Vite + TailwindCSS + shadcn-style UI primitives
- Deployment: single backend binary can serve frontend static assets

## Current implementation (V1 base)

- JWT admin auth (`/api/auth/*`)
- Access CRUD (`/api/accesses`)
- Workflow CRUD + runs (`/api/workflows`, `/api/workflows/:id/runs`)
- Workflow dispatcher queue with cancel/stats
- Certificate list/download/revoke (`/api/certificates/*`)
- Statistics + notification test endpoints
- PostgreSQL migration runner (`go run ./cmd/migrate`)

Default bootstrap admin:

- email: `admin@easyssl.local`
- password: `1234567890`

## Quick start

### 0) Prepare global config

```bash
cp .env.example .env
```

Edit `.env` for your global settings (especially `PG_*` and `JWT_SECRET`).

### 1) Start DB + API

```bash
docker compose up -d db
cd server
GOPROXY=https://goproxy.cn,direct go run ./cmd/migrate
GOPROXY=https://goproxy.cn,direct go run ./cmd/api
```

### 2) Start web

```bash
cd web
npm install
npm run dev
```

Visit: `http://127.0.0.1:5173`

## OpenAPI (X-API-Key) quick flow

1. Login and create API Key:
   - `POST /api/auth/login`
   - `POST /api/auth/api-keys` (returns one-time `token`)
2. Use `X-API-Key` to apply certificate:

```bash
curl -X POST "http://127.0.0.1:8090/api/open/certificates/apply" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: esk_xxx" \
  -d '{
    "provider": "tencentcloud",
    "accessId": "your-access-id",
    "domains": ["ssl1.example.com", "*.example.com"],
    "caProvider": "letsencrypt",
    "keyAlgorithm": "RSA2048"
  }'
```

3. Poll run status and events:
   - `GET /api/open/certificates/runs/{runId}`
   - `GET /api/open/certificates/runs/{runId}/events`

Swagger UI: `http://127.0.0.1:8090/swagger/index.html`

## Install EasySSL Skill / 安装 EasySSL Skill

Expand one language section below to switch:

<details open>
<summary><strong>中文说明</strong></summary>

仓库内置 skill 目录：`.agents/skills/easyssl`。

1) 安装到本地 skills 目录：

```bash
mkdir -p ~/.codex/skills
rm -rf ~/.codex/skills/easyssl
cp -R .agents/skills/easyssl ~/.codex/skills/easyssl
```

2) 验证安装：

```bash
ls ~/.codex/skills/easyssl/SKILL.md
```

在 Codex 对话中可直接提到 `easyssl` skill 触发使用。

</details>

<details>
<summary><strong>English</strong></summary>

Built-in skill source path: `.agents/skills/easyssl`.

1) Install to local skills directory:

```bash
mkdir -p ~/.codex/skills
rm -rf ~/.codex/skills/easyssl
cp -R .agents/skills/easyssl ~/.codex/skills/easyssl
```

2) Verify installation:

```bash
ls ~/.codex/skills/easyssl/SKILL.md
```

Mention `easyssl` in Codex chat to trigger this skill.

</details>

## Provider 配置案例

先在「Access 授权管理」创建授权，工作流节点里的 `accessId` 引用对应授权记录 `id`。

### Access: Aliyun

```yaml
name: aliyun-main
provider: aliyun
config:
  accessKeyId: "LTAIxxxx"
  accessKeySecret: "xxxx"
  region: "cn-hangzhou"        # optional
  resourceGroupId: "rg-xxxx"   # optional
```

### Access: TencentCloud

```yaml
name: tencent-dns
provider: tencentcloud
config:
  secretId: "AKIDxxxx"
  secretKey: "xxxx"
  region: "ap-guangzhou"       # optional, default ap-guangzhou
  sessionToken: ""             # optional
```

### Access: Qiniu

```yaml
name: qiniu-ssl
provider: qiniu
config:
  accessKey: "xxxx"
  secretKey: "xxxx"
```

### Access: SSH

```yaml
name: prod-nginx
provider: ssh
config:
  host: "10.0.0.12"
  port: 22
  username: "root"
  authMethod: "password"       # password or key
  password: "xxxx"             # required when authMethod=password
  # key: "-----BEGIN OPENSSH PRIVATE KEY-----..."   # required when authMethod=key
  # keyPassphrase: "xxxx"                            # optional
```

## Workflow YAML 案例

### 1) 仅申请证书（不部署）

```yaml
version: 1
options:
  failFast: true
nodes:
  - id: apply-1
    name: Apply Cert
    action: apply
    provider: tencentcloud
    accessId: "your-tencent-access-id"
    config:
      domains: "ssl1.example.com;*.example.com"
      caProvider: "letsencrypt"
      contactEmail: ""              # optional, auto-generate if empty
      keyAlgorithm: "RSA2048"       # optional: RSA2048/RSA4096/EC256/EC384
      dnsPropagationTimeout: 120    # optional
      dnsTTL: 0                     # optional
edges: []
```

### 2) 申请（腾讯云）+ 部署到 Aliyun CAS

```yaml
version: 1
options:
  failFast: true
nodes:
  - id: apply-1
    name: Apply By Tencent
    action: apply
    provider: tencentcloud
    accessId: "tencent-access-id"
    config:
      domains: "ssl1.example.com"
      caProvider: "letsencrypt"

  - id: deploy-aliyun
    name: Deploy To Aliyun CAS
    action: deploy
    provider: aliyun-cas
    accessId: "aliyun-access-id"
    config:
      region: "cn-hangzhou"         # optional
      resourceGroupId: "rg-xxxx"    # optional
edges:
  - source: apply-1
    target: deploy-aliyun
```

### 3) 申请（阿里云 DNS）+ 部署到 Qiniu

```yaml
version: 1
options:
  failFast: true
nodes:
  - id: apply-1
    name: Apply By Aliyun DNS
    action: apply
    provider: aliyun
    accessId: "aliyun-access-id"
    config:
      domains: "ssl1.example.com"
      caProvider: "letsencrypt"

  - id: deploy-qiniu
    name: Deploy To Qiniu
    action: deploy
    provider: qiniu
    accessId: "qiniu-access-id"
    config:
      certName: "easyssl-ssl1-example-com"  # optional, auto-generate if empty
      commonName: "ssl1.example.com"        # optional
edges:
  - source: apply-1
    target: deploy-qiniu
```

### 4) 申请 + 部署到 SSH (Nginx)

```yaml
version: 1
options:
  failFast: true
nodes:
  - id: apply-1
    name: Apply Cert
    action: apply
    provider: tencentcloud
    accessId: "tencent-access-id"
    config:
      domains: "ssl1.example.com"
      caProvider: "letsencrypt"

  - id: deploy-ssh
    name: Deploy To Nginx
    action: deploy
    provider: ssh
    accessId: "ssh-access-id"
    config:
      certPath: "/etc/nginx/ssl/fullchain.pem"
      keyPath: "/etc/nginx/ssl/privkey.pem"
      preCommand: "mkdir -p /etc/nginx/ssl"   # optional
      postCommand: "nginx -s reload"          # optional
edges:
  - source: apply-1
    target: deploy-ssh
```

## Project structure

- `server/`: Gin API, domain services, workflow dispatcher, migrations
- `web/`: Vite React app, dashboard/access/workflow/certificate/settings pages
- `docker-compose.yml`: local PostgreSQL
- `.env`: global runtime config (local file, not committed)
- `.env.example`: global config template

## Notes

- Workflow execution currently includes a mock certificate issuance step for end-to-end verification.
- Provider adapters are scaffold-ready via `provider` and workflow node config but not yet fully expanded to all target providers.
