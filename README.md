# Access Control 权限管理系统

基于 go-zero + Docker + Kubernetes 的微服务权限管理系统

## 🚀 特性

- **高性能**: 基于 Go-Zero 框架，支持高并发
- **容器化**: 完整的 Docker 支持，一键部署
- **K8s 原生**: 支持 Kubernetes 部署，自动扩缩容
- **监控完善**: 集成 Prometheus + Grafana 监控
- **负载均衡**: Nginx 反向代理，自动负载均衡
- **安全可靠**: JWT 认证，接口限流

## 📁 项目结构

```
access-control/
├── main.go                      # 程序入口
├── etc/
│   ├── config.yaml             # 本地开发配置
│   ├── config-docker.yaml      # Docker 配置
│   └── config-prod.yaml        # 生产环境配置
├── internal/
│   ├── config/                 # 配置结构体
│   ├── handler/                # HTTP 处理器
│   ├── logic/                  # 业务逻辑
│   ├── middleware/             # 中间件
│   ├── model/                  # 数据模型
│   ├── svc/                    # 服务上下文
│   └── types/                  # 类型定义
├── deploy/
│   ├── k8s/                    # Kubernetes 部署文件
│   ├── nginx/                  # Nginx 配置
│   ├── prometheus/             # Prometheus 配置
│   ├── grafana/                # Grafana 配置
│   └── sql/                    # 数据库脚本
├── Dockerfile                   # Docker 构建文件 (开发)
├── Dockerfile.prod             # Docker 构建文件 (生产)
├── docker-compose.yml          # Docker Compose (开发)
├── docker-compose.prod.yml     # Docker Compose (生产)
├── deploy.ps1                  # 部署脚本 (Windows)
├── deploy.sh                   # 部署脚本 (Linux/Mac)
├── Makefile                    # 构建脚本
└── README.md
```

## 🛠 快速开始

### 环境要求

- Go 1.20+
- Docker 20.0+
- Docker Compose 2.0+
- PostgreSQL 14+ (或使用 Docker)
- Redis 6+ (或使用 Docker)

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/yourname/access-control.git
cd access-control

# 2. 安装依赖
go mod tidy

# 3. 启动数据库和 Redis (使用 Docker)
docker-compose up -d postgres redis

# 4. 运行应用
go run main.go -f etc/config.yaml
# 或
make run
```

### Docker 部署 (开发环境)

```bash
# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### Docker 部署 (生产环境)

```bash
# 1. 复制环境变量配置
cp .env.example .env

# 2. 修改 .env 文件中的配置 (必须修改!)
#    - DB_PASSWORD: 数据库密码
#    - JWT_SECRET: JWT 密钥
#    - GRAFANA_PASSWORD: Grafana 密码

# 3. 启动所有服务
docker-compose -f docker-compose.prod.yml up -d

# 或使用部署脚本 (Windows)
.\deploy.ps1 start

# 或使用部署脚本 (Linux/Mac)
./deploy.sh start
```

### Kubernetes 部署

```bash
# 1. 修改 Secret 配置
vim deploy/k8s/secret.yaml

# 2. 部署到 K8s
make deploy-k8s

# 或使用脚本
./k8s-deploy.sh deploy

# 3. 查看状态
make status-k8s
```

## 📊 监控访问

| 服务 | 地址 | 说明 |
|------|------|------|
| API | http://localhost:8080 | API 服务 |
| Nginx | http://localhost:80 | 反向代理 |
| Prometheus | http://localhost:9091 | 指标监控 |
| Grafana | http://localhost:3000 | 可视化面板 |

Grafana 默认账号: `admin` / `admin123`

## 🔧 常用命令

### Make 命令

```bash
make help           # 查看所有命令
make build          # 构建应用
make test           # 运行测试
make lint           # 代码检查
make docker-prod    # 构建生产镜像
make docker-deploy  # 部署生产环境
make db-backup      # 备份数据库
```

### 部署脚本命令 (Windows)

```powershell
.\deploy.ps1 start      # 启动服务
.\deploy.ps1 stop       # 停止服务
.\deploy.ps1 restart    # 重启服务
.\deploy.ps1 update     # 更新服务
.\deploy.ps1 logs       # 查看日志
.\deploy.ps1 status     # 查看状态
.\deploy.ps1 backup     # 备份数据库
```

## 📡 API 接口

### 用户登录

```bash
POST /api/user/login
Content-Type: application/json

{
    "phoneNumber": "18888888888",
    "password": "123456"
}
```

### 响应示例

```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIs...",
        "accessExpire": 1769075546,
        "userInfo": {
            "id": 1,
            "phoneNumber": "18888888888",
            "createdAt": "2026-01-21 08:26:44",
            "updatedAt": "2026-01-21 08:26:44"
        }
    }
}
```

## 🔒 安全配置

生产环境必须修改以下配置：

1. **JWT 密钥**: 使用强随机密钥
   ```bash
   openssl rand -base64 32
   ```

2. **数据库密码**: 使用强密码

3. **Redis 密码**: 如果暴露在公网，必须设置密码

## 📈 性能优化

- 使用连接池管理数据库连接
- Redis 缓存热点数据
- Nginx 限流保护后端服务
- K8s HPA 自动扩缩容

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 License

MIT License
