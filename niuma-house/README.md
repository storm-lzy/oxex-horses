# 牛马之家 (NiuMa House)

职场人的避风港，分享工作经验，曝光坑逼公司！🐴

## 项目架构

```
niuma-house/
├── config/                  # 配置文件
│   ├── config.yaml         # 主配置
│   ├── rbac_model.conf     # Casbin RBAC 模型
│   └── rbac_policy.csv     # Casbin 权限策略
├── server/                  # Go 后端
│   ├── cmd/main.go         # 入口
│   ├── internal/           # 业务代码
│   │   ├── model/          # GORM 实体
│   │   ├── service/        # 业务逻辑
│   │   ├── handler/        # HTTP 控制器
│   │   ├── repository/     # 数据访问层
│   │   ├── middleware/     # 中间件
│   │   ├── router/         # 路由
│   │   ├── mq/             # RabbitMQ
│   │   ├── ws/             # WebSocket
│   │   └── task/           # Cron 任务
│   └── pkg/                # 工具包
├── web-client/             # Vue3 用户端
└── web-admin/              # Vue3 管理端
```

## 技术栈

### 后端
- **Go 1.21+** + Gin + GORM v2
- **MySQL 8.0** + Redis 7
- **RabbitMQ** (经验值异步处理)
- **MinIO** (对象存储)
- **Gorilla WebSocket** (实时私信)
- **robfig/cron** (定时任务)
- **Casbin** (RBAC 权限)

### 前端
- **Vue 3** + TypeScript + Vite
- **Pinia** (状态管理)
- **Element Plus** (UI 组件)
- **md-editor-v3** (Markdown 编辑器)
- **ECharts** (管理端图表)

## 快速开始

### 1. 启动基础设施

确保本地已安装并运行：
- MySQL 8.0 (创建数据库 `niuma_house`)
- Redis 7
- RabbitMQ
- MinIO

或使用 Docker Compose:
```bash
docker-compose up -d
```

### 2. 启动后端

```bash
cd server
go mod tidy
go run cmd/main.go -config ../config/config.yaml
```

服务运行在 `http://localhost:8080`

### 3. 启动客户端前端

```bash
cd web-client
npm install
npm run dev
```

访问 `http://localhost:3000`

### 4. 启动管理端前端

```bash
cd web-admin
npm install
npm run dev
```

访问 `http://localhost:3001`

## 默认账号

### 管理员
- 用户名: `admin`
- 密码: `admin123`

## API 概览

### 公开 API
| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/auth/register` | POST | 用户注册 |
| `/api/auth/login` | POST | 用户登录 |
| `/api/occupations` | GET | 获取职业列表 |

### 认证 API (需 Bearer Token)
| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/user/profile` | GET | 获取用户资料 |
| `/api/posts` | GET/POST | 帖子列表/创建 |
| `/api/posts/:id` | GET/PUT/DELETE | 帖子详情/编辑/删除 |
| `/api/posts/:id/like` | POST/DELETE | 点赞/取消 |
| `/api/companies` | GET/POST | 公司列表/添加 |
| `/api/companies/search` | GET | 搜索公司 |
| `/api/upload/presign` | POST | 获取上传预签名 URL |
| `/ws/chat` | WebSocket | 私信连接 |

### 管理 API (需管理员权限)
| 端点 | 方法 | 说明 |
|------|------|------|
| `/admin/dashboard/stats` | GET | 统计数据 |
| `/admin/users` | GET | 用户列表 |
| `/admin/users/:id/ban` | POST | 封禁用户 |
| `/admin/posts` | GET | 帖子管理 |
| `/admin/companies` | GET | 公司管理 |

## 等级系统

| 等级 | 名称 | 所需经验 |
|------|------|----------|
| Lv.1 | 普通牛马 | 0 |
| Lv.2 | 内卷牛马 | 100 |
| Lv.3 | 精英牛马 | 500 |
| Lv.4 | 天选牛马 | 2000 |
| Lv.5 | 核动力牛马 | 10000 |

**经验获取:**
- 发布帖子: +5 EXP
- 获得点赞: +2 EXP
- 获得评论: +1 EXP

## License

MIT
