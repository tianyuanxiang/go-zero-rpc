# go-zero-admin 接口测试用例

> 基础地址: `http://localhost:8888`
>
> 认证方式: 除登录和刷新令牌外，所有接口需要在 Header 中携带 `Authorization: Bearer {accessToken}`
>
> 测试前提: 先调用登录接口获取 token，后续接口都带上 token
>
> 角色基础 CRUD 和登录登出已跳过，此文档不包含这些接口

---

## 一、认证模块

> 说明: 登录接口已跳过，此处仅包含登录后的认证操作

### 1.1 获取当前用户信息

```
GET /api/auth/userInfo
Authorization: Bearer {accessToken}
```

**用例1: 正常获取** -- 直接发送请求，无参数

**用例2: token 过期** -- 使用一个过期的 token 发送请求

**用例3: 无 token** -- 不携带 Authorization 头

---

### 1.2 刷新令牌

```
POST /api/auth/refresh
Content-Type: application/json
```

**用例1: 正常刷新**
```json
{
    "refreshToken": "{登录时返回的refreshToken}"
}
```

**用例2: 无效的 refreshToken**
```json
{
    "refreshToken": "invalid_token_string"
}
```

**用例3: 空 refreshToken**
```json
{
    "refreshToken": ""
}
```

---

### 1.3 修改密码

```
POST /api/auth/changePassword
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 正常修改密码**
```json
{
    "oldPassword": "123456",
    "newPassword": "Abc@12345"
}
```

**用例2: 旧密码错误**
```json
{
    "oldPassword": "wrong_password",
    "newPassword": "Abc@12345"
}
```

**用例3: 新密码过短**
```json
{
    "oldPassword": "123456",
    "newPassword": "123"
}
```

**用例4: 新密码与旧密码相同**
```json
{
    "oldPassword": "123456",
    "newPassword": "123456"
}
```

**用例5: 缺少必填字段**
```json
{
    "oldPassword": "123456"
}
```

---

### 1.4 登出

```
POST /api/auth/logout
Authorization: Bearer {accessToken}
```

**用例1: 正常登出** -- 直接发送请求，无 body

**用例2: 无 token 登出** -- 不携带 Authorization 头

---

## 二、菜单管理

> 测试顺序: 查询菜单树 -> 创建菜单 -> 获取当前用户菜单 -> 更新菜单 -> 删除菜单
>
> 菜单是其他模块（角色权限分配）的基础数据，优先测试

### 2.1 获取菜单树✅️

```
GET /api/system/menu/tree
Authorization: Bearer {accessToken}
```

**用例1: 正常获取** -- 直接请求，无参数，确认初始化菜单数据存在

---

### 2.2 创建菜单

```
POST /api/system/menu
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 创建目录（顶级）✅️**

```json
{
    "parentId": 0,
    "menuName": "测试目录",
    "menuType": 0,
    "path": "/test",
    "component": "Layout",
    "icon": "example",
    "sort": 99,
    "status": 1,
    "remark": "测试用目录"
}
```

**用例2: 创建菜单（挂在系统管理下）✅️**

```json
{
    "parentId": 1,
    "menuName": "测试页面",
    "menuType": 1,
    "path": "/system/test",
    "component": "system/test/index",
    "icon": "test",
    "sort": 99,
    "perms": "system:test:list",
    "status": 1,
    "remark": "测试用菜单"
}
```

**用例3: 创建按钮（挂在用户管理下）✅️**

```json
{
    "parentId": 2,
    "menuName": "测试按钮",
    "menuType": 2,
    "perms": "system:user:test",
    "sort": 99,
    "status": 1,
    "remark": "测试用按钮"
}
```

**用例4: 缺少菜单名称✅️**

```json
{
    "parentId": 0,
    "menuType": 0,
    "path": "/noname",
    "status": 1,
    "remark": ""
}
```

**用例5: 不存在的父菜单✅️**

```json
{
    "parentId": 99999,
    "menuName": "孤儿菜单",
    "menuType": 1,
    "path": "/orphan",
    "status": 1,
    "remark": ""
}
```

**用例6: 隐藏状态菜单✅️**

```json
{
    "parentId": 1,
    "menuName": "隐藏页面",
    "menuType": 1,
    "path": "/system/hidden",
    "component": "system/hidden/index",
    "sort": 99,
    "status": 0,
    "remark": "隐藏菜单"
}
```

---

### 2.3 获取当前用户菜单✅️

```
GET /api/system/menu/current
Authorization: Bearer {accessToken}
```

**用例1: 正常获取** -- 返回当前登录用户有权限的菜单树，确认刚创建的菜单是否出现

---

### 2.4 更新菜单

```
PUT /api/system/menu/:id
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 更新菜单名称和图标✅️**

```json
// PUT /api/system/menu/{之前创建的菜单id}
{
    "parentId": 1,
    "menuName": "更新后名称",
    "menuType": 1,
    "path": "/system/test",
    "component": "system/test/index",
    "icon": "new-icon",
    "sort": 50,
    "perms": "system:test:list",
    "status": 1,
    "remark": "更新后备注"
}
```

**用例2: 修改父菜单（移动菜单位置）**✅️

```json
// PUT /api/system/menu/{菜单id}
{
    "parentId": 7,
    "menuName": "移动到日志管理下",
    "menuType": 1,
    "path": "/log/test",
    "component": "log/test/index",
    "sort": 99,
    "status": 1,
    "remark": ""
}
```

**用例3: 更新不存在的菜单**✅️

```json
// PUT /api/system/menu/99999
{
    "parentId": 0,
    "menuName": "不存在",
    "menuType": 0,
    "status": 1,
    "remark": ""
}
```

**用例4: 将菜单设为不可见**✅

```
// PUT /api/system/menu/7
{
    "visible": 0
}
```

### 2.5 删除菜单

```
DELETE /api/system/menu/:id
Authorization: Bearer {accessToken}
```

**用例1: 删除叶子菜单（按钮/无子菜单）✅**

```
DELETE /api/system/menu/{之前创建的按钮菜单id}
```

**用例2: 删除有子菜单的菜单（应失败或级联删除）✅**

```
DELETE /api/system/menu/1
```

**用例3: 删除不存在的菜单✅**

```
DELETE /api/system/menu/99999
```

---

## 三、接口管理

> 测试顺序: 创建接口 -> 接口列表 -> 全部接口 -> 更新接口 -> 删除接口
>
> 接口也是角色权限分配的基础数据

### 3.1 创建接口

```
POST /api/system/api
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 完整参数**✅️

```json
{
    "apiName": "测试接口-获取列表",
    "apiPath": "/api/test/list",
    "method": "GET",
    "group": "test",
    "remark": "测试用接口"
}
```

**用例2: 最小必填参数**✅️

```json
{
    "apiName": "测试接口-创建",
    "apiPath": "/api/test/create",
    "method": "POST",
    "group": "test"
}
```

**用例3: 带路径参数的接口**✅️

```json
{
    "apiName": "测试接口-详情",
    "apiPath": "/api/test/:id",
    "method": "GET",
    "group": "test",
    "remark": "根据ID获取详情"
}
```

**用例4: 各种 HTTP 方法 -- PUT✅️**

```json
{
    "apiName": "测试接口-更新",
    "apiPath": "/api/test/:id",
    "method": "PUT",
    "group": "test"
}
```

**用例5: 各种 HTTP 方法 -- DELETE✅️**

```json
{
    "apiName": "测试接口-删除",
    "apiPath": "/api/test/:id",
    "method": "DELETE",
    "group": "test"
}
```

**用例6: 缺少接口路径✅️**

```json
{
    "apiName": "缺路径接口",
    "method": "GET",
    "group": "test"
}
```

**用例7: 缺少 HTTP 方法✅️**

```json
{
    "apiName": "缺方法接口",
    "apiPath": "/api/test/nomethod",
    "group": "test"
}
```

---

### 3.2 接口列表

```
GET /api/system/api
Authorization: Bearer {accessToken}
```

**用例1: 默认分页**✅️

```
GET /api/system/api
```

**用例2: 指定分页**✅️

```
GET /api/system/api?page=1&pageSize=5
```

**用例3: 按关键词搜索**✅️

```
GET /api/system/api?keyword=user
```

**用例4: 按分组过滤**✅️

```
GET /api/system/api?group=system
```

**用例5: 组合条件**✅️

```
GET /api/system/api?page=1&pageSize=10&keyword=list&group=system
```

---

### 3.3 获取全部接口（不分页）✅️

```
GET /api/system/api/all
Authorization: Bearer {accessToken}
```

**用例1: 正常获取** -- 直接请求，用于下拉选择

---

### 3.4 更新接口

```
PUT /api/system/api/:id
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 更新全部字段**✅️

```json
// PUT /api/system/api/{接口id}
{
    "apiName": "更新后接口名",
    "apiPath": "/api/test/updated",
    "method": "POST",
    "group": "test_updated",
    "remark": "更新后备注"
}
```

**用例2: 只更新名称和备注**✅️

```json
// PUT /api/system/api/{接口id}
{
    "apiName": "只改名称",
    "apiPath": "/api/test/list",
    "method": "GET",
    "group": "test",
    "remark": "只更新了名称"
}
```

**用例3: 更新不存在的接口**✅️

```json
// PUT /api/system/api/99999
{
    "apiName": "不存在",
    "apiPath": "/api/nothing",
    "method": "GET",
    "group": "none"
}
```

---

### 3.5 删除接口

```
DELETE /api/system/api/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常删除**✅️

```
DELETE /api/system/api/{之前创建的测试接口id}
```

**用例2: 删除不存在的接口**✅️

```
DELETE /api/system/api/99999
```

---

## 四、字典管理

> 测试顺序: 创建字典类型 -> 字典类型列表 -> 创建字典数据 -> 按类型查询字典数据 -> 更新字典类型 -> 更新字典数据 -> 删除字典数据 -> 删除字典类型
>
> 先创建类型再创建数据，先删除数据再删除类型

### 4.1 创建字典类型

```
POST /api/system/dict/type
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 完整参数**✅️

```json
{
    "dictName": "测试字典类型",
    "dictCode": "test_dict_type",
    "status": 1,
    "remark": "用于测试的字典类型"
}
```

**用例2: 最小必填参数**✅️

```json
{
    "dictName": "最简字典",
    "dictCode": "test_simple",
    "status": 1
}
```

**用例3: 禁用状态**✅️

```json
{
    "dictName": "禁用字典",
    "dictCode": "test_disabled",
    "status": 0,
    "remark": "创建时即禁用"
}
```

**用例4: 重复的 dictCode**✅️

```json
{
    "dictName": "重复类型",
    "dictCode": "sys_user_status",
    "status": 1
}
```

**用例5: 缺少 dictName✅️**

```json
{
    "dictCode": "test_no_name",
    "status": 1
}
```

---

### 4.2 字典类型列表

```
GET /api/system/dict/type
Authorization: Bearer {accessToken}
```

**用例1: 默认分页**✅️

```
GET /api/system/dict/type
```

**用例2: 指定分页**✅️

```
GET /api/system/dict/type?page=1&pageSize=5
```

**用例3: 关键词搜索**✅️

```
GET /api/system/dict/type?keyword=用户
```

**用例4: 组合条件**✅️

```
GET /api/system/dict/type?page=1&pageSize=10&keyword=状态
```

---

### 4.3 创建字典数据

```
POST /api/system/dict/data
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 完整参数**✅️

```json
{
    "dictTypeId": 1,
    "dictLabel": "测试标签",
    "dictValue": "test_value",
    "sort": 1,
    "status": 1,
    "remark": "测试字典数据"
}
```

**用例2: 最小必填参数**✅️

```json
{
    "dictTypeId": 1,
    "dictLabel": "简单标签",
    "dictValue": "simple",
    "status": 1
}
```

**用例3: 不存在的字典类型ID**✅️

```json
{
    "dictTypeId": 99999,
    "dictLabel": "孤儿数据",
    "dictValue": "orphan",
    "status": 1
}
```

**用例4: 禁用状态**✅️

```json
{
    "dictTypeId": 1,
    "dictLabel": "禁用标签",
    "dictValue": "disabled_value",
    "sort": 99,
    "status": 0
}
```

**用例5: 缺少 dictLabel**✅️

```json
{
    "dictTypeId": 1,
    "dictValue": "no_label",
    "status": 1
}
```

---

### 4.4 按字典类型查询字典数据

```
GET /api/system/dict/data/:dictType
Authorization: Bearer {accessToken}
```

**用例1: 查询存在的字典类型**✅️

```
GET /api/system/dict/data/sys_user_status
```

**用例2: 查询不存在的字典类型**✅️

```
GET /api/system/dict/data/not_exist_type
```

**用例3: 查询菜单类型字典**✅️

```
GET /api/system/dict/data/sys_menu_type
```

---

### 4.5 更新字典类型

```
PUT /api/system/dict/type/:id
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 更新全部字段**✅️

```json
// PUT /api/system/dict/type/{字典类型id}
{
    "dictName": "更新后字典名",
    "dictCode": "test_dict_type_updated",
    "status": 1,
    "remark": "更新后备注"
}
```

**用例2: 只更新名称**✅️

```json
// PUT /api/system/dict/type/{字典类型id}
{
    "dictName": "只改名称",
    "dictCode": "test_dict_type",
    "status": 1
}
```

**用例3: 禁用字典类型**✅️

```json
// PUT /api/system/dict/type/{字典类型id}
{
    "dictName": "测试字典类型",
    "dictCode": "test_dict_type",
    "status": 0
}
```

**用例4: 更新不存在的字典类型**✅️

```json
// PUT /api/system/dict/type/99999
{
    "dictName": "不存在",
    "dictCode": "not_exist",
    "status": 1
}
```

---

### 4.6 更新字典数据

```
PUT /api/system/dict/data/:id
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 更新全部字段**✅️

```json
// PUT /api/system/dict/data/{字典数据id}
{
    "dictTypeId": 1,
    "dictLabel": "更新后标签",
    "dictValue": "updated_value",
    "sort": 10,
    "status": 1,
    "remark": "更新后备注"
}
```

**用例2: 只更新标签**✅️

```json
// PUT /api/system/dict/data/{字典数据id}
{
    "dictTypeId": 1,
    "dictLabel": "只改标签",
    "dictValue": "test_value",
    "status": 1
}
```

**用例3: 更新不存在的字典数据**✅️

```json
// PUT /api/system/dict/data/99999
{
    "dictTypeId": 1,
    "dictLabel": "不存在",
    "dictValue": "none",
    "status": 1
}
```

---

### 4.7 删除字典数据

```
DELETE /api/system/dict/data/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常删除**✅️

```
DELETE /api/system/dict/data/{之前创建的字典数据id}
```

**用例2: 删除不存在的字典数据**✅️

```
DELETE /api/system/dict/data/99999
```

---

### 4.8 删除字典类型

```
DELETE /api/system/dict/type/:id
Authorization: Bearer {accessToken}
```

**用例1: 删除有字典数据关联的类型（测试级联行为）**✅️

```
DELETE /api/system/dict/type/1
```

**用例2: 正常删除（无关联数据的类型）**✅️

```
DELETE /api/system/dict/type/{之前创建的字典类型id}
```

**用例3: 删除不存在的字典类型**✅️

```
DELETE /api/system/dict/type/99999
```

---

## 五、用户管理

> 测试顺序: 创建用户 -> 用户列表 -> 用户详情 -> 更新用户 -> 重置密码 -> 删除用户
>
> 依赖: 角色数据（用于 roleIds 关联）

### 5.1 创建用户

```
POST /api/system/user
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 完整参数创建✅️**

```json
{
    "username": "testuser01",
    "password": "Test@12345",
    "nickname": "测试用户01",
    "email": "test01@example.com",
    "phone": "13800138001",
    "status": 1,
    "avatar": "https://example.com/avatar.png",
    "remark": "测试创建的用户",
    "roleIds": [2, 3]
}
```

**用例2: 最小必填参数**✅️

```json
{
    "username": "testuser02",
    "password": "Test@12345",
    "status": 1
}
```

**用例3: 用户名重复**✅️

```json
{
    "username": "admin",
    "password": "Test@12345",
    "status": 1
}
```

**用例4: 不存在的角色ID✅️**

```json
{
    "username": "testuser03",
    "password": "Test@12345",
    "status": 1,
    "roleIds": [9999]
}
```

**用例5: 创建禁用状态用户**✅️

```json
{
    "username": "testuser04",
    "password": "Test@12345",
    "nickname": "禁用用户",
    "status": 0
}
```

**用例6: 缺少用户名✅️**

```json
{
    "password": "Test@12345",
    "status": 1
}
```

**用例7: 缺少密码✅️**

```json
{
    "username": "testuser05",
    "status": 1
}
```

---

### 5.2 用户列表

```
GET /api/system/user
Authorization: Bearer {accessToken}
```

**用例1: 默认分页✅️**

```
GET /api/system/user
```

**用例2: 指定分页**✅️

```
GET /api/system/user?page=1&pageSize=5
```

**用例3: 第二页**✅️

```
GET /api/system/user?page=2&pageSize=5
```

**用例4: 按关键词搜索✅️**

```
GET /api/system/user?keyword=admin
```

**用例5: 按状态筛选 -- 启用**✅️

```
GET /api/system/user?status=1
```

**用例6: 按状态筛选 -- 禁用✅️**

```
GET /api/system/user?status=0
```

**用例7: 组合条件**✅️

```
GET /api/system/user?page=1&pageSize=10&keyword=test&status=1
```

**用例8: 超大页码（无数据）**✅️

```
GET /api/system/user?page=9999&pageSize=20
```

---

### 5.3 用户详情

```
GET /api/system/user/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常查询**✅️

```
GET /api/system/user/1
```

**用例2: 不存在的用户✅️**

```
GET /api/system/user/99999
```

---

### 5.4 更新用户

```
PUT /api/system/user/:id
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 更新全部字段✅️**

```json
// PUT /api/system/user/2
{
    "nickname": "更新后昵称",
    "email": "updated@example.com",
    "phone": "13900139001",
    "status": 1,
    "avatar": "https://example.com/new-avatar.png",
    "remark": "更新备注",
    "roleIds": [2]
}
```

**用例2: 只更新昵称✅️**

```json
// PUT /api/system/user/2
{
    "nickname": "只改昵称",
    "status": 1
}
```

**用例3: 禁用用户✅️**

```json
// PUT /api/system/user/2
{
    "status": 0
}
```

**用例4: 更新角色关联✅️**

```json
// PUT /api/system/user/2
{
    "status": 1,
    "roleIds": [2, 3]
}
```

**用例5: 清空角色✅️**

```json
// PUT /api/system/user/2
{
    "status": 1,
    "roleIds": []
}
```

**用例6: 更新不存在的用户**✅️

```json
// PUT /api/system/user/99999
{
    "nickname": "不存在",
    "status": 1
}
```

---

### 5.5 重置密码

```
POST /api/system/user/:id/reset-password
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 正常重置**✅️

```json
// POST /api/system/user/2/reset-password
{
    "newPassword": "Reset@12345"
}
```

**用例2: 重置不存在用户的密码**✅️

```json
// POST /api/system/user/99999/reset-password
{
    "newPassword": "Reset@12345"
}
```

**用例3: 空密码✅️**

```json
// POST /api/system/user/2/reset-password
{
    "newPassword": ""
}
```

---

### 5.6 删除用户

```
DELETE /api/system/user/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常删除**✅️

```
DELETE /api/system/user/{之前创建的testuser的id}
```

**用例2: 删除不存在的用户✅️**

```
DELETE /api/system/user/99999
```

**用例3: 删除管理员用户（应被禁止或特殊处理）**✅️

```
DELETE /api/system/user/1
```

---

## 六、角色权限分配

> 测试顺序: 角色详情 -> 全部角色列表 -> 更新角色权限
>
> 依赖: 菜单数据、接口数据（前面模块已创建）
>
> 角色基础 CRUD 已跳过

### 6.1 角色详情

```
GET /api/system/role/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常查询**✅️

```
GET /api/system/role/1
```

**用例2: 不存在的角色**✅️

```
GET /api/system/role/99999
```

---

### 6.2 角色列表（全部，不分页）

```
GET /api/system/role/all
Authorization: Bearer {accessToken}
```

**用例1: 正常获取** -- 直接请求，用于下拉选择✅️

---

### 6.3 更新角色权限（分配菜单和接口）

```
PUT /api/system/role/:id/permissions
Content-Type: application/json
Authorization: Bearer {accessToken}
```

**用例1: 同时分配菜单和接口**✅️

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [1, 2, 3, 20, 21, 22],
    "apiIds": [1, 2, 3]
}
```

**用例2: 只分配菜单**✅️

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [1, 2, 20, 21]
}
```

**用例3: 只分配接口**✅️

```json
// PUT /api/system/role/2/permissions
{
    "apiIds": [1, 2]
}
```

**用例4: 清空所有权限**✅️

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [],
    "apiIds": []
}
```

**用例5: 包含不存在的菜单ID**✅️

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [1, 2, 99999],
    "apiIds": []
}
```

**用例6: 只分配叶子菜单（测试祖先链补全）✅️**

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [20, 21, 22],
    "apiIds": []
}
```
> 预期: 后端自动补全父菜单 2（用户管理）和祖父菜单 1（系统管理）

**用例7: 分配跨模块菜单（测试祖先链补全）✅️**

```json
// PUT /api/system/role/2/permissions
{
    "menuIds": [2, 20, 8, 42],
    "apiIds": []
}
```
> 预期: 自动补全父菜单 1（系统管理）和 7（日志管理）

**用例8: 更新不存在的角色**✅️
```json
// PUT /api/system/role/99999/permissions
{
    "menuIds": [1],
    "apiIds": [1]
}
```

---

## 七、文件管理

> 测试顺序: 上传文件 -> 文件列表 -> 删除文件

### 7.1 上传文件

```
POST /api/system/file/upload
Content-Type: multipart/form-data
Authorization: Bearer {accessToken}
```

**用例1: 上传图片** -- Postman 中 Body 选 form-data, key 为 `file`, 类型选 File, 选择一张 jpg/png 图片

**用例2: 上传文档** -- 选择一个 pdf/docx 文件

**用例3: 不选文件直接提交** -- Body 中不添加 file 字段

**用例4: 上传大文件** -- 选择一个超过服务端限制大小的文件（测试上传限制）

---

### 7.2 文件列表

```
GET /api/system/file
Authorization: Bearer {accessToken}
```

**用例1: 默认分页**
```
GET /api/system/file
```

**用例2: 指定分页**
```
GET /api/system/file?page=1&pageSize=5
```

**用例3: 按文件名搜索**
```
GET /api/system/file?keyword=avatar
```

**用例4: 组合条件**
```
GET /api/system/file?page=1&pageSize=10&keyword=test
```

---

### 7.3 删除文件

```
DELETE /api/system/file/:id
Authorization: Bearer {accessToken}
```

**用例1: 正常删除**
```
DELETE /api/system/file/{之前上传的文件id}
```

**用例2: 删除不存在的文件**
```
DELETE /api/system/file/99999
```

---

## 八、日志管理

> 测试顺序: 登录日志列表 -> 操作日志列表 -> 清空登录日志 -> 清空操作日志
>
> 前面所有模块的操作已经产生了日志数据，放在最后测试，先查后清

### 8.1 登录日志列表

```
GET /api/system/log/login
Authorization: Bearer {accessToken}
```

**用例1: 默认分页**
```
GET /api/system/log/login
```

**用例2: 指定分页**
```
GET /api/system/log/login?page=1&pageSize=5
```

**用例3: 按状态筛选 -- 成功**
```
GET /api/system/log/login?status=1
```

**用例4: 按状态筛选 -- 失败**
```
GET /api/system/log/login?status=0
```

**用例5: 按时间范围**
```
GET /api/system/log/login?startTime=2026-01-01 00:00:00&endTime=2026-12-31 23:59:59
```

**用例6: 组合条件**
```
GET /api/system/log/login?page=1&pageSize=10&status=1&startTime=2026-04-01 00:00:00&endTime=2026-04-30 23:59:59
```

---

### 8.2 操作日志列表

```
GET /api/system/log/oper
Authorization: Bearer {accessToken}
```

**用例1: 默认分页**
```
GET /api/system/log/oper
```

**用例2: 指定分页**
```
GET /api/system/log/oper?page=1&pageSize=5
```

**用例3: 按状态筛选**
```
GET /api/system/log/oper?status=1
```

**用例4: 按时间范围**
```
GET /api/system/log/oper?startTime=2026-01-01 00:00:00&endTime=2026-12-31 23:59:59
```

**用例5: 组合条件**
```
GET /api/system/log/oper?page=1&pageSize=10&status=1&startTime=2026-04-01 00:00:00
```

---

### 8.3 清空登录日志

```
DELETE /api/system/log/login/clear
Authorization: Bearer {accessToken}
```

**用例1: 正常清空** -- 直接请求

**用例2: 重复清空（已无数据时再清空）** -- 再次请求确认无报错

---

### 8.4 清空操作日志

```
DELETE /api/system/log/oper/clear
Authorization: Bearer {accessToken}
```

**用例1: 正常清空** -- 直接请求
