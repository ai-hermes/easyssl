# EasySSL

[![English](https://img.shields.io/badge/Language-English-blue)](./README.md)
[![简体中文](https://img.shields.io/badge/语言-简体中文-red)](./README.zh-CN.md)

基于 Certimate 风格重构的 SSL 自动化系统，技术栈：

- 后端：Go + Gin + PostgreSQL
- 前端：React + Vite + TailwindCSS + shadcn 风格 UI 组件
- 部署：单个后端二进制可直接托管前端静态资源

## 当前实现（V1 基线）

- JWT 管理员认证（`/api/auth/*`）
- Access 授权 CRUD（`/api/accesses`）
- Workflow CRUD 与执行（`/api/workflows`, `/api/workflows/:id/runs`）
- 工作流调度队列（支持取消与统计）
- 证书列表/下载/吊销（`/api/certificates/*`）
- 统计与通知测试接口
- PostgreSQL 迁移执行器（`go run ./cmd/migrate`）

默认初始化管理员：

- email: `admin@easyssl.local`
- password: `1234567890`

## 快速开始

### 0) 准备全局配置

```bash
cp .env.example .env
```

按需编辑 `.env`（重点关注 `PG_*` 与 `JWT_SECRET`）。

### 1) 启动 DB + API

```bash
docker compose up -d db
cd server
GOPROXY=https://goproxy.cn,direct go run ./cmd/migrate
GOPROXY=https://goproxy.cn,direct go run ./cmd/api
```

### 2) 启动 Web

```bash
cd web
npm install
npm run dev
```

访问：`http://127.0.0.1:5173`

## OpenAPI（X-API-Key）快速流程

1. 登录并创建 API Key：
   - `POST /api/auth/login`
   - `POST /api/auth/api-keys`（返回一次性 `token`）
2. 使用 `X-API-Key` 申请证书：

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

3. 轮询执行状态与事件：
   - `GET /api/open/certificates/runs/{runId}`
   - `GET /api/open/certificates/runs/{runId}/events`

Swagger 地址：`http://127.0.0.1:8090/swagger/index.html`

## 安装 EasySSL Skill（给 Codex/Agents 使用）

仓库内置 skill 目录：`.agents/skills/easyssl`。

选择以下任一安装方式：

**方式一 - 从 GitHub 仓库安装（推荐）：**

```bash
npx skills add ai-hermes/easyssl -y -g
```

**方式二 - 从本地源码安装：**

```bash
mkdir -p ~/.codex/skills
rm -rf ~/.codex/skills/easyssl
cp -R .agents/skills/easyssl ~/.codex/skills/easyssl
```

安装后校验：

```bash
ls ~/.codex/skills/easyssl/SKILL.md
```

在 Codex 对话中可直接提到 `easyssl` skill 触发使用。

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
