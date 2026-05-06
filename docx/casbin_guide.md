# Casbin 权限管理完全指南

> 本文档结合本项目（plating）的实际代码，从零讲�?Casbin 的核心概念和完整使用方式�?> 读完此文档，你将彻底掌握 Casbin �?RBAC 权限模型，以及它�?go-zero 项目中的集成方法�?
---

## 一、Casbin 是什么？

Casbin 是一个开源的 Go 权限管理框架，核心职责只有一件事�?
> **回答一个问题：用户 X 是否有权限对资源 Y 执行操作 Z�?*

它不管理用户账号和密码，那是认证（Authentication）的事�?Casbin 只管授权（Authorization）：**谁能做什�?*�?
### 1.1 认证 vs 授权

先搞清两个概念，否则后面会混淆：

| 概念 | 英文 | 职责 | 本项目中谁负�?|
|------|------|------|---------------|
| 认证 | Authentication | 你是谁？（验证身份） | `AuthMiddleware`（JWT�?|
| 授权 | Authorization | 你能做什么？（验证权限） | `CasbinMiddleware`（Casbin�?|

**先认证，后授�?*。请求进来先�?AuthMiddleware 确认"你是�?，再�?CasbinMiddleware 确认"你能不能做这件事"�?
### 1.2 和传统硬编码鉴权的对�?
```go
// 传统做法：每个接口里硬编码权限检查，改一个接口就要改一次代�?func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    role := getRole(r)
    if role != "admin" && role != "manager" {
        http.Error(w, "没有权限", 403)
        return
    }
    // 业务逻辑...
}

// Casbin做法：统一中间件拦截，业务代码完全不用管权�?// 只需在数据库里维�?哪个角色能访问哪个接�?的规则即�?// 新增接口？往数据库加一行规则就行，不用改代�?```

---

## 二、核心概念：PERM 模型

Casbin 基于 **PERM 模型**（Policy, Effect, Request, Matchers），四个字母分别对应四个概念�?
| 组件 | 含义 | 类比 |
|------|------|------|
| **R**equest（请求） | 待鉴权的三元组：谁、什么资源、什么操�?| 有人刷卡想进�?|
| **P**olicy（策略） | 存储在数据库中的权限规则 | 门禁白名�?|
| **M**atchers（匹配器�?| 如何�?Request �?Policy 进行比较 | 门禁比对规则 |
| **E**ffect（效果） | 多条策略命中时的最终裁�?| 最终放行还是拦�?|

**通俗理解整个流程**�?
```
有人刷卡想进A栋（Request�?  �?门禁系统拿着刷卡信息去比对白名单（Policy�?  �?按照预设的比对规则逐条检查（Matchers�?  �?只要白名单里有一条匹配就放行（Effect�?```

---

## 三、本项目�?model.conf 详解

文件位置：`etc/rbac_model.conf`

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

### 逐行解读

#### 3.1 `[request_definition]` -- 请求长什么样

```ini
r = sub, obj, act
```

定义了每次鉴权请求包含三个要素：

| 字段 | 含义 | 本项目中的实际�?| 示例 |
|------|------|-----------------|------|
| `sub` | 主体（Subject），谁在请求 | 角色编码（role.Code�?| `"admin"`, `"operator"` |
| `obj` | 对象（Object），访问什�?| HTTP 请求路径 | `"/api/system/user"` |
| `act` | 操作（Action），做什么动�?| HTTP 方法 | `"GET"`, `"POST"`, `"DELETE"` |

**注意**：`sub` 不是用户ID，是**角色编码**。一个用户可以有多个角色，中间件会逐个角色去检查�?
#### 3.2 `[policy_definition]` -- 规则长什么样

```ini
p = sub, obj, act
```

数据�?`casbin_rule` 表中每一行存储的就是一条策略，格式和请求一样是三个字段�?
**对应到数据库�?*�?
| 表字�?| 对应 | 含义 | 示例�?|
|--------|------|------|--------|
| `ptype` | 策略类型 | `p` 表示普通策�?| `"p"` |
| `v0` | sub（角色编码） | 哪个角色 | `"admin"` |
| `v1` | obj（路径） | 能访问什么路�?| `"/api/system/user"` |
| `v2` | act（方法） | 用什么HTTP方法 | `"GET"` |

#### 3.3 `[role_definition]` -- 角色继承

```ini
g = _, _
```

`g` 定义角色之间的继承关系。`_, _` 表示支持两个参数：子角色和父角色�?
**本项目当前没有使用角色继�?*。如果未来需�?经理继承操作员的所有权�?，可以在 `casbin_rule` 表中添加�?
| ptype | v0 | v1 |
|-------|----|----|
| g | manager | operator |

这样 manager 自动拥有 operator 的所有权限，不用重复配置�?
#### 3.4 `[policy_effect]` -- 最终裁决规�?
```ini
e = some(where (p.eft == allow))
```

含义�?*只要存在任意一条匹配的策略，就允许访问**�?这是"白名�?模式 -- 默认拒绝，只有明确授权的才放行�?
#### 3.5 `[matchers]` -- 匹配规则（最核心�?
```ini
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

这一行决定了"请求"�?规则"怎么比对。拆开看：

| 表达�?| 含义 | 解释 |
|--------|------|------|
| `g(r.sub, p.sub)` | 角色匹配 | 请求者的角色 == 策略中的角色（含继承关系�?|
| `keyMatch2(r.obj, p.obj)` | **路径通配匹配** | 请求路径能匹配策略中的路径模�?|
| `r.act == p.act` | 方法精确匹配 | HTTP方法完全一�?|
| `p.act == "*"` | 方法通配 | 策略中方法为`*`时匹配任意方�?|

**三个条件�?`&&` 连接，必须同时满足才算匹�?*�?
##### keyMatch2 路径匹配详解

`keyMatch2` �?Casbin 内置的路径匹配函数，支持 `:param` 风格的参数通配�?
| 策略中的路径（p.obj�?| 请求路径（r.obj�?| 是否匹配 | 说明 |
|----------------------|-------------------|---------|------|
| `/api/system/user` | `/api/system/user` | 匹配 | 精确匹配 |
| `/api/system/user/:id` | `/api/system/user/123` | 匹配 | `:id` 通配任意�?|
| `/api/system/user/:id` | `/api/system/user/456/detail` | 不匹�?| 多了一层路�?|
| `/api/*` | `/api/system/user` | 匹配 | `*` 通配所有子路径 |
| `/api/*` | `/api/plating/event/dosing` | 匹配 | `*` 通配所有子路径 |

##### 方法通配详解

`(r.act == p.act || p.act == "*")` 的含义：

| 策略中的方法（p.act�?| 请求方法（r.act�?| 是否匹配 |
|---------------------|------------------|---------|
| `GET` | `GET` | 匹配 |
| `GET` | `POST` | 不匹�?|
| `*` | `GET` | 匹配 |
| `*` | `POST` | 匹配 |
| `*` | 任意方法 | 都匹�?|

---

## 四、本项目的数据库表结�?
### 4.1 casbin_rule 表（Casbin 自动管理�?
文件位置：`schema/plating.sql`

```sql
CREATE TABLE `casbin_rule` (
  `id`    bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ptype` varchar(100) NULL DEFAULT NULL,  -- 策略类型: p=普通策�? g=角色继承
  `v0`    varchar(100) NULL DEFAULT NULL,  -- 对应 sub（角色编码）
  `v1`    varchar(100) NULL DEFAULT NULL,  -- 对应 obj（路径）
  `v2`    varchar(100) NULL DEFAULT NULL,  -- 对应 act（方法）
  `v3`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用�?  `v4`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用�?  `v5`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用�?  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
) COMMENT = 'Casbin权限规则�?;
```

**初始数据**（超级管理员拥有所有权限）�?
```sql
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');
```

这一条规则的含义：`admin` 角色可以�?*任意HTTP方法**访问 `/api/` 下的**所有路�?*�?
### 4.2 相关联的业务�?
Casbin 只管 `casbin_rule` 这一张表，但鉴权流程中还会查询以下表�?
| 表名 | 作用 | 关键字段 |
|------|------|---------|
| `sys_user` | 用户�?| id, username, status |
| `sys_role` | 角色�?| id, code, name, status |
| `sys_user_role` | 用户-角色关联�?| user_id, role_id |
| `sys_api` | API接口注册�?| path, method, group, description |
| `sys_role_api` | 角色-API关联�?| role_id, api_id |

**这些表之间的关系**�?
```
sys_user  ←──  sys_user_role  ──�? sys_role
                                      �?                                      �?                                sys_role_api  ──�? sys_api
                                      �?                                      �?                                casbin_rule（Casbin自动读写�?```

通俗来说�?- `sys_user_role` 记录"用户有哪些角�?
- `sys_role_api` 记录"角色能访问哪些接�?（业务层面的关联关系�?- `casbin_rule` 记录"角色能访问哪些接�?（Casbin 层面的策略，�?sys_role_api 是同步的�?
为什�?`sys_role_api` �?`casbin_rule` 看起来重复？因为 `sys_role_api` 是给前端页面展示用的�?这个角色勾选了哪些接口"），`casbin_rule` 是给 Casbin 引擎做鉴权用的。两者通过 "为角色分配API" 的接口保持同步�?
---

## 五、完整鉴权链条：从请求到放行

下面用一个真实的例子，把整条链路串起来�?
### 场景：操作员张三要录入一条加药事�?
```
POST /api/plating/event/dosing
Authorization: Bearer eyJhbGci...（张三的JWT令牌�?Content-Type: application/json

{"tankId": "T001", "eventTime": "2026-04-15 10:30:00", ...}
```

### 第一步：路由匹配

go-zero 根据 `routes.go` 找到这个路由，它属于带有 `[AuthMiddleware, CasbinMiddleware]` 的路由组�?
```go
// routes.go �?6-123�?server.AddRoutes(
    rest.WithMiddlewares(
        []rest.Middleware{serverCtx.AuthMiddleware, serverCtx.CasbinMiddleware},
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/event/dosing",
                Handler: plateevent.CreateDosingEventHandler(serverCtx),
            },
            // ...
        }...,
    ),
    rest.WithPrefix("/api/plating"),
)
```

**中间件按顺序执行**：先 AuthMiddleware，再 CasbinMiddleware，最后才�?Handler�?
### 第二步：AuthMiddleware -- 确认"你是�?

文件位置：`internal/middleware/auth_middleware.go`

```
1. 从请求头提取 Authorization: Bearer eyJhbGci...
2. 解析JWT令牌，提�?userId=5, username="zhangsan"
3. �?userId �?username 写入请求�?Context �?4. 放行，交给下一个中间件
```

如果JWT无效或过期，直接返回 401，不会走�?CasbinMiddleware�?
### 第三步：CasbinMiddleware -- 确认"你能不能�?

文件位置：`internal/middleware/casbin_middleware.go`

```
1. �?Context 读取 userId=5
   （调�?GetUserIdFromCtx(r.Context())�?
2. �?sys_user_role 表：userId=5 的角色有哪些�?   �?得到 roleIds = [2]（张三只有一个角色）

3. �?sys_role 表：roleId=2 的角色编码是什么？
   �?得到 role.Code = "operator"

4. 调用 Casbin 鉴权�?   enforcer.Enforce("operator", "/api/plating/event/dosing", "POST")

5. Casbin 在内存中匹配 casbin_rule 表的规则�?   找到一条：("p", "operator", "/api/plating/event/dosing", "POST")
   
   �?matchers 比对�?   - g("operator", "operator") �?角色匹配
   - keyMatch2("/api/plating/event/dosing", "/api/plating/event/dosing") �?路径匹配
   - "POST" == "POST" �?方法匹配
   
   三个条件都满�?�?返回 true

6. hasPermission = true �?放行，交�?Handler 处理业务逻辑
```

**如果张三没有这个权限呢？**

```
Casbin 匹配不到任何规则 �?返回 false
�?中间件返�?403 Forbidden，请求到此结束，不会进入 Handler
```

### 流程图总结

```bash
客户端请�?    �?    �?路由匹配（routes.go�?    �?    �?AuthMiddleware（JWT认证�?    �?失败 �?401 Unauthorized
    �?成功
CasbinMiddleware（权限校验）
    �?    ├─ 1. 从Context取userId
    ├─ 2. 查sys_user_role �?得到roleIds
    ├─ 3. 查sys_role �?得到role.Code
    ├─ 4. enforcer.Enforce(role.Code, path, method)
    �?    └─ 在内存中匹配casbin_rule的规�?    �?    �?无权�?�?403 Forbidden
    �?有权�?Handler（业务逻辑�?    �?    �?返回响应
```

---

## 六、哪些路由受 Casbin 保护�?
�?`routes.go` 可以看出，路由分为三个安全等级：

### 6.1 完全公开（无任何中间件）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/auth/login` | POST | 登录 |
| `/api/auth/refresh` | POST | 刷新令牌 |

### 6.2 仅需登录（只�?AuthMiddleware，无 Casbin�?
| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/auth/changePassword` | POST | 修改自己的密�?|
| `/api/auth/logout` | POST | 登出 |
| `/api/auth/userInfo` | GET | 获取当前用户信息 |

**这些接口只要登录了就能访问，不检查角色权�?*。因为它们是每个用户都应该有的基础操作�?
### 6.3 需要登�?+ 角色授权（AuthMiddleware + CasbinMiddleware�?
以下所有接口都�?Casbin 保护，必须在 `casbin_rule` 表中有对应规则才能访问：

**系统管理 -- 用户**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/user` | POST | 创建用户 |
| `/api/system/user` | GET | 用户列表 |
| `/api/system/user/:id` | PUT | 更新用户 |
| `/api/system/user/:id` | DELETE | 删除用户 |
| `/api/system/user/:id` | GET | 用户详情 |
| `/api/system/user/:id/reset-password` | POST | 重置用户密码 |

**系统管理 -- 角色**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/role` | POST | 创建角色 |
| `/api/system/role` | GET | 角色列表 |
| `/api/system/role/:id` | PUT | 更新角色 |
| `/api/system/role/:id` | DELETE | 删除角色 |
| `/api/system/role/:id` | GET | 角色详情 |
| `/api/system/role/all` | GET | 所有角色（不分页） |

**系统管理 -- 菜单**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/menu` | POST | 创建菜单 |
| `/api/system/menu/:id` | PUT | 更新菜单 |
| `/api/system/menu/:id` | DELETE | 删除菜单 |
| `/api/system/menu/current` | GET | 当前用户的菜�?|
| `/api/system/menu/tree` | GET | 菜单�?|

**系统管理 -- 接口**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/api` | POST | 注册接口 |
| `/api/system/api` | GET | 接口列表 |
| `/api/system/api/:id` | PUT | 更新接口 |
| `/api/system/api/:id` | DELETE | 删除接口 |
| `/api/system/api/all` | GET | 所有接口（不分页） |

**系统管理 -- 字典**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/dict/type` | POST | 创建字典类型 |
| `/api/system/dict/type` | GET | 字典类型列表 |
| `/api/system/dict/type/:id` | PUT | 更新字典类型 |
| `/api/system/dict/type/:id` | DELETE | 删除字典类型 |
| `/api/system/dict/data` | POST | 创建字典数据 |
| `/api/system/dict/data/:dictType` | GET | 按类型查字典数据 |
| `/api/system/dict/data/:id` | PUT | 更新字典数据 |
| `/api/system/dict/data/:id` | DELETE | 删除字典数据 |

**系统管理 -- 文件**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/file` | GET | 文件列表 |
| `/api/system/file/:id` | DELETE | 删除文件 |
| `/api/system/file/upload` | POST | 上传文件 |

**系统管理 -- 日志**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/log/login` | GET | 登录日志列表 |
| `/api/system/log/login/clear` | DELETE | 清空登录日志 |
| `/api/system/log/oper` | GET | 操作日志列表 |
| `/api/system/log/oper/clear` | DELETE | 清空操作日志 |

**业务功能 -- 槽液事件**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/event/dosing` | POST | 录入加药事件 |
| `/api/plating/event/dosing` | GET | 加药事件列表 |
| `/api/plating/event/dosing/:id` | DELETE | 删除加药事件 |
| `/api/plating/event/production` | POST | 录入生产事件 |
| `/api/plating/event/production` | GET | 生产事件列表 |
| `/api/plating/event/production/:id` | DELETE | 删除生产事件 |
| `/api/plating/event/water` | POST | 录入换水事件 |
| `/api/plating/event/water` | GET | 换水事件列表 |
| `/api/plating/event/water/:id` | DELETE | 删除换水事件 |
| `/api/plating/event/trigger-calc` | POST | 触发计算 |

**业务功能 -- 槽体状�?*

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/state/:tankId` | GET | 最新状�?|
| `/api/plating/state/export` | GET | 导出报表 |
| `/api/plating/state/override` | POST | 覆盖状�?|
| `/api/plating/state/trend` | GET | 趋势数据 |

**业务功能 -- 槽体配置**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/tank` | POST | 创建槽体 |
| `/api/plating/tank` | GET | 槽体列表 |
| `/api/plating/tank/:tankId` | PUT | 更新槽体 |
| `/api/plating/tank/:tankId` | DELETE | 删除槽体 |
| `/api/plating/tank/:tankId` | GET | 槽体详情 |
| `/api/plating/tank/init-state` | POST | 初始化模型状�?|

---

## 七、超级管理员是怎么拥有所有权限的�?
初始化SQL中只�?admin 角色写了一条规则：

```sql
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');
```

翻译一下：

| 字段 | �?| 含义 |
|------|-----|------|
| ptype | `p` | 这是一条策�?|
| v0 (sub) | `admin` | admin 角色 |
| v1 (obj) | `/api/*` | 所�?`/api/` 下的路径 |
| v2 (act) | `*` | 所有HTTP方法 |

�?admin 角色访问任何接口时：

```
enforcer.Enforce("admin", "/api/system/user/5/reset-password", "POST")

匹配过程�?1. g("admin", "admin") �?角色匹配
2. keyMatch2("/api/system/user/5/reset-password", "/api/*") �?路径匹配�?通配所有子路径�?3. "POST" == "*" �?不等... �?p.act == "*" �?方法通配匹配

三个条件都满�?�?返回 true �?放行
```

**所�?admin 不需要一条条配置，一条通配规则搞定一切�?*

---

## 八、本项目�?Casbin 初始化代�?
### 8.1 Casbin 引擎初始�?
文件位置：`pkg/casbin/casbin.go`

```go
func NewCasbin(db *gorm.DB, modelPath string) (*casbinv2.Enforcer, error) {
    // 1. 使用 gorm-adapter 连接 MySQL
    //    Casbin 规则自动读写 casbin_rule �?    //    如果表不存在，gorm-adapter 会自动建�?    adapter, err := gormadapter.NewAdapterByDB(db)
    
    // 2. 加载模型配置文件（etc/rbac_model.conf�?    //    定义了请求格式、策略格式、匹配规�?    enforcer, err := casbinv2.NewEnforcer(modelPath, adapter)
    
    // 3. 从数据库加载所有策略到内存
    //    鉴权时直接在内存中匹配，不查数据库，性能极高
    enforcer.LoadPolicy()
    
    return enforcer, nil
}
```

**关键理解**：Casbin 把数据库中的规则**一次性加载到内存**，后续所�?`Enforce()` 调用都在内存中比对，不会频繁查数据库。这就是它性能高的原因。但代价是：**修改了数据库中的规则后，必须调用 `LoadPolicy()` 重新加载，否则内存中还是旧规则�?*

### 8.2 注入�?ServiceContext

文件位置：`internal/svc/service_context.go`

```go
type ServiceContext struct {
    // ... 其他字段
    Enforcer         *casbinv2.Enforcer                        // Casbin执行�?    AuthMiddleware   func(handlerFunc http.HandlerFunc) http.HandlerFunc
    CasbinMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc
    // ...
}

// 初始化时
enforcer, err := casbin.NewCasbin(db, c.CasbinModelPath)
// ...
return &ServiceContext{
    Enforcer:         enforcer,
    AuthMiddleware:   middleware.AuthMiddleware(c),
    CasbinMiddleware: middleware.CasbinMiddleware(enforcer, conn, c.CacheRedis, db),
    // ...
}
```

**配置文件**（`etc/plating-api.yaml`）：

```yaml
CasbinModelPath: "etc/rbac_model.conf"
```

### 8.3 鉴权中间件的完整代码

文件位置：`internal/middleware/casbin_middleware.go`

```go
func CasbinMiddleware(enforcer *casbinv2.Enforcer, conn sqlx.SqlConn,
    c cache.CacheConf, db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
    
    // 初始化需要查询的Model
    userRoleModel := systemmodel.NewSysUserRoleModel(conn, c, db)
    roleModel := systemmodel.NewSysRoleModel(conn, c, db)

    return func(next http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 第一步：从Context获取当前用户ID（AuthMiddleware写入的）
            userId := GetUserIdFromCtx(r.Context())
            if userId == 0 {
                response.FailUnauthorized(w, r)
                return
            }

            reqPath := r.URL.Path   // 如：/api/system/user
            reqMethod := r.Method   // 如：GET

            // 第二步：查询该用户的所有角色ID
            roleIds, err := userRoleModel.GetRoleIdsByUserId(r.Context(), userId)
            // ...（错误处理）

            if len(roleIds) == 0 {
                // 用户没有任何角色，直接拒�?                response.FailForbidden(w, r)
                return
            }

            // 第三步：逐个角色检查权限（OR逻辑：任意一个角色有权限就放行）
            hasPermission := false
            for _, roleId := range roleIds {
                role, err := roleModel.FindOneByRoleId(r.Context(), roleId)
                if err != nil || role == nil {
                    continue
                }

                // 用角色编码（code）调用Casbin鉴权
                allowed, err := enforcer.Enforce(role.Code, reqPath, reqMethod)
                if err != nil {
                    continue
                }

                if allowed {
                    hasPermission = true
                    break  // 有一个角色有权限就够了，不用继续检�?                }
            }

            if !hasPermission {
                response.FailForbidden(w, r)
                return
            }

            // 有权限，放行到下一个Handler
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 九、casbin_rule 表数据配置实�?
### 9.1 场景设定

一个电镀管理系统，有四种角色�?
| id | code | name | 定位 |
|----|------|------|------|
| 1 | admin | 超级管理�?| 什么都能做，不受限�?|
| 2 | manager | 管理�?| 能管人、管角色、管接口，但不碰业务数据 |
| 3 | operator | 操作�?| 能录入业务数据、查看状态，但不能管系统 |
| 4 | viewer | 查看�?| 只能看，什么都不能�?|

初始状态下 `casbin_rule` 表只有一�?admin 的通配规则�?
```sql
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');
```

下面按角色逐个配置�?
### 9.2 �?管理�?配置系统管理权限

管理员的职责：管理用户、管理角色、管理菜单、管理接口注册、查看日志。不参与业务数据操作�?
```sql
-- ========== 用户管理（增删改查） ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'manager', '/api/system/user',                    'POST'),   -- 创建用户
('p', 'manager', '/api/system/user',                    'GET'),    -- 用户列表
('p', 'manager', '/api/system/user/:id',                'GET'),    -- 用户详情
('p', 'manager', '/api/system/user/:id',                'PUT'),    -- 更新用户
('p', 'manager', '/api/system/user/:id',                'DELETE'), -- 删除用户
('p', 'manager', '/api/system/user/:id/reset-password', 'POST');   -- 重置用户密码

-- ========== 角色管理（增删改查） ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'manager', '/api/system/role',     'POST'),   -- 创建角色
('p', 'manager', '/api/system/role',     'GET'),    -- 角色列表
('p', 'manager', '/api/system/role/:id', 'GET'),    -- 角色详情
('p', 'manager', '/api/system/role/:id', 'PUT'),    -- 更新角色
('p', 'manager', '/api/system/role/:id', 'DELETE'), -- 删除角色
('p', 'manager', '/api/system/role/all', 'GET');    -- 所有角�?不分�?

-- ========== 菜单管理 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'manager', '/api/system/menu',         'POST'),   -- 创建菜单
('p', 'manager', '/api/system/menu/:id',     'PUT'),    -- 更新菜单
('p', 'manager', '/api/system/menu/:id',     'DELETE'), -- 删除菜单
('p', 'manager', '/api/system/menu/current', 'GET'),    -- 当前用户菜单
('p', 'manager', '/api/system/menu/tree',    'GET');    -- 菜单�?
-- ========== 接口管理 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'manager', '/api/system/api',     'POST'),   -- 注册接口
('p', 'manager', '/api/system/api',     'GET'),    -- 接口列表
('p', 'manager', '/api/system/api/:id', 'PUT'),    -- 更新接口
('p', 'manager', '/api/system/api/:id', 'DELETE'), -- 删除接口
('p', 'manager', '/api/system/api/all', 'GET');    -- 所有接�?不分�?

-- ========== 日志查看（只看不清） ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'manager', '/api/system/log/login', 'GET'),  -- 查看登录日志
('p', 'manager', '/api/system/log/oper',  'GET');   -- 查看操作日志
```

**注意**：管理员没有 `/api/system/log/login/clear` �?`/api/system/log/oper/clear` �?DELETE 权限，所以管理员只能查看日志，不能清空日志。清空日志只�?admin 能做�?
**验证一�?*：管理员尝试删除操作日志会怎样�?
```
请求: manager DELETE /api/system/log/oper/clear

Casbin匹配过程:
1. 遍历manager的所有p规则
2. 没有任何一条规则的 v1 能匹�?/api/system/log/oper/clear �?v2 == DELETE
3. 无匹�?�?返回 false �?403 Forbidden
```

### 9.3 �?操作�?配置业务操作权限

操作员的职责：录入加�?换水/生产事件，查看槽体状态和趋势，触发计算。不能管理系统�?
```sql
-- ========== 加药事件（录�?+ 查看 + 删除�?==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/event/dosing',     'POST'),   -- 录入加药
('p', 'operator', '/api/plating/event/dosing',     'GET'),    -- 查看加药列表
('p', 'operator', '/api/plating/event/dosing/:id', 'DELETE'); -- 删除加药记录

-- ========== 换水事件 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/event/water',     'POST'),
('p', 'operator', '/api/plating/event/water',     'GET'),
('p', 'operator', '/api/plating/event/water/:id', 'DELETE');

-- ========== 生产事件 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/event/production',     'POST'),
('p', 'operator', '/api/plating/event/production',     'GET'),
('p', 'operator', '/api/plating/event/production/:id', 'DELETE');

-- ========== 触发计算 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/event/trigger-calc', 'POST');

-- ========== 查看槽体状态（只读�?==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/state/:tankId', 'GET'),  -- 某槽体最新状�?('p', 'operator', '/api/plating/state/trend',   'GET'),  -- 趋势数据
('p', 'operator', '/api/plating/state/export',  'GET');   -- 导出报表

-- ========== 查看槽体配置（只读，不能增删改） ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/tank',          'GET'),  -- 槽体列表
('p', 'operator', '/api/plating/tank/:tankId',  'GET');   -- 槽体详情

-- ========== 基础功能：查看自己的菜单 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/system/menu/current', 'GET');
```

**验证一�?*：操作员尝试创建新用户会怎样�?
```
请求: operator POST /api/system/user

Casbin匹配过程:
1. 遍历operator的所有p规则
2. 没有任何一条规则的 v1 == /api/system/user �?v2 == POST
3. 无匹�?�?返回 false �?403 Forbidden

结论：操作员无法访问系统管理功能，被成功拦截�?```

**验证一�?*：操作员尝试删除槽体会怎样�?
```
请求: operator DELETE /api/plating/tank/T001

Casbin匹配过程:
1. 找到 (operator, /api/plating/tank/:tankId, GET)
2. keyMatch2("/api/plating/tank/T001", "/api/plating/tank/:tankId") �?路径匹配
3. �?"DELETE" != "GET" �?"GET" != "*" �?方法不匹�?4. 无完整匹�?�?返回 false �?403 Forbidden

结论：操作员能查看槽体详情，但不能删除槽体。路径匹配了，方法没匹配�?```

### 9.4 �?查看�?配置只读权限

查看者的职责：只能看数据，一�?POST/PUT/DELETE 权限都没有�?
```sql
-- ========== 槽体状态（只读�?==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'viewer', '/api/plating/state/:tankId', 'GET'),
('p', 'viewer', '/api/plating/state/trend',   'GET'),
('p', 'viewer', '/api/plating/state/export',  'GET');

-- ========== 事件记录（只读，只能看不能录入和删除�?==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'viewer', '/api/plating/event/dosing',     'GET'),
('p', 'viewer', '/api/plating/event/water',      'GET'),
('p', 'viewer', '/api/plating/event/production', 'GET');

-- ========== 槽体配置（只读） ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'viewer', '/api/plating/tank',         'GET'),
('p', 'viewer', '/api/plating/tank/:tankId', 'GET');

-- ========== 基础功能 ==========
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'viewer', '/api/system/menu/current', 'GET');
```

**验证一�?*：查看者尝试录入加药事件会怎样�?
```
请求: viewer POST /api/plating/event/dosing

Casbin匹配过程:
1. 找到 (viewer, /api/plating/event/dosing, GET)
2. 路径匹配，但 "POST" != "GET" �?方法不匹�?3. 无完整匹�?�?返回 false �?403 Forbidden

结论：查看者只有GET权限，所有写操作都被拦截�?```

### 9.5 完整�?casbin_rule 表总览

| id | ptype | v0 | v1 | v2 | 说明 |
|----|-------|---------|----------------------------------------------|--------|------|
| 1 | p | admin | /api/* | * | 超管通配，拥有一切权�?|
| | | | | | |
| | | | **管理�?-- 系统管理** | | |
| 2 | p | manager | /api/system/user | POST | 创建用户 |
| 3 | p | manager | /api/system/user | GET | 用户列表 |
| 4 | p | manager | /api/system/user/:id | GET | 用户详情 |
| 5 | p | manager | /api/system/user/:id | PUT | 更新用户 |
| 6 | p | manager | /api/system/user/:id | DELETE | 删除用户 |
| 7 | p | manager | /api/system/user/:id/reset-password | POST | 重置密码 |
| 8 | p | manager | /api/system/role | POST | 创建角色 |
| 9 | p | manager | /api/system/role | GET | 角色列表 |
| 10 | p | manager | /api/system/role/:id | GET | 角色详情 |
| 11 | p | manager | /api/system/role/:id | PUT | 更新角色 |
| 12 | p | manager | /api/system/role/:id | DELETE | 删除角色 |
| 13 | p | manager | /api/system/role/all | GET | 所有角�?|
| 14 | p | manager | /api/system/menu | POST | 创建菜单 |
| 15 | p | manager | /api/system/menu/:id | PUT | 更新菜单 |
| 16 | p | manager | /api/system/menu/:id | DELETE | 删除菜单 |
| 17 | p | manager | /api/system/menu/current | GET | 当前菜单 |
| 18 | p | manager | /api/system/menu/tree | GET | 菜单�?|
| 19 | p | manager | /api/system/api | POST | 注册接口 |
| 20 | p | manager | /api/system/api | GET | 接口列表 |
| 21 | p | manager | /api/system/api/:id | PUT | 更新接口 |
| 22 | p | manager | /api/system/api/:id | DELETE | 删除接口 |
| 23 | p | manager | /api/system/api/all | GET | 所有接�?|
| 24 | p | manager | /api/system/log/login | GET | 登录日志 |
| 25 | p | manager | /api/system/log/oper | GET | 操作日志 |
| | | | | | |
| | | | **操作�?-- 业务操作** | | |
| 26 | p | operator | /api/plating/event/dosing | POST | 录入加药 |
| 27 | p | operator | /api/plating/event/dosing | GET | 加药列表 |
| 28 | p | operator | /api/plating/event/dosing/:id | DELETE | 删除加药 |
| 29 | p | operator | /api/plating/event/water | POST | 录入换水 |
| 30 | p | operator | /api/plating/event/water | GET | 换水列表 |
| 31 | p | operator | /api/plating/event/water/:id | DELETE | 删除换水 |
| 32 | p | operator | /api/plating/event/production | POST | 录入生产 |
| 33 | p | operator | /api/plating/event/production | GET | 生产列表 |
| 34 | p | operator | /api/plating/event/production/:id | DELETE | 删除生产 |
| 35 | p | operator | /api/plating/event/trigger-calc | POST | 触发计算 |
| 36 | p | operator | /api/plating/state/:tankId | GET | 槽体状�?|
| 37 | p | operator | /api/plating/state/trend | GET | 趋势数据 |
| 38 | p | operator | /api/plating/state/export | GET | 导出报表 |
| 39 | p | operator | /api/plating/tank | GET | 槽体列表 |
| 40 | p | operator | /api/plating/tank/:tankId | GET | 槽体详情 |
| 41 | p | operator | /api/system/menu/current | GET | 当前菜单 |
| | | | | | |
| | | | **查看�?-- 只读** | | |
| 42 | p | viewer | /api/plating/state/:tankId | GET | 槽体状�?|
| 43 | p | viewer | /api/plating/state/trend | GET | 趋势数据 |
| 44 | p | viewer | /api/plating/state/export | GET | 导出报表 |
| 45 | p | viewer | /api/plating/event/dosing | GET | 加药记录 |
| 46 | p | viewer | /api/plating/event/water | GET | 换水记录 |
| 47 | p | viewer | /api/plating/event/production | GET | 生产记录 |
| 48 | p | viewer | /api/plating/tank | GET | 槽体列表 |
| 49 | p | viewer | /api/plating/tank/:tankId | GET | 槽体详情 |
| 50 | p | viewer | /api/system/menu/current | GET | 当前菜单 |

### 9.6 四种角色的权限对比总结

| 功能模块 | admin | manager | operator | viewer |
|---------|-------|---------|----------|--------|
| 用户管理（增删改查） | 全部 | 全部 | �?| �?|
| 角色管理（增删改查） | 全部 | 全部 | �?| �?|
| 菜单管理（增删改�?| 全部 | 全部 | �?| �?|
| 接口管理（增删改查） | 全部 | 全部 | �?| �?|
| 日志查看 | 全部 | 只看 | �?| �?|
| 日志清空 | 全部 | **�?* | �?| �?|
| 字典管理 | 全部 | �?| �?| �?|
| 文件管理 | 全部 | �?| �?| �?|
| 事件录入（加�?换水/生产�?| 全部 | �?| �?�?�?| **只查** |
| 触发计算 | 全部 | �?| �?| �?|
| 槽体状�?趋势/导出 | 全部 | �?| 只读 | 只读 |
| 槽体配置（增删改查） | 全部 | �?| **只查** | 只查 |
| 槽体初始化模型状�?| 全部 | �?| �?| �?|

**设计思路**�?- `admin`：一条通配规则搞定，拥有一切权限，包括字典管理、文件管理、日志清空等敏感操作
- `manager`：专注系统管理，能管人、管角色、管菜单、管接口，但不碰业务数据，也不能清空日志
- `operator`：专注业务操作，能录入事件、查看状态、触发计算，但不能管理系统，也不能增删槽体配�?- `viewer`：纯只读，所有权限都�?GET，任何写操作都会被拦�?
---

## 十、本项目中完整的权限配置流程（通过API操作�?
### 场景：给"操作�?角色配置"可以录入加药事件"的权�?
**第一步：在系统中注册 API 接口**（只需做一次）

```
POST /api/system/api
{
    "path": "/api/plating/event/dosing",
    "method": "POST",
    "group": "槽液事件",
    "description": "录入加药事件"
}
```

接口信息写入 `sys_api` 表，得到 id=10�?
**第二步：确认"操作�?角色存在**

```
POST /api/system/role
{
    "name": "操作�?,
    "code": "operator",
    "status": 1
}
```

角色编码 `operator` 将作�?Casbin 中的 sub�?
**第三步：为角色绑�?API 权限**

这一步会同时操作两张表：

1. �?`sys_role_api` 表中记录"角色2绑定了接�?0"（给前端展示用）
2. �?`casbin_rule` 表中写入策略 `("p", "operator", "/api/plating/event/dosing", "POST")`（给Casbin鉴权用）
3. 调用 `enforcer.LoadPolicy()` 刷新内存中的策略

**第四步：为用户分�?操作�?角色**

�?userId �?roleId 的关系写�?`sys_user_role` 表�?
**第五步：用户发起请求，自动鉴�?*

```
POST /api/plating/event/dosing
Authorization: Bearer eyJhbGci...
```

CasbinMiddleware 自动完成鉴权，业务代码无需关心权限�?
---

## 十一、常�?Casbin API

本项目在 `pkg/casbin/casbin.go` 中封装了常用操作�?
### 11.1 添加策略

```go
// 单条添加：允�?operator 角色 POST /api/plating/event/dosing
casbin.AddPolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// 批量添加（全量覆盖模式：先删旧的，再写新的）
rules := [][]string{
    {"/api/plating/event/dosing", "POST"},
    {"/api/plating/event/dosing", "GET"},
    {"/api/plating/state/:tankId", "GET"},
}
casbin.AddRolePolicies(enforcer, "operator", rules)
```

### 11.2 删除策略

```go
// 删除单条
casbin.RemovePolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// 删除某角色的所有策略（删除角色时用�?casbin.RemoveAllPoliciesForRole(enforcer, "operator")
```

### 11.3 查询策略

```go
// 查询某角色的所有权�?policies, _ := casbin.GetRolePolicies(enforcer, "operator")
// 返回: [["operator", "/api/plating/event/dosing", "POST"], ...]

// 获取所有策略（调试用）
allPolicies := enforcer.GetPolicy()
```

### 11.4 检查权�?
```go
// 检�?operator 是否有权�?POST /api/plating/event/dosing
allowed, _ := casbin.CheckPermission(enforcer, "operator", "/api/plating/event/dosing", "POST")
```

### 11.5 重新加载策略（重要）

```go
// 每次修改casbin_rule表后必须调用，否则内存中的策略不会更�?casbin.ReloadPolicy(enforcer)
```

---

## 十二、常见问题排�?
### Q1：添加了策略�?Enforce 返回 false�?
**最常见原因**：没有调�?`enforcer.LoadPolicy()` 重新加载�?
Casbin 在内存中匹配，数据库写入后必须重新加载才生效。每次通过 API 修改策略后，代码中已经自动调用了 `LoadPolicy()`。但如果你直接操作数据库（比如手�?INSERT），必须重启服务或手动触发加载�?
### Q2：路径带参数时匹配不上？

本项目的 matchers 使用 `keyMatch2`，支�?`:param` 风格通配�?
**正确的策略写�?*�?
```sql
-- 策略中用 :id 占位
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'operator', '/api/system/user/:id', 'GET');
```

这样 `/api/system/user/123`、`/api/system/user/456` 都能匹配�?
**常见错误**：策略中写了精确的ID，如 `/api/system/user/123`，这样只能匹配用�?23�?
### Q3：多个角色时如何判定�?
用户可以同时拥有多个角色（比如既�?operator 又是 viewer）。中间件�?*逐个角色检�?*，只要有一个角色有权限就放行（OR 逻辑）�?
```
用户角色: [operator, viewer]

检�?operator �?enforcer.Enforce("operator", path, method) �?false
检�?viewer  �?enforcer.Enforce("viewer", path, method) �?true
�?有权限，放行
```

### Q4：角色被删除后，casbin_rule 中的规则会残留吗�?
删除角色时，应该同步清理 Casbin 策略�?
```go
// 删除该角色在 casbin_rule 表中的所有规�?casbin.RemoveAllPoliciesForRole(enforcer, roleCode)
casbin.ReloadPolicy(enforcer)
```

### Q5：如何调试当前内存中有哪些策略？

```go
// 打印所有策�?policies := enforcer.GetPolicy()
for _, p := range policies {
    fmt.Printf("角色:%s  路径:%s  方法:%s\n", p[0], p[1], p[2])
}
```

### Q6：为什么我直接在数据库加了规则但不生效�?
因为 Casbin 在内存中做匹配。你往数据库插了数据，但内存还是旧的。两种解法：
1. 重启服务（服务启动时�?`LoadPolicy()`�?2. 通过代码调用 `enforcer.LoadPolicy()` 重新加载

**建议通过系统的API接口来管理权限，API内部会自动刷新策略�?*

### Q7：路径中 `*` �?`:param` 的区别？

| 模式 | 示例 | 匹配范围 |
|------|------|---------|
| `:param` | `/api/system/user/:id` | 只匹配一层：`/api/system/user/123` |
| `*` | `/api/*` | 匹配所有子路径：`/api/system/user/123/detail` 也能匹配 |

所�?admin �?`/api/*` 能匹配所�?API 接口�?
---

## 十三、进阶：角色继承（role_definition�?
### 13.1 什么是角色继承

当前 `rbac_model.conf` 中已经声明了 `g = _, _`，表示支持角色继承，但项目目前没有使用�?
角色继承解决的问题：**避免重复配置权限**。比�?经理"应该拥有"操作�?的所有权限，再加上一些管理权限。如果不用继承，就得把操作员的每条规则都给经理复制一遍�?
### 13.2 如何启用角色继承

�?`casbin_rule` 表中插入 `ptype = 'g'` 的记录即可，不需要改 model.conf�?
```sql
-- 表示 manager 继承 operator 的所有权�?INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'manager', 'operator');
```

| ptype | v0 (子角�? | v1 (父角�? | 含义 |
|-------|------------|------------|------|
| g | manager | operator | manager 拥有 operator 的所有权�?|

### 13.3 继承后的鉴权过程

假设 casbin_rule 表中有以下数据：

```
ptype=p, v0=operator, v1=/api/plating/event/dosing, v2=POST    -- operator 能录入加�?ptype=p, v0=manager,  v1=/api/system/user,          v2=GET     -- manager 能查看用户列�?ptype=g, v0=manager,  v1=operator                               -- manager 继承 operator
```

�?manager 请求 `POST /api/plating/event/dosing` 时：

```
matchers: g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")

1. 检查策�? (operator, /api/plating/event/dosing, POST)
   - g("manager", "operator") �?manager 继承�?operator �?true
   - keyMatch2("/api/plating/event/dosing", "/api/plating/event/dosing") �?true
   - "POST" == "POST" �?true
   - 三个条件都满�?�?匹配成功 �?放行
```

**manager 没有直接配置加药权限，但通过继承 operator 自动获得了�?*

### 13.4 多级继承

支持链式继承�?
```sql
INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'operator', 'viewer');
INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'manager', 'operator');
```

继承链：`manager �?operator �?viewer`

- viewer 只能查看
- operator 拥有 viewer 的权�?+ 自己的录入权�?- manager 拥有 operator 的权限（包含 viewer 的）+ 自己的管理权�?
### 13.5 通过代码管理角色继承

```go
// 添加角色继承关系：manager 继承 operator
enforcer.AddGroupingPolicy("manager", "operator")

// 删除角色继承关系
enforcer.RemoveGroupingPolicy("manager", "operator")

// 查询某角色继承了哪些角色
roles, _ := enforcer.GetRolesForUser("manager")
// 返回: ["operator"]

// 查询某角色被哪些角色继承
users, _ := enforcer.GetUsersForRole("operator")
// 返回: ["manager"]
```

---

## 十四、进阶：policy_effect 策略效果详解

### 14.1 所有可用的 effect 规则

| 规则 | 含义 | 模式名称 | 适用场景 |
|------|------|---------|---------|
| `some(where (p.eft == allow))` | 有任意一条允许就放行 | 白名单模�?| **本项目当前使�?*，最常用 |
| `!some(where (p.eft == deny))` | 没有任何一条拒绝就放行 | 黑名单模�?| 默认全部允许，只�?禁止�? |
| `some(where (p.eft == allow)) && !some(where (p.eft == deny))` | 有允许且没有拒绝才放�?| 允许+拒绝并存 | 需要精细控�?|
| `priority(p.eft) \|\| deny` | 按优先级决定，无匹配则拒�?| 优先级模�?| 策略有优先级排序 |

### 14.2 白名单模式（当前项目�?
```ini
e = some(where (p.eft == allow))
```

- 默认拒绝一�?- 只有�?casbin_rule �?*明确配置了允许规�?*的请求才放行
- 安全性最高，推荐大多数项目使�?
### 14.3 黑名单模�?
```ini
e = !some(where (p.eft == deny))
```

- 默认允许一�?- 只有�?casbin_rule �?*明确配置了拒绝规�?*的请求才拦截
- 适合"大部分接口都公开，只有少数需要禁�?的场�?- 不推荐用于管理系统（安全风险高）

### 14.4 允许+拒绝并存模式（实战案例）

```ini
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
```

使用此模式时，策略需要带�?`eft`（effect）字段。model.conf 需要改为：

```ini
[policy_definition]
p = sub, obj, act, eft
```

casbin_rule 表中的数据示例：

| ptype | v0 | v1 | v2 | v3 |
|-------|----|----|----|----|
| p | operator | /api/plating/* | * | allow |
| p | operator | /api/plating/tank/:tankId | DELETE | deny |

含义：operator 可以访问所有业务接口，**但禁止删除槽�?*�?
```
请求: operator DELETE /api/plating/tank/T001

1. 匹配�?allow 规则: (operator, /api/plating/*, *, allow) �?路径匹配
2. 匹配�?deny 规则:  (operator, /api/plating/tank/:tankId, DELETE, deny) �?路径和方法都匹配
3. effect 判定: �?allow 但也�?deny �?deny 优先 �?拒绝
```

**注意**：启用此模式需要修�?policy_definition 并重新设�?casbin_rule 表的数据，是一个较大的改动。当前项目的白名单模式已经够用，建议等确实需�?禁止特定操作"的需求时再升级�?
---

## 十五、casbin_rule 表的增删改查操作方式

### 15.1 推荐方式：通过系统 API 接口操作（强烈推荐）

这是本项目设计的正规流程，通过角色管理接口间接维护 casbin_rule�?
```
管理员在前端页面操作
    �?调用"创建/更新角色"接口，传�?apiIds
    �?后端代码同时写入 sys_role_api �?casbin_rule
    �?自动调用 enforcer.LoadPolicy() 刷新内存
```

**优点**：数据库和内存自动保持同步，不会出现"改了数据库但不生�?的问题�?
### 15.2 通过代码操作（开�?调试时使用）

项目�?`pkg/casbin/casbin.go` 中封装了完整�?CRUD 操作�?
```go
// �?-- 添加单条策略
casbin.AddPolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// �?-- 批量替换某角色的所有策略（先删旧的，再写新的）
rules := [][]string{
    {"/api/plating/event/dosing", "POST"},
    {"/api/plating/event/dosing", "GET"},
    {"/api/plating/state/:tankId", "GET"},
}
casbin.AddRolePolicies(enforcer, "operator", rules)

// �?-- 删除单条策略
casbin.RemovePolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// �?-- 删除某角色的所有策�?casbin.RemoveAllPoliciesForRole(enforcer, "operator")

// �?-- 查询某角色的所有策�?policies, _ := casbin.GetRolePolicies(enforcer, "operator")
// 返回: [["operator", "/api/plating/event/dosing", "POST"], ...]

// �?-- 查询所有策�?allPolicies := enforcer.GetPolicy()

// �?-- 没有直接�?Update，用"删旧增新"实现
// AddRolePolicies 内部就是这个逻辑：先 RemoveAll，再批量 Add

// 鉴权检�?allowed, _ := casbin.CheckPermission(enforcer, "operator", "/api/plating/event/dosing", "POST")

// 刷新内存（通过上述封装函数操作时会自动刷新，一般不需要手动调�?casbin.ReloadPolicy(enforcer)
```

### 15.3 直接操作数据库（不推荐，仅限紧急排查）

可以直接 SQL 操作 casbin_rule 表，�?*必须重启服务**或手动触�?`LoadPolicy()` 才能生效�?
```sql
-- 查看所有规�?SELECT * FROM casbin_rule;

-- 查看某角色的规则
SELECT * FROM casbin_rule WHERE ptype = 'p' AND v0 = 'operator';

-- 手动添加规则（添加后必须重启服务或调�?LoadPolicy�?INSERT INTO casbin_rule (ptype, v0, v1, v2) 
VALUES ('p', 'operator', '/api/plating/event/dosing', 'POST');

-- 手动删除规则（删除后必须重启服务或调�?LoadPolicy�?DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 = 'operator' AND v1 = '/api/plating/event/dosing';
```

**再次强调**：直接操作数据库后如果不重启服务，内存中的策略不会更新，鉴权结果不会改变�?
---

## 十六、sys_api 接口注册�?Casbin 鉴权联动机制

### 16.1 整体关系�?
```
                          前端管理页面
                              �?            ┌─────────────────┼─────────────────�?            �?                �?                �?    注册API接口          创建/编辑角色        分配用户角色
  POST /api/system/api   PUT /api/system/role  (sys_user_role)
            �?                �?            �?                �?        sys_api�?        同时写入两张�?
    (接口注册表，供前�?    ├→ sys_role_api（给前端展示勾选状态）
     展示"有哪些接�?      └→ casbin_rule （给Casbin引擎鉴权�?     可以分配给角�?)            �?                               �?                        enforcer.LoadPolicy()
                        （刷新内存中的策略）
                               �?            ┌──────────────────�?            �?    用户发起业务请求
    POST /api/plating/event/dosing
            �?            �?    AuthMiddleware（JWT认证 �?得到userId�?            �?            �?    CasbinMiddleware
    ├─ �?sys_user_role �?得到 roleIds
    ├─ �?sys_role �?得到 role.Code
    └─ enforcer.Enforce(role.Code, path, method)
       └─ 在内存中匹配 casbin_rule �?放行/拒绝
```

### 16.2 各张表的分工

| �?| 存什�?| 谁写�?| 谁读�?| 作用 |
|---|---|---|---|---|
| `sys_api` | 接口的路径、方法、分组、描�?| 管理员通过接口注册API | 前端角色编辑页面 | 展示"有哪些接口可以分�? |
| `sys_role` | 角色名称、编�?code)、状�?| 管理员创建角�?| CasbinMiddleware | 提供角色编码给Casbin |
| `sys_user_role` | 用户ID + 角色ID | 管理员分配角�?| CasbinMiddleware | 查询用户有哪些角�?|
| `sys_role_api` | 角色ID + 接口ID | 管理员为角色分配接口 | 前端角色编辑页面 | 展示"这个角色勾选了哪些接口" |
| `casbin_rule` | 角色编码 + 路径 + 方法 | 管理员为角色分配接口（代码自动同步） | Casbin引擎（内存匹配） | 运行时鉴权判�?|

### 16.3 完整操作流程（从注册接口到鉴权生效）

**场景**：系统新增了一个接�?导出加药报表"，需要让操作员能用�?
**第一步：注册接口�?sys_api**

```
POST /api/system/api
{
    "apiName": "导出加药报表",
    "apiPath": "/api/plating/event/dosing/export",
    "method": "GET",
    "group": "槽液事件",
    "remark": "导出加药事件数据为Excel"
}
�?写入 sys_api 表，得到 api_id = 15
```

此时只是"登记"了这个接口的存在，还没有分配给任何角色�?
**第二步：编辑"操作�?角色，勾选新接口**

管理员在前端打开"操作�?角色编辑页面，看到接口列表（�?sys_api 查出来的），勾�?导出加药报表"，点保存�?
```
PUT /api/system/role/2
{
    "name": "操作�?,
    "code": "operator",
    "apiIds": [3, 4, 5, 6, 7, 8, 9, 10, 15]  // 原有�?+ 新增�?5
}
```

后端代码执行�?1. 更新 sys_role 表的角色信息
2. 清空 sys_role_api �?role_id=2 的旧记录，写入新的关�?3. 调用 `casbin.AddRolePolicies(enforcer, "operator", rules)` —�?自动清空旧策略，写入新策略到 casbin_rule，并刷新内存

**第三步：鉴权自动生效**

操作员张三请�?`GET /api/plating/event/dosing/export`，CasbinMiddleware 自动匹配到新规则，放行�?
### 16.4 sys_role �?Code 字段说明

Code �?*自定义的角色编码**，没有固定的枚举值，由项目按业务需要定义。但有以下约束：

- 全局唯一（数据库有唯一索引 `uk_sys_role_code`�?- 使用英文小写 + 下划线，见名知义
- 它就�?Casbin 中的 `sub`（主体），直接写�?casbin_rule 表的 v0 字段

本项目建议的角色编码规范�?
| code | name | 权限范围 |
|------|------|---------|
| `admin` | 超级管理�?| 通配所有接口（`/api/*` + `*`�?|
| `manager` | 管理�?| 用户管理、角色管理、日志查看等系统管理功能 |
| `operator` | 操作�?| 业务数据录入（加药、换水、生产事件）+ 状态查�?|
| `viewer` | 查看�?| 只读权限：查看状态、趋势、导出报�?|

可以根据业务需要自由扩展，比如 `quality_inspector`（质检员）、`shift_leader`（班组长）等�?
---

## 十七、sys_role_api �?casbin_rule 的关系详解（为什么不是重复）

### 17.1 两张表存储的数据格式完全不同

**sys_role_api** 存的是业务ID的关联：

| id | role_id | api_id |
|----|---------|--------|
| 1 | 2 | 10 |

**casbin_rule** 存的是鉴权三元组�?
| ptype | v0 | v1 | v2 |
|-------|----|----|-----|
| p | operator | /api/plating/event/dosing | POST |

同一�?操作员能录入加药"的权限，在两张表中的表达方式完全不同�?
### 17.2 两张表的使用者不�?
```
┌─────────────────────────────────────────────�?�?              前端管理页面                    �?�?                                            �?�? 角色编辑页面需要展示：                       �?�? "操作�?当前勾选了哪些接口�?                 �?�?                                            �?�? 查询: SELECT api_id FROM sys_role_api       �?�?       WHERE role_id = 2                     �?�? �?得到 api_id: [3, 4, 5, 10]               �?�? �?�?api_id 关联 sys_api 表展示接口名�?     �?�?                                            �?�? Casbin引擎不认�?api_id，它需要的是：         �?�? ("operator", "/api/plating/event/dosing",   �?�?  "POST") 这种三元组格�?                     �?�?                                            �?�? 所�?casbin_rule 表存了一份Casbin能理解的数据  �?└─────────────────────────────────────────────�?```

### 17.3 为什么不能只用一张表

**假设只保�?casbin_rule，去�?sys_role_api**�?- 前端角色编辑页面需要展�?勾选了哪些接口"
- casbin_rule 里存的是 (operator, /api/plating/event/dosing, POST)
- 前端需要展示接口名称、分组等信息，就得用 path+method 反查 sys_api
- 这种反查不可靠：如果接口路径改了，关联就断了
- �?ID 关联（role_id + api_id）才是关系数据库的正确做�?
**假设只保�?sys_role_api，去�?casbin_rule**�?- CasbinMiddleware 鉴权时需�?(role_code, path, method) 三元�?- 如果没有 casbin_rule，就得在每次请求时：�?sys_role_api �?�?sys_api 拿到 path �?method �?再做匹配
- 这会在每个请求上增加多次数据库查询，性能极差
- 而且无法使用 Casbin �?keyMatch2 路径通配等高级功�?
### 17.4 一句话总结

| �?| 面向�?| 解决什么问�?|
|---|---|---|
| `sys_role_api` | 面向**前端/业务�?* | �?ID 关联，方便展示和管理 |
| `casbin_rule` | 面向**Casbin 引擎** | 用三元组格式，支持内存高速匹配和路径通配 |

**两张表是同一份权限数据的两种表达形式**，通过"为角色分配API"的业务逻辑保持同步。不是重复，是各司其职�?
---

## 十八、本项目 Casbin 相关文件一览（原十三章�?
| 文件 | 作用 |
|------|------|
| `etc/rbac_model.conf` | RBAC 模型定义文件，定义请求格式、策略格式、匹配规�?|
| `etc/plating-api.yaml` | 配置文件，`CasbinModelPath` 指定模型文件路径 |
| `pkg/casbin/casbin.go` | Casbin 初始�?+ 策略操作封装函数 |
| `internal/middleware/casbin_middleware.go` | HTTP 鉴权中间件，拦截请求并执�?Enforce |
| `internal/middleware/auth_middleware.go` | JWT 认证中间件，�?Casbin 提供 userId |
| `internal/svc/service_context.go` | 初始�?Enforcer 并注入中间件 |
| `internal/handler/routes.go` | 路由注册，决定哪些路由受 Casbin 保护 |
| `internal/model/system/sys_role_model.go` | 角色Model，提供按ID查角色编�?|
| `internal/model/system/sys_user_role_model.go` | 用户角色关联Model，提供按用户ID查角色列�?|
| `internal/model/system/sys_role_api_model.go` | 角色API关联Model，记录角色绑定了哪些接口 |
| `internal/model/system/sys_api_model.go` | API接口Model，存储所有可配置的接口信�?|
| `schema/plating.sql` | 数据库表结构，包�?casbin_rule 表定义和初始数据 |

---

## 十九、快速备忘录

```json
鉴权的三元组
(角色code, 请求路径, HTTP方法)  �? Enforce()  �? true/false

增删策略
enforcer.AddPolicy(role, path, method)           // 添加单条
enforcer.AddPolicies(rules)                       // 批量添加
enforcer.RemovePolicy(role, path, method)         // 删除单条
enforcer.RemoveFilteredPolicy(0, role)            // 删除角色的所有规�?
查询
enforcer.GetPolicy()                              // 所有规�?enforcer.GetFilteredPolicy(0, role)               // 某角色的规则
enforcer.Enforce(role, path, method)              // 执行鉴权

刷新（每次改完casbin_rule后必须调用）
enforcer.LoadPolicy()

路径匹配规则（本项目使用 keyMatch2�?/api/system/user      精确匹配 /api/system/user
/api/system/user/:id  通配匹配 /api/system/user/123（一层）
/api/*                通配匹配 /api/ 下所有路径（多层�?
方法匹配
GET/POST/PUT/DELETE   精确匹配
*                     匹配任意方法
```

# QA�?
当前rbac_model.conf中的[matchers]规则，就是角色匹�?r.sub和p.sub匹配)且请求路由匹�?r.obj和p.obj)且HTTP方法匹配(r.act和p.act)就放行呗�?但是还有几个疑问，如下所示：
**1.我现在只明白了当前rbac_model.conf配置文件里的东西，那如果项目需要更复杂的鉴权规则，或者后面我的项目要设置更严格的鉴权规则，那怎么办？例如需要角色继承，[role_definition], 那应该怎么设计，流程是什么？**

```bash
答：
当前 rbac_model.conf 已经声明�?g = _, _，说明已经预留了角色继承能力，只是还没用�?流程很简单：�?casbin_rule 表中插入 ptype=g 的记录即可。比�?经理继承操作员所有权�?�?INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'manager', 'operator');
这样 manager 自动拥有 operator 的全部权限，不需要重复配 p 规则�?支持多级继承：manager �?operator �?viewer，manager同时拥有 operator �?viewer 的权限�?
这里大家可能有疑问：
新建角色继承时，要在casbin_rule表中插入ptype=g记录，INSERT INTO casbin_rule (ptype,
v0, v1) VALUES ('g', 'manager', 'operator'); 但是当前表中现有的内容type,v0,v1的值是p,admin,/api/*。难道v1字段填写�?operator'之后，就说明'manager'集成了所�?operator'可访问的路由？例如：/api/*

解释�?casbin_rule 表的 v0、v1、v2 这些字段不是固定含义，它们的含义取决�?ptype 的值：
ptype 就是"这行数据该怎么解读"的标�?
�?）当 ptype = 'p'（策略规则）时：
ptype  �?   v0    �?  v1   �?   v2   �? p     �? 角色编码 �?路径   �?HTTP方法   �? p     �?admin    �?/api/* �?*       �? 
�?）ptype = 'g'（角色继承）时：
�?ptype �?  v0    �?   v1    �?  v2   �?�?g     �?子角�?  �? 父角�?  �?(不用)  �?�?g     �?manager �?operator �?(�?   �?
INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'manager', 'operator');
Casbin 读到 ptype='g'，就知道这不是一条路由策略，而是一条继承关系：manager 继承 operator�?然后�?matchers �?g(r.sub, p.sub) 这个函数做匹配时�?
  请求: manager 访问 POST /api/plating/event/dosing
  策略: (p, operator, /api/plating/event/dosing, POST)

  g("manager", "operator")
  �?�?ptype='g' 的记录，发现 manager �?operator
  �?返回 true（manager 等价�?operator�?
  manager 自动拥有 operator 的所�?p 规则，不需要重复配置�?```

(2) [policy_effect]有哪些规则？

![image-20260416093949898](C:\Users\17878\AppData\Roaming\Typora\typora-user-images\image-20260416093949898.png)

(3) 我如果想针对`casbin_rule`表中的规则进行增删改查，通常都是怎么操作,直接操作数据库吗�?
```
正确做法是通过代码中封装的 Casbin API 操作，项目已经在 pkg/casbin/casbin.go 中封装好了：
  - AddPolicyForRole() / AddRolePolicies() —�?�?  - RemovePolicyForRole() / RemoveAllPoliciesForRole() —�?�?  - GetRolePolicies() —�?�?  - �?= 先删后增（AddRolePolicies 内部就是这么做的：先删旧的，再批量写新的�?  这些函数内部会自动同步数据库和内存。在业务层面，就是通过"创建/更新角色时传�?ApiIds"这个接口来触发的
```

  2.当前的api/desc/system/api.api的接口是要结合casbin鉴权一起用的吗？怎么用呢�?
```
sys_api 表是"接口注册�?—�?把系统里所有需要管控权限的接口登记进去（路径、方法、分组、描述）。它本身不参与鉴权，是给前端页面用的，让管理员在界面上看�?系统有哪些接口可以分�?�?  联动流程�?  1. 通过 /api/system/api (POST) 注册接口 �?写入 sys_api �?     （告诉系统："我有这个接口，它可以被分配给角色"�?  2. 通过 创建/更新角色接口，传�?apiIds �?同时写入两张表：
     - sys_role_api（role_id + api_id）→ 前端展示用："这个角色勾选了哪些接口"
     - casbin_rule（role_code + path + method）→ Casbin鉴权�?  3. 请求进来 �?CasbinMiddleware �?casbin_rule 做实时鉴�?```

sys_role 表中的Code字段不能随意写吧，是不是有固定的几个值，有哪些？

```
sys_role �?Code 字段不是随便写的，但也没有固定值——它是你自己定义的角色编码，作为 Casbin 策略中的
  sub。约定是�?  - 使用英文小写 + 下划�?  - 要有语义，方便辨�?  - 全局唯一（有唯一索引约束�?```

 常见的设计：

![image-20260416094500997](C:\Users\17878\AppData\Roaming\Typora\typora-user-images\image-20260416094500997.png)

 3.`sys_role_api` 记录"角色能访问哪些接�?（业务层面的关联关系）和`casbin_rule`记录"角色能访问哪些接�?重复的问题我还是不理解，为啥是重复的呢？

它们不是重复，是各管各的事，分工明确�?
![image-20260416094545529](C:\Users\17878\AppData\Roaming\Typora\typora-user-images\image-20260416094545529.png)

```
举个例子：前端管理员打开"操作�?角色的编辑页面，看到一个接口列表，上面有勾选框。这个勾选状态从哪来？从
  sys_role_api 查的（SELECT api_id FROM sys_role_api WHERE role_id=2）�?
  但是 Casbin 不认�?role_id=2, api_id=10 这种数据，它只认 ("operator", "/api/plating/event/dosing", "POST")
  这种三元组。所以需�?casbin_rule 存一�?Casbin 能理解的格式�?
  两张表通过"分配API权限"的业务逻辑保持同步：管理员勾选接�?�?代码同时�?sys_role_api �?casbin_rule�?
  如果只有 casbin_rule，前端没法用 ID 关联展示；如果只�?sys_role_api，Casbin
  引擎没法做鉴权。所以两张表缺一不可�?```

