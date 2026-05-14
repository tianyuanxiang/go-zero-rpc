# go-zero-rpc

基于 **go-zero v1.10** 的微服务架构仓库，集成 RBAC 权限管理、JWT 双Token认证、WebSocket 实时通信、定时任务调度等企业级后台管理系统所需的核心能力。该仓库可直接作为新项目的脚手架，快速启动微服务开发。

## 技术栈介绍

| 技术            | 版本        | 用途                                 |
| --------------- | ----------- | ------------------------------------ |
| **go-zero**     | v1.10.1     | 微服务核心框架（REST + gRPC）        |
| **GORM**        | v1.31.1     | ORM（复杂查询、事务、软删除）        |
| **PostgreSQL**  | 14.x        | 主数据库（pgx v5 驱动）              |
| **Redis**       | go-redis v9 | 缓存、Token黑名单、分布式锁          |
| **Casbin**      | v2.100.0    | RBAC 权限管理（gorm-adapter 持久化） |
| **JWT**         | v5.3.1      | Token 认证（双Token机制）            |
| **etcd**        | v3.5.21     | 服务注册与发现、配置中心             |
| **robfig/cron** | v3.0.1      | 定时任务调度                         |

## 架构设计

```
                    HTTP Clients (Browser / Mobile)
                            |
                    +-------v--------+
                    |   Gateway      |  :8888 (REST API)
                    |  go-zero rest  |
                    +-------+--------+
                            |
                  gRPC (etcd 服务发现)
                            |
          +-----------------+------------------+
          |                                    |
    +-----+-----+                      +------+---+
    |  sys.rpc   |                      |   job    |
    |  gRPC      |  :9100               |  cron    |
    +-----+------+                      +----------+
          |
    +-----+------+
    | PostgreSQL  |
    +-----+------+
          |
    +-----+------+
    |   Redis     |
    +-------------+
```

### 分层设计

```
gateway/
  handler/     -- HTTP 请求解析、参数校验、响应封装
  logic/       -- 业务编排：调用 RPC 客户端，拼装返回数据
  svc/         -- 服务上下文：持有所有 RPC 客户端、中间件、第三方依赖

service/sys/rpc/
  server/      -- gRPC 入口（goctl 生成）
  logic/       -- 核心业务逻辑
  model/       -- 数据访问层（goctl 生成 + GORM 自定义扩展）
  svc/         -- 服务上下文：持有 DB、Redis、Casbin 等依赖
```

## 功能清单

### 认证模块 (AuthService)

- 用户登录（密码 bcrypt 校验，返回 AccessToken + RefreshToken）
- Token 刷新（RefreshToken 换取新 AccessToken）
- 退出登录（AccessToken 加入 Redis 黑名单）
- 修改密码（需提供旧密码）
- 获取当前用户信息（含角色、权限）

### 系统管理 (SystemService)

| 模块         | 功能                                                  |
| ------------ | ----------------------------------------------------- |
| **用户管理** | 创建、编辑、删除（软删除）、列表查询、重置密码        |
| **角色管理** | 创建、编辑、删除、分配菜单权限、分配 API 权限         |
| **菜单管理** | 树形菜单 CRUD、支持目录/菜单/按钮三种类型、可见性控制 |
| **API管理**  | 接口注册、路径+方法唯一性校验                         |
| **字典管理** | 字典类型 CRUD + 字典数据 CRUD                         |
| **文件管理** | 文件上传、列表查询、删除                              |
| **日志管理** | 登录日志查询/清理、操作日志查询/清理                  |

### 权限管理 (PermissionService)

- RBAC 权限校验（Casbin，支持 `keyMatch2` 路径匹配）
- 用户查询（ID、用户名）
- 用户角色查询
- Token 吊销状态检查

### 实时通信

- WebSocket 连接管理（基于 Gorilla WebSocket）
- 用户级连接路由（Hub 模式）
- 心跳保活机制

### 定时任务

- 登录日志定期清理
- 操作日志定期清理
- Redis 分布式锁防止多实例重复执行
- Cron 表达式可配置

## 项目结构

```
go-zero-rpc/
├── common/                          # 共享模块
│   ├── constants/                   #   角色常量（RoleCodeAdmin）
│   ├── jwtx/                        #   JWT 生成与解析
│   ├── middleware/                  #   AuthMiddleware / CasbinMiddleware
│   ├── response/                    #   统一 HTTP 响应格式
│   ├── rpcerr/                      #   gRPC 错误解码（客户端侧）
│   └── xerr/                        #   业务错误码定义（服务端侧）
│
├── gateway/                         # HTTP 网关服务
│   ├── gateway.go                   #   入口
│   ├── api/                         #   .api 定义文件（goctl 生成代码的源）
│   │   └── desc/                    #     按模块拆分的 API 定义
│   ├── etc/                         #   配置文件
│   │   ├── gateway-api.yaml         #     生产配置
│   │   └── gateway-api.example.yaml #     配置模板（git 跟踪）
│   └── internal/
│       ├── config/                  #   配置结构体
│       ├── handler/                 #   HTTP 处理器
│       │   ├── routes.go            #     路由注册
│       │   ├── auth/                #     认证相关 handler
│       │   ├── sys/                 #     系统管理 handler
│       │   └── ws/                  #     WebSocket handler
│       ├── logic/                   #   业务编排逻辑
│       ├── svc/                     #   服务上下文
│       ├── types/                   #   请求/响应类型
│       └── ws/                      #   WebSocket 核心
│
├── service/
│   ├── sys/rpc/                     # 系统 RPC 服务
│   │   ├── sys.go                   #   入口
│   │   ├── client/                  #   生成的 gRPC 客户端（供 gateway 引用）
│   │   ├── pb/                      #   Protobuf 定义 + 生成代码
│   │   ├── etc/                     #   配置文件
│   │   └── internal/
│   │       ├── config/              #   配置结构体
│   │       ├── logic/               #   核心业务逻辑
│   │       │   ├── authservice/     #     认证逻辑
│   │       │   ├── permissionservice/ #   权限逻辑
│   │       │   └── systemservice/   #     系统管理逻辑
│   │       ├── model/               #   数据访问层
│   │       │   ├── *_model_gen.go   #     goctl 生成的 CRUD
│   │       │   └── *_model.go       #     GORM 自定义扩展
│   │       └── server/              #   gRPC 服务端实现
│   │
│   └── job/                         # 定时任务服务
│       ├── job.go                   #   入口
│       ├── etc/                     #   配置文件
│       └── internal/
│           ├── lockx/               #   Redis 分布式锁
│           ├── scheduler/           #   Cron 调度器
│           └── logic/               #   任务执行逻辑
│
├── deploy/                          # 数据库脚本
│   └── system_pgsql.sql                   #   PostgreSQL 初始化 DDL + DML
│
├── etc/                             # 全局配置
│   └── rbac_model.conf              #   Casbin RBAC 模型定义
│
├── go.work                          # Go Workspace 配置
├── go.mod                           # 根模块
└── .gitignore
```

## 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 14+
- Redis 6+
- etcd 3.5+
- goctl（go-zero 代码生成工具）

### 安装 goctl

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.1
```

### 1. 克隆仓库

```bash
git clone https://github.com/tianyuanxiang/go-zero-rpc.git
cd go-zero-rpc
```

### 2. 初始化数据库

在 PostgreSQL 中创建数据库并导入初始化脚本：

```bash
createdb system
psql -d system -f deploy/system_pgsql.sql
```

### 3. 配置服务

参照各服务的 `*.example.yaml` 模板，创建实际配置文件：

**gateway/etc/gateway-api.yaml**
```yaml
Name: gateway
Host: 0.0.0.0
Port: 8888
Mode: dev
Timeout: 30000

Auth:
  AccessSecret: "your-jwt-secret-key"
  AccessExpire: 1200

SysRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sys.rpc
```

**service/sys/rpc/etc/sys.yaml**
```yaml
Name: sys.rpc
ListenOn: 0.0.0.0:9100
Mode: dev

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: sys.rpc

JwtAuth:
  AccessSecret: "your-jwt-secret-key"
  AccessExpire: 1200
  RefreshExpire: 10080

DB:
  DataSource: "postgres://user:password@127.0.0.1:5432/system?sslmode=disable&TimeZone=Asia/Shanghai"

CacheRedis:
  - Host: "127.0.0.1:6379"
    Pass: ""
    Type: node

BizRedis:
  Host: 127.0.0.1:6379
  Pass: ""
  DB: 10

CasbinModelPath: "etc/rbac_model.conf"
UploadPath: "uploads"
```

**service/job/etc/job.yaml**
```yaml
Name: job
Mode: dev

SysRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sys.rpc

BizRedis:
  Host: 127.0.0.1:6379
  Pass: ""
  DB: 10

Jobs:
  ClearLoginLog:
    Enable: true
    Cron: "0 0 3 * * *"
    LockExpireSeconds: 600
  ClearOperLog:
    Enable: true
    Cron: "0 */10 * * * *"
    LockExpireSeconds: 600
```

### 4. 初始化 Casbin 策略

```sql
-- 为管理员角色添加 Casbin 策略（v0 = role_code）
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ('p', 'admin', '/*', '*', '', '', '');
```

### 5. 启动服务

按顺序启动三个服务：

```bash
# 终端 1 - 启动 RPC 服务
cd service/sys/rpc
go run sys.go -f etc/sys.yaml

# 终端 2 - 启动定时任务服务
cd service/job
go run job.go -f etc/job.yaml

# 终端 3 - 启动 API 网关
cd gateway
go run gateway.go -f etc/gateway-api.yaml
```

### 6. 验证

```bash
# 登录获取 Token
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 使用返回的 AccessToken 访问受保护接口
curl http://localhost:8888/api/v1/system/user \
  -H "Authorization: Bearer <access_token>"
```

## 统一响应格式

所有 HTTP 接口返回统一格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

- `code = 0` 表示成功
- `code != 0` 表示业务错误（通过 gRPC `ErrorInfo` 元数据传递业务错误码，网关统一转为 HTTP 200 响应）
- 认证失败返回 `401 Unauthorized`
- 权限不足返回 `403 Forbidden`

## 认证与授权

### JWT 双Token机制

- **AccessToken**：短期有效（默认 20 分钟），用于 API 请求认证
- **RefreshToken**：长期有效（默认 7 天），用于换取新的 AccessToken
- 退出登录时 AccessToken 加入 Redis 黑名单，带 TTL 自动过期

### Casbin RBAC

- 模型定义：`g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")`
- 用户通过角色获得权限：`用户 -> 角色 -> 权限`
- `keyMatch2` 支持 RESTful 路径通配，如 `/system/user/:id` 匹配 `/system/user/1`
- `admin` 角色拥有全部权限（硬编码跳过 Casbin 检查）
- 角色权限变更时自动重建 Casbin 策略

### 中间件链

```
HTTP Request
  -> [AuthMiddleware]  JWT 解析 + 黑名单校验 -> ctx.Set(userId, username)
  -> [CasbinMiddleware]  从 ctx 取 userId -> gRPC 调用 CheckPermission
  -> [Business Handler]  执行业务逻辑
```

## 作为新项目脚手架使用

该仓库设计为**架构仓**，可直接复用至新项目：

### 复用步骤

参考
**[Project Introduction.md](https://github.com/tianyuanxiang/go-zero-rpc/blob/dev/Project%20Introduction.md)**



### 模块依赖关系

```
common/           -- 被 gateway、sys.rpc 引用（通过 go.work）
service/sys/rpc/  -- 提供 gRPC 客户端给 gateway 和 job 引用
gateway/          -- HTTP 入口，依赖 sys.rpc 的 gRPC 客户端
service/job/      -- 定时任务，依赖 sys.rpc 的 gRPC 客户端
```

### 错误码体系

业务错误码统一在 `common/xerr/errcode.go` 中定义，通过 gRPC `status.ErrorInfo` 元数据透明传递到网关层，网关全局错误处理器将其转为 HTTP 200 + 业务错误码的响应格式，确保业务错误不触发 HTTP 5xx。

如需添加新错误码，在 `errcode.go` 中定义新的 `CodeError` 即可。

## 配置说明

### Gateway 配置参数

| 参数                | 说明                    | 默认值  |
| ------------------- | ----------------------- | ------- |
| `Name`              | 服务名称                | gateway |
| `Host`              | 监听地址                | 0.0.0.0 |
| `Port`              | 监听端口                | 8888    |
| `Mode`              | 运行模式 (dev/test/pro) | dev     |
| `Timeout`           | 请求超时 (ms)           | 30000   |
| `Auth.AccessSecret` | JWT 签名密钥            | -       |
| `Auth.AccessExpire` | AccessToken 有效期 (秒) | 1200    |
| `SysRpc.Etcd.Hosts` | etcd 集群地址           | -       |
| `SysRpc.Etcd.Key`   | etcd 服务注册 Key       | sys.rpc |

### SysRPC 配置参数

| 参数                    | 说明                        | 默认值              |
| ----------------------- | --------------------------- | ------------------- |
| `Name`                  | 服务名称                    | sys.rpc             |
| `ListenOn`              | 监听地址                    | 0.0.0.0:9100        |
| `JwtAuth.AccessSecret`  | JWT 签名密钥                | -                   |
| `JwtAuth.AccessExpire`  | AccessToken 有效期 (秒)     | 1200                |
| `JwtAuth.RefreshExpire` | RefreshToken 有效期 (分钟)  | 10080               |
| `DB.DataSource`         | PostgreSQL 连接字符串       | -                   |
| `CacheRedis`            | go-zero 内置缓存 Redis      | -                   |
| `BizRedis`              | 业务 Redis（Token黑名单等） | -                   |
| `CasbinModelPath`       | Casbin 模型文件路径         | etc/rbac_model.conf |
| `UploadPath`            | 文件上传目录                | uploads             |

## License

MIT License.
