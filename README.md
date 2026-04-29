# SMS Dashboard 📱

一个轻量级、自托管的短信转发接收面板。配合 Android 短信转发器使用，可以方便地在网页端查看、管理转发的短信。

## ✨ 特性

- **单文件部署**：后端使用 Go 开发，前端构建产物通过 `embed` 嵌入，最终只生成一个二进制文件。
- **Android 转发器兼容**：提供标准的 HTTP POST 接口，完美对接各类 Android 短信转发应用。
- **安全可靠**：
  - API 接收接口支持 Token 鉴权。
  - 管理面板支持账号密码登录，使用 JWT 进行身份验证。
- **简洁美观**：前端采用 Vue 3 + Tailwind CSS 构建，支持响应式布局。
- **轻量存储**：使用 SQLite 数据库，无需安装复杂的数据库环境。
- **持久化**：支持 Docker 部署及数据卷挂载。

## 🛠️ 技术栈

- **后端**: Go (Gin, GORM)
- **前端**: Vue 3, TypeScript, Vite, Tailwind CSS
- **数据库**: SQLite

## 🚀 快速开始

### 本地开发

1. **克隆仓库**
   ```bash
   git clone https://github.com/your-username/sms-dashboard.git
   cd sms-dashboard
   ```

2. **前端构建**
   ```bash
   cd web
   npm install
   npm run build
   cd ..
   ```

3. **启动后端**
   ```bash
   go run cmd/server/main.go
   ```
   访问 `http://localhost:8080`，默认账号 `admin`，默认密码 `admin`。

### Docker 部署 (推荐)

我们推荐使用 Docker 进行生产环境部署，以实现数据的持久化。

**单端口模式（前后端一体）:**
```bash
docker run -d \
  -p 8080:8080 \
  -v /your/local/path/data:/app/data \
  -e API_TOKEN=your_secret_token \
  --name sms-dashboard \
  ghcr.io/jjss520/sms-dashboard-dd:latest
```

**双端口模式（前后端分离）:**
```bash
docker run -d \
  -p 8080:8080 \
  -p 3000:3000 \
  -v /your/local/path/data:/app/data \
  -e API_TOKEN=your_secret_token \
  -e PORT=8080 \
  -e WEB_PORT=3000 \
  --name sms-dashboard \
  ghcr.io/jjss520/sms-dashboard-dd:latest
```

- **前端面板**: `http://localhost:3000`
- **API 服务**: `http://localhost:8080`
- **默认账号**: `admin` / `admin`

## 📲 Android 转发器配置

在你的 Android 短信转发器 App 中，配置发送方：

- **URL**: `http://your-server-ip:8080/api/sms`
- **Method**: `POST`
- **参数说明**:
  - `content`: 短信内容（必填）
  - `sendTime`: 短信时间（可选，格式：2026-04-29 09:16:05）
  - `sender`: 发件人号码（可选）
  - `phone`: 卡槽/手机号（可选，如：SIM2_13345694791）
  - `device`: 设备型号（可选，如：MIX 3）
  - `token`: API Token（或使用查询参数 `?token=xxx`）
- **Content-Type**: `application/json` 或 `application/x-www-form-urlencoded`
## 📄 开源协议

MIT License
