# 槽液智能分析平台 — 业务接口说明文档

> 本文档说明 plating-platform 项目的完整目录架构、各层职责关系，以及所有 HTTP 业务接口的路由、参数和业务逻辑。



## 一、项目目录架构

```
plating-platform/
├── api/                    # go-zero API定义文件（接口契约）
│   ├── plating.api         # 主入口，引入所有子API定义
│   └── desc/               # API分类定义目录
│       ├── base.api        # 通用基础结构（分页请求/响应、ID请求等）
│       ├── auth.api        # 认证相关接口定义（登录/刷新/登出/修改密码）
│       ├── system/         # 系统管理接口定义
│       │   ├── user.api    # 用户管理（CRUD + 角色分配 + 密码重置）
│       │   ├── role.api    # 角色管理（CRUD + 菜单/API权限绑定）
│       │   ├── menu.api    # 菜单管理（CRUD + 树形查询）
│       │   ├── api.api     # API接口管理（CRUD + 分组查询）
│       │   ├── dict.api    # 字典管理（类型 + 数据的 CRUD）
│       │   ├── log.api     # 日志管理（登录日志 + 操作日志查询/清空）
│       │   └── file.api    # 文件上传管理（上传/删除/列表）
│       └── plating/        # 槽液业务接口定义
│           ├── tank.api    # 槽体配置管理（CRUD + 初始化模型状态）
│           ├── event.api   # 生产/加药/补水事件录入查询 + 触发计算
│           └── state.api   # 模型状态查询/趋势/化验校正/导出
├── docs/                   # 项目文档
│   ├── casbin_guide.md     # Casbin权限系统使用指南
│   └── business_api_interfaces.md  # 业务接口说明文档（本文档）
├── etc/                    # 配置文件目录
│   ├── plating.yaml        # 应用主配置文件（端口、数据库、JWT密钥等）
│   └── rbac_model.conf     # Casbin RBAC模型配置（策略规则定义）
├── internal/               # 私有业务代码（不对外暴露）
│   ├── config/             # 配置结构定义
│   │   └── config.go       # 配置项结构体（嵌入go-zero RestConf，扩展DB/JWT/Redis等）
│   ├── handler/            # HTTP请求处理层（薄层，只负责参数解析和响应）
│   │   ├── routes.go       # 路由注册总入口（由goctl自动生成，手动维护）
│   │   ├── auth/           # 认证相关Handler（login/logout/refresh等）
│   │   ├── system/         # 系统管理Handler（user/role/menu/api/dict/log/file）
│   │   └── plating/        # 槽液业务Handler（tank/event/state）
│   ├── logic/              # 业务逻辑层（核心层，所有业务规则在此实现）
│   │   ├── auth/           # 认证业务逻辑（JWT生成、密码验证、登录日志写入）
│   │   ├── system/         # 系统管理业务逻辑（RBAC管理、字典缓存等）
│   │   └── plating/        # 槽液业务逻辑
│   │       ├── tankconfiglogic.go      # 槽体配置CRUD逻辑 + 初始化模型状态
│   │       ├── modelstatelogic.go      # 模型状态查询/化验校正/导出逻辑
│   │       └── eventlogic.go           # 生产/加药/补水事件录入 + 触发计算逻辑
│   ├── middleware/         # HTTP中间件
│   │   ├── authmiddleware.go    # JWT认证中间件（解析令牌、注入用户信息到ctx）
│   │   └── casbinmiddleware.go  # Casbin权限校验中间件（校验路由+方法权限）
│   ├── model/              # 数据访问层（纯SQL，无业务逻辑）
│   │   ├── sysuser.go          # 系统用户模型（sys_user表）
│   │   ├── sysrole.go          # 系统角色模型（sys_role表）
│   │   ├── sysmenu.go          # 系统菜单模型（sys_menu表）
│   │   ├── sysapi.go           # 系统API接口模型（sys_api表）
│   │   ├── sysdicttype.go      # 系统字典类型模型（sys_dict_type表）
│   │   ├── sysdictdata.go      # 系统字典数据模型（sys_dict_data表）
│   │   ├── sysloginlog.go      # 登录日志模型（sys_login_log表）
│   │   ├── sysoperlog.go       # 操作日志模型（sys_oper_log表）
│   │   ├── sysfile.go          # 文件管理模型（sys_file表）
│   │   ├── sysroleapi.go       # 角色-API关联模型（sys_role_api表）
│   │   ├── sysrolemenu.go      # 角色-菜单关联模型（sys_role_menu表）
│   │   ├── sysuserrole.go      # 用户-角色关联模型（sys_user_role表）
│   │   ├── tankconfig.go       # 槽体配置模型（tank_config表）
│   │   ├── modelstate.go       # 机理模型状态模型（model_state表）
│   │   ├── productionevent.go  # 生产事件模型（production_event表）
│   │   ├── dosingevent.go      # 加药事件模型（dosing_event表）
│   │   └── waterevent.go       # 补水事件模型（water_event表）
│   ├── svc/                # 服务上下文（依赖注入中心）
│   │   └── servicecontext.go   # ServiceContext：初始化数据库连接、Model实例、JWT工具等并注入
│   └── types/              # 请求/响应结构体定义（DTO层）
│       └── types.go            # 所有Handler层使用的请求/响应DTO类型（由goctl生成）
├── pkg/                    # 可被外部包导入的公共工具库
│   ├── jwtx/               # JWT工具：生成AccessToken/RefreshToken、解析、刷新令牌
│   ├── xerr/               # 统一错误码定义与业务错误类型封装
│   └── utils/              # 通用工具函数（时间格式化、加密等）
├── schema/                 # 数据库初始化脚本
│   └── init.sql            # 完整建表SQL + 初始数据（角色/菜单/字典种子数据）
├── go.mod                  # Go模块定义（模块名、Go版本、依赖声明）
├── go.sum                  # 依赖校验文件（由go mod tidy自动维护）
└── main.go                 # 程序入口：加载配置、初始化ServiceContext、注册路由、启动HTTP服务
```

---

## 二、分层架构与调用链路

```
HTTP请求
   |
   v
[Handler层] internal/handler/
   - 使用 httpx.Parse() 解析请求参数
   - 调用对应 Logic 方法
   - 使用 httpx.OkJson() 返回统一响应
   - 不包含任何业务判断
   |
   v
[Logic层] internal/logic/
   - 包含所有业务规则和计算逻辑
   - 通过 ServiceContext 访问 Model 层
   - 处理参数校验、权限过滤、数据组装
   - 机理模型公式在此层实现
   |
   v
[Model层] internal/model/
   - 封装原生 SQL 语句
   - 通过 sqlx.SqlConn 执行数据库操作
   - 只做数据存取，不含业务规则
   |
   v
[数据库] MySQL
   - 存储所有业务数据
   - 通过 schema/init.sql 初始化表结构
```

**中间件调用位置**：在 Handler 层之前由 go-zero 框架自动调用。

- `AuthMiddleware`：解析 JWT，将 userId/roleCode 写入 context，令牌非法则直接返回 401。
- `CasbinMiddleware`：从 context 取出 roleCode，结合请求路径和 HTTP 方法查询 Casbin 策略，无权限则返回 403。

**ServiceContext**（`internal/svc/servicecontext.go`）是依赖注入中心，在服务启动时初始化：
- 数据库连接（sqlx.SqlConn）
- 所有 Model 实例（TankConfigModel、ModelStateModel 等）
- JWT 工具实例（jwtx）
- Casbin Enforcer 实例

---

## 三、认证接口（/api/auth）

### 3.1 用户登录

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/auth/login |
| 认证 | 无需认证 |

**请求体：**

```json
{
    "username": "admin",
    "password": "Admin@123"
}
```

**响应体：**

```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 7200,
    "userId": 1,
    "username": "admin",
    "nickname": "超级管理员"
}
```

---

### 3.2 刷新访问令牌

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/auth/refresh |
| 认证 | 无需认证（使用RefreshToken） |

**请求体：**

```json
{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应体：**

```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 7200
}
```

---

### 3.3 退出登录

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/auth/logout |
| 认证 | 需要 JWT |

**业务逻辑**：将当前 AccessToken 加入 Redis 黑名单，后续请求携带该 Token 将被拦截。

---

### 3.4 修改密码

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/auth/change-password |
| 认证 | 需要 JWT |

**请求体：**

```json
{
    "oldPassword": "Admin@123",
    "newPassword": "NewPass@456"
}
```

---

### 3.5 获取当前用户信息

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/auth/me |
| 认证 | 需要 JWT |

响应结构同登录响应（LoginResp）。

---

## 四、系统管理接口（/api/system）

> 所有系统管理接口均需要 JWT 认证 + Casbin 权限校验。

### 4.1 用户管理

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/user | 创建用户 |
| PUT | /api/system/user/:id | 更新用户信息 |
| DELETE | /api/system/user/:id | 删除用户 |
| GET | /api/system/user/:id | 查询用户详情 |
| GET | /api/system/user | 分页查询用户列表 |
| POST | /api/system/user/:id/reset-password | 重置用户密码 |

**创建用户请求体示例：**

```json
{
    "username": "operator01",
    "password": "Oper@123",
    "nickname": "操作员01",
    "email": "op01@example.com",
    "phone": "13800138001",
    "status": 1,
    "roleIds": [2, 3]
}
```

**用户列表查询参数：** `page`、`pageSize`、`keyword`（用户名/昵称模糊）、`status`（1启用/0禁用）

---

### 4.2 角色管理

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/role | 创建角色 |
| PUT | /api/system/role/:id | 更新角色 |
| DELETE | /api/system/role/:id | 删除角色 |
| GET | /api/system/role/:id | 查询角色详情 |
| GET | /api/system/role | 分页查询角色列表 |
| GET | /api/system/role/all | 获取所有角色（不分页，用于下拉选择） |

**创建角色请求体示例：**

```json
{
    "roleName": "操作员",
    "roleCode": "operator",
    "remark": "生产现场操作员，负责录入生产事件",
    "menuIds": [10, 11, 12],
    "apiIds": [30, 31, 32]
}
```

---

### 4.3 菜单管理

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/menu | 创建菜单 |
| PUT | /api/system/menu/:id | 更新菜单 |
| DELETE | /api/system/menu/:id | 删除菜单 |
| GET | /api/system/menu/tree | 获取完整菜单树（含所有子菜单） |
| GET | /api/system/menu/current | 获取当前登录用户有权限的菜单树 |

**menuType 枚举值：** 1=目录 2=菜单 3=按钮

---

### 4.4 API接口管理

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/api | 创建API记录 |
| PUT | /api/system/api/:id | 更新API记录 |
| DELETE | /api/system/api/:id | 删除API记录 |
| GET | /api/system/api | 分页查询API列表 |
| GET | /api/system/api/all | 获取所有API（不分页，用于角色权限配置） |

**列表查询参数：** `page`、`pageSize`、`keyword`（接口名/路径模糊）、`group`（分组过滤）

---

### 4.5 字典管理

**字典类型接口：**

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/dict/type | 创建字典类型 |
| PUT | /api/system/dict/type/:id | 更新字典类型 |
| DELETE | /api/system/dict/type/:id | 删除字典类型 |
| GET | /api/system/dict/type | 分页查询字典类型列表 |

**字典数据接口：**

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/dict/data | 创建字典数据 |
| PUT | /api/system/dict/data/:id | 更新字典数据 |
| DELETE | /api/system/dict/data/:id | 删除字典数据 |
| GET | /api/system/dict/data/:dictType | 根据字典类型编码查询字典数据（前端下拉框使用） |

**内置字典类型示例：**

| dictType | 说明 | 常用值 |
|----------|------|--------|
| plating_workpiece_type | 工件类型 | A=平板件, B=普通件, C=复杂件 |
| plating_drug_type | 药品类型 | CrO3=铬酸酐, H2SO4=硫酸 |
| sys_user_status | 用户状态 | 1=启用, 0=禁用 |

---

### 4.6 日志管理

| 方法 | 路由 | 说明 |
|------|------|------|
| GET | /api/system/log/login | 分页查询登录日志 |
| DELETE | /api/system/log/login/clear | 清空所有登录日志 |
| GET | /api/system/log/oper | 分页查询操作日志 |
| DELETE | /api/system/log/oper/clear | 清空所有操作日志 |

**登录日志查询参数：** `page`、`pageSize`、`username`、`status`（1成功/0失败）、`startTime`、`endTime`

**操作日志查询参数：** `page`、`pageSize`、`operName`、`operType`（create/update/delete/query）、`status`、`startTime`、`endTime`

---

### 4.7 文件管理

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | /api/system/file/upload | 上传文件（multipart/form-data） |
| DELETE | /api/system/file/:id | 删除文件 |
| GET | /api/system/file | 分页查询文件列表 |

**上传响应示例：**

```json
{
    "fileId": 42,
    "fileName": "report_20240115.xlsx",
    "fileUrl": "/uploads/2024/01/15/abc123.xlsx",
    "fileSize": 10240,
    "mimeType": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}
```

---

## 五、槽液业务接口 — 槽体配置（/api/plating/tank）

槽体配置是机理模型计算的基本单元，存储每个槽体的物理和化学参数。

> 所有槽液业务接口均需要 JWT 认证 + Casbin 权限校验。

### 5.1 创建槽体配置

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/tank |
| 说明 | 工程师首次部署时录入，一般不频繁变动 |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "tankName": "1号镀铬槽",
    "volumeLiters": 500.0,
    "densityGPerL": 1280.0,
    "dCrO3GPerAh": 0.4,
    "dCr3GPerAh": 0.08,
    "carryoverA": 1.0,
    "carryoverB": 2.0,
    "carryoverC": 6.5,
    "isActive": 1,
    "remark": "主生产线一号槽"
}
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tankId | string | 是 | 槽体唯一标识，自定义，不可修改，如 TANK-001 |
| volumeLiters | float64 | 是 | 槽液工作体积（升），从设备铭牌读取 |
| densityGPerL | float64 | 是 | 镀液密度（g/L），铬镀槽约1250~1350，默认1280 |
| dCrO3GPerAh | float64 | 是 | 每Ah消耗的CrO3（g），论文范围0.3~0.5，初始用0.4 |
| dCr3GPerAh | float64 | 是 | 每Ah生成的Cr3+（g），论文范围0.05~0.12，初始用0.08 |
| carryoverA | float64 | 是 | 平板件带液系数（g/dm²），Leiden论文Table4参考值：1.0 |
| carryoverB | float64 | 是 | 普通件带液系数（g/dm²），Leiden论文Table4参考值：2.0 |
| carryoverC | float64 | 是 | 复杂件带液系数（g/dm²），Leiden论文Table4参考值：6.5 |

---

### 5.2 更新槽体配置

| 项目 | 说明 |
|------|------|
| 路由 | PUT /api/plating/tank/:tankId |
| 说明 | 修改参数标定值（dCrO3GPerAh/dCr3GPerAh）时也使用此接口 |

---

### 5.3 删除槽体配置

| 项目 | 说明 |
|------|------|
| 路由 | DELETE /api/plating/tank/:tankId |
| 说明 | 物理删除，生产中的槽体请先停用（isActive=0）而非直接删除 |

---

### 5.4 查询槽体配置详情

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/tank/:tankId |

---

### 5.5 分页查询槽体配置列表

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/tank |
| 查询参数 | `page`、`pageSize`、`keyword`（槽体ID/名称模糊）、`isActive`（-1不过滤） |

---

### 5.6 初始化槽体模型状态（首次上线必用）

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/tank/init-state |
| 说明 | 系统上线时，化验员化验后录入真实浓度，作为机理模型计算的起点，source="init" |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "crO3GPerL": 248.0,
    "cr3GPerL": 4.2,
    "tankVolume": 500.0,
    "calcTime": "2024-01-15 09:00:00",
    "remark": "系统上线初始化化验数据"
}
```

**业务逻辑：**
1. 校验 tankId 对应的槽体配置存在且已激活
2. 若 tankVolume 为 0 或未填，自动从 tank_config.volume_liters 读取
3. 写入 model_state 表，source = "init"
4. 若已存在初始化状态记录，视为重新初始化（写入新记录，不覆盖历史）

---

## 六、槽液业务接口 — 事件管理（/api/plating/event）

事件是驱动机理模型计算的输入数据。三类事件（生产/加药/补水）录入后，调度器按时间顺序逐一处理，更新模型状态。

### 6.1 生产事件

#### 6.1.1 记录生产事件

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/event/production |
| 说明 | 电镀批次完成后由操作员录入，记录消耗的电量及工件信息 |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "startTime": "2024-01-15 09:00:00",
    "endTime": "2024-01-15 09:45:00",
    "ampHour": 150.0,
    "areaDm2": 100.0,
    "workpieceType": "B",
    "remark": "普通工件，整流器读数150Ah"
}
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ampHour | float64 | 是 | 本批次安时数（从整流器面板读取，或由TDEngine积分计算） |
| areaDm2 | float64 | 否 | 工件总表面积（dm²）= 单件面积 × 件数，用于带液量计算 |
| workpieceType | string | 否 | 工件类型：A=平板件（带液系数1.0），B=普通件（2.0），C=复杂件（6.5） |

**机理模型计算（生产事件触发）：**

```
公式C（电化学消耗CrO3）：CrO3_new = CrO3_old - (ampHour × dCrO3GPerAh) / volume
公式D（电化学生成Cr3+）：Cr3_new  = Cr3_old  + (ampHour × dCr3GPerAh)  / volume
公式B（工件带液带出）：
    带液质量 = areaDm2 × carryoverCoeff（由workpieceType选择A/B/C系数）
    带液体积 = 带液质量 / densityGPerL
    volume_new = volume_old - 带液体积
    浓度不变（带走的是原浓度液体）
```

#### 6.1.2 删除生产事件

| 项目 | 说明 |
|------|------|
| 路由 | DELETE /api/plating/event/production/:id |
| 说明 | 仅允许删除 is_processed=0（未处理）的事件，已处理事件不可删除 |

#### 6.1.3 查询生产事件列表

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/event/production |
| 查询参数 | `tankId`、`page`、`pageSize`、`startTime`、`endTime` |

---

### 6.2 加药事件

#### 6.2.1 记录加药事件

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/event/dosing |
| 说明 | 操作员每次向槽内加入化学品后立即录入，是提升模型准确性的关键操作 |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "eventTime": "2024-01-15 11:00:00",
    "drugType": "CrO3",
    "massGram": 2000.0,
    "remark": "日常补充铬酸酐"
}
```

**drugType 枚举值：** CrO3=铬酸酐，H2SO4=硫酸，Cr2O3=三氧化二铬

**机理模型计算（加药事件触发）：**

```
公式A（加药）：CrO3_new = CrO3_old + massGram / volume
              Cr3不变，体积不变
              （固体粉末溶解后体积变化极小，500L槽加500g误差<0.1%，忽略不计）
```

#### 6.2.2 删除加药事件

| 项目 | 说明 |
|------|------|
| 路由 | DELETE /api/plating/event/dosing/:id |
| 说明 | 仅允许删除未处理（is_processed=0）的事件 |

#### 6.2.3 查询加药事件列表

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/event/dosing |
| 查询参数 | `tankId`、`page`、`pageSize`、`startTime`、`endTime` |

---

### 6.3 补水事件

#### 6.3.1 记录补水事件

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/event/water |
| 说明 | 操作员补纯水后立即录入，或液位传感器自动触发（结合TDEngine液位数据可实现自动化录入） |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "eventTime": "2024-01-15 14:00:00",
    "volumeLiter": 5.0,
    "remark": "蒸发补水"
}
```

**机理模型计算（补水事件触发）：**

```
公式E（补水稀释）：
    volume_new   = volume_old + volumeLiter
    CrO3_new     = CrO3_old × volume_old / volume_new
    Cr3_new      = Cr3_old  × volume_old / volume_new
    （CrO3和Cr3+均等比例稀释，体积增加）
```

#### 6.3.2 删除补水事件

| 项目 | 说明 |
|------|------|
| 路由 | DELETE /api/plating/event/water/:id |
| 说明 | 仅允许删除未处理（is_processed=0）的事件 |

#### 6.3.3 查询补水事件列表

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/event/water |
| 查询参数 | `tankId`、`page`、`pageSize`、`startTime`、`endTime` |

---

### 6.4 手动触发模型计算

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/event/trigger-calc |
| 说明 | 手动触发未处理事件的计算，正常情况由后台调度器自动执行，调试或补录数据后使用 |

**请求体：**

```json
{
    "tankId": "TANK-001"
}
```

`tankId` 为空字符串时处理所有活跃槽体的未处理事件。

**响应体：**

```json
{
    "processedCount": 15
}
```

**事件处理顺序**（同一槽体多个未处理事件按时间升序逐一处理）：
1. 生产批次（production）：先执行公式C/D计算消耗和生成，再执行公式B计算带液带出
2. 加药（dosing）：执行公式A
3. 补水（water）：执行公式E

---

## 七、槽液业务接口 — 模型状态（/api/plating/state）

模型状态是机理模型的输出结果，是前端实时看板和历史趋势图的数据来源。

### 7.1 获取槽体最新模型状态

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/state/:tankId |
| 说明 | 返回该槽体 calc_time 最新的一条状态记录，用于实时看板展示 |

**响应体：**

```json
{
    "id": 1024,
    "tankId": "TANK-001",
    "calcTime": "2024-01-15 14:30:00",
    "crO3GPerL": 249.88,
    "cr3GPerL": 5.024,
    "tankVolume": 499.844,
    "source": "calc",
    "remark": "",
    "createdAt": "2024-01-15 14:30:05"
}
```

**source 字段含义：**

| 值 | 含义 |
|----|------|
| init | 系统上线初始化时录入的化验值 |
| calc | 机理模型推算的结果 |
| lab_override | 化验员手工校准覆盖的化验值 |

---

### 7.2 查询状态历史趋势（图表数据）

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/state/trend |
| 说明 | 返回时间序列数据，供前端 CrO3/Cr3+ 浓度折线图渲染，按 calc_time 升序排列 |

**查询参数：**

| 参数 | 必填 | 说明 |
|------|------|------|
| tankId | 是 | 槽体ID |
| startTime | 否 | 开始时间（格式：2006-01-02 15:04:05） |
| endTime | 否 | 结束时间 |
| limit | 否 | 最多返回数据点数量，默认200，最大建议500 |

**响应体：**

```json
{
    "tankId": "TANK-001",
    "list": [
        {
            "id": 1020,
            "tankId": "TANK-001",
            "calcTime": "2024-01-15 09:00:00",
            "crO3GPerL": 248.0,
            "cr3GPerL": 4.2,
            "tankVolume": 500.0,
            "source": "init"
        },
        {
            "id": 1021,
            "calcTime": "2024-01-15 09:45:00",
            "crO3GPerL": 247.88,
            "cr3GPerL": 4.224,
            "tankVolume": 499.844,
            "source": "calc"
        }
    ]
}
```

**前端使用方式**：提取 `calcTime` 作为 X 轴，`crO3GPerL` 和 `cr3GPerL` 分别作为两条折线数据。

---

### 7.3 化验校正（以实验室结果修正模型状态）

| 项目 | 说明 |
|------|------|
| 路由 | POST /api/plating/state/override |
| 说明 | 化验员出具化验结果后，将真实值写入模型状态，作为后续推算的新起点 |

**请求体：**

```json
{
    "tankId": "TANK-001",
    "crO3GPerL": 245.0,
    "cr3GPerL": 4.9,
    "calcTime": "2024-01-15 18:00:00",
    "remark": "日常化验校准，化验报告编号：LAB-20240115-001"
}
```

**业务说明：**
1. 写入 model_state 表，source = "lab_override"
2. 不对历史数据做回溯修正（取样到出结果可能间隔数小时，回溯重算成本高且意义有限）
3. 调度器以此化验值为新起点继续推算后续事件
4. calcTime 为空时使用当前时间作为校正时间点

---

### 7.4 导出状态报告

| 项目 | 说明 |
|------|------|
| 路由 | GET /api/plating/state/export |
| 说明 | 将指定时间范围的状态历史导出为 Excel 文件，响应头设置附件下载 |

**查询参数：** `tankId`（必填）、`startTime`、`endTime`

**响应**：Excel 文件流（Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet）

---

## 八、接口权限规划（Casbin角色分配建议）

| 角色 | roleCode | 说明 | 建议权限范围 |
|------|----------|------|------------|
| 超级管理员 | admin | 系统内置角色 | 全部接口（model.conf中跳过Casbin检查） |
| 管理员 | manager | 平台管理人员 | 系统管理全部 + 槽液分析全部 |
| 工程师 | engineer | 工艺工程师 | 槽体配置CRUD + 参数修改 + 全部查看 |
| 操作员 | operator | 生产现场操作员 | 录入生产/加药/补水事件 + 查看状态和趋势 |
| 只读 | viewer | 管理层查看人员 | 只能查看看板、趋势和历史数据 |
| 化验员 | lab | 化验室人员 | 化验校正接口 + 初始化状态 + 查看 |

**权限策略说明（Casbin格式）：**

```
p, engineer,  /api/plating/tank,        POST
p, engineer,  /api/plating/tank/*,      GET
p, engineer,  /api/plating/tank/*,      PUT
p, operator,  /api/plating/event/*, POST
p, operator,  /api/plating/state/*, GET
p, lab,       /api/plating/state/override, POST
p, lab,       /api/plating/tank/init-state, POST
p, viewer,    /api/plating/state/*,     GET
```

---

## 九、机理模型四大核心公式速查

根据 Leiden et al. 2020 论文及《机理模型技术文档.md》，每次处理事件时对应以下公式：

```
公式A（加药）：   CrO3_new = CrO3_old + M(g) / V(L)
                  Cr3不变，体积不变

公式B（带液带出）：V_new = V_old - (areaDm2 × coeff / density)
                  coeff 由 workpieceType 决定：A=1.0, B=2.0, C=6.5 g/dm²
                  浓度不变（带走同浓度液体）

公式C（电化学消耗CrO3）：CrO3_new = CrO3_old - (Ah × dCrO3GPerAh) / V

公式D（电化学生成Cr3+）：Cr3_new  = Cr3_old  + (Ah × dCr3GPerAh)  / V

公式E（补水）：   V_new    = V_old + V_water
                  CrO3_new = CrO3_old × V_old / V_new
                  Cr3_new  = Cr3_old  × V_old / V_new
```

**事件处理顺序**（同一槽体、同一时刻多类事件的处理优先级）：
1. 生产事件（production）：先执行公式C/D，再执行公式B
2. 加药事件（dosing）：执行公式A
3. 补水事件（water）：执行公式E

---

## 十、后台调度器说明（非HTTP接口）

调度器是后台 goroutine，不暴露 HTTP 接口，在 main.go 中随服务启动。

**实现位置**：建议新建 `internal/scheduler/scheduler.go`

**调度逻辑**：
1. 每5分钟轮询一次所有活跃槽体（is_active=1）
2. 对每个槽体，查询所有 is_processed=0 的事件（三种类型合并后按时间升序排列）
3. 从 model_state 取最新状态作为计算起点
4. 按顺序逐事件应用上述公式，每处理一个事件写入一条 model_state 记录（source="calc"）
5. 标记已处理事件为 is_processed=1

**手动触发**：通过 `POST /api/plating/event/trigger-calc` 可立即触发指定槽体的计算，不等待调度器轮询周期。

## 十一、Go-Zero生成说明

### 1、生成Api文件和接口的handler、logic

**相关接口需要单独放在一个文件夹下，需要加一个`group参数`**

```
@server (
	prefix: /api/system
	group: system/role
	middleware: AuthMiddleware,CasbinMiddleware
)
```

然后执行：

```bash
goctl api go --api .\api\plating.api --dir . --style go_zero
```

#### 生成的时候让 types 按模块分组

```bash
goctl api go --api .\api\plating.api --dir . --type-group --style go_zero
```

#### 生成结果：

```
etc目录下的配置文件
路由注册文件 handler/routes.go 和 所有的handler
所有的logic文件
types文件
```

### 2、生成model文件和基础CRUD	

如果表名已经分类：

```bash
goctl model mysql datasource -url="tianyx:bk147258.@tcp(172.16.90.70:3307)/plating" -t 'plate_*' -dir="./internal/model/plate" --style go_zero
```

如果表名没有分类 --- 在 -t 参数下列举相关表：

```bash
goctl model mysql datasource `
    --url='tianyx:bk147258.@tcp(172.16.90.70:3307)/plating?charset=utf8mb4&parseTime=true&loc=Local' `
    -t 'tank_config,dosing_event,production_event,water_event,state_snapshot' `
    -d './internal/model/plate' `
    -c `
    --style go_zero
```

### 3、改某个接口的前端请求参数或后端返回参数，流程

```bash
1. 在对应 .api 文件里修改 type XxxReq / type XxxResp
2. 如果是新增接口，就在api文件中对应 service 里新增路由定义
3. 重新执行 goctl api go --api .\api\plating.api --dir . --type-group --style go_zero
注意：每次修改任何 .api 文件后，始终用总入口文件生成：
```

### 4、格式化api文件

```bash
 goctl api format --dir .\api\desc\auth.api 
```



```
那也就是说：
1.这个四个模型都和AI没关系呗！前三个模型本质上就是根据获取的传感器数据来判断，第四个模型也是一个规则引擎，你推荐是使用决策树还是规则引擎来实现呢？
2.如果通过三个运行参数反推净化效果，按照正常的业务流程，前三个模型的输入数据是同步的吗？如果测点的采集频率不一致，那么前三个模型的结果不能同步，那第4个引擎什么时候执行呢？
3.你说多参数同时异常时需要有主次排序逻辑（比如pH异常的影响权重可能高于填料老化），这个排序规则是什么？
4.现在经常听到大家说PLC，PLC 到底是啥啊？
5.你仔细分析文档，告诉我到底需要处理哪些塔呀？分别需要处理什么指标？分别需要采集什么指标？多个塔如何处理？
6.你把需要找客户确认的关键问题再细化一下告诉我。
```

