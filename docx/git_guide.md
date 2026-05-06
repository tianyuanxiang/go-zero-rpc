# Git 日常开发完整指�?

本文档涵�?`.gitignore` 语法详解、日常开发工作流、分支管理、敏感信息保护等内容，适用于本项目（go-zero-admin-api）的日常协作开发�?

---

## 目录

- [一�?gitignore 文件语法详解](#一gitignore-文件语法详解)
  - [1.1 基本规则](#11-基本规则)
  - [1.2 通配符语法](#12-通配符语�?
  - [1.3 完整示例](#13-完整示例)
  - [1.4 常见问题](#14-常见问题)
- [二、敏感信息保护](#二敏感信息保�?
  - [2.1 本项目的敏感文件](#21-本项目的敏感文件)
  - [2.2 新增敏感文件时的操作](#22-新增敏感文件时的操作)
  - [2.3 已经�?Git 跟踪的敏感文件如何移除](#23-已经�?git-跟踪的敏感文件如何移�?
- [三、Git 基础操作](#三git-基础操作)
  - [3.1 查看状态](#31-查看状�?
  - [3.2 添加文件到暂存区](#32-添加文件到暂存区)
  - [3.3 提交代码](#33-提交代码)
  - [3.4 查看提交历史](#34-查看提交历史)
  - [3.5 查看文件差异](#35-查看文件差异)
  - [3.6 撤销与回退](#36-撤销与回退)
- [四、分支管理](#四分支管�?
  - [4.1 分支基本操作](#41-分支基本操作)
  - [4.2 本项目的分支策略](#42-本项目的分支策略)
- [五、日常开发工作流（dev 分支）](#五日常开发工作流dev-分支)
  - [5.1 标准开发流程](#51-标准开发流�?
  - [5.2 提交信息规范](#52-提交信息规范)
- [六、合�?dev �?master（发布流程）](#六合�?dev-�?master发布流程)
  - [6.1 方式一：命令行合并](#61-方式一命令行合�?
  - [6.2 方式二：GitHub Pull Request](#62-方式二github-pull-request)
- [七、远程仓库操作](#七远程仓库操�?
- [八、冲突解决](#八冲突解�?
- [九、标签管理（版本发布）](#九标签管理版本发�?
- [十、实用技巧](#十实用技�?

---

## 一�?gitignore 文件语法详解

`.gitignore` 文件告诉 Git 哪些文件或文件夹不需要被版本控制跟踪。每一行写一条规则�?

### 1.1 基本规则

| 语法 | 含义 | 示例 |
|------|------|------|
| `文件名` | 排除所有同名文件（不论在哪个目录下�?| `config.yaml` 会排除项目中所有叫 config.yaml 的文�?|
| `目录�?` | 排除整个目录（末尾加 `/`�?| `uploads/` 会排�?uploads 目录及其所有内�?|
| `/文件名` | 仅排除项目根目录下的该文�?| `/config.yaml` 只排除根目录�?config.yaml，不影响子目录的 |
| `路径/文件名` | 排除指定路径下的文件 | `etc/config.yaml` 只排�?etc 目录下的 config.yaml |
| `!文件名` | 取反，表�?不排�?（强制包含） | `!important.yaml` 即使上面的规则排除了 *.yaml，这个文件也会被保留 |
| `# 注释` | �?`#` 开头的行是注释，Git 会忽�?| `# 这是注释` |
| 空行 | 空行会被忽略，用于分组增加可读�?| |

### 1.2 通配符语�?

| 通配�?| 含义 | 示例 | 匹配结果 |
|--------|------|------|----------|
| `*` | 匹配任意多个字符（不�?`/`�?| `*.log` | app.log、error.log、debug.log |
| `**` | 匹配任意多层目录 | `**/logs` | logs、src/logs、src/app/logs |
| `?` | 匹配单个任意字符 | `file?.txt` | file1.txt、fileA.txt，但不匹�?file12.txt |
| `[abc]` | 匹配方括号内的任意一个字�?| `file[0-9].txt` | file0.txt �?file9.txt |
| `[!abc]` | 匹配不在方括号内的任意一个字�?| `file[!0-9].txt` | fileA.txt，但不匹�?file1.txt |

**详细示例�?*

```gitignore
# --- 星号 * ：匹配任意多个字符（不跨目录�?---

*.log                  # 排除所�?.log 文件：app.log, error.log, tmp/debug.log
*.exe                  # 排除所�?.exe 文件：main.exe, tools/build.exe
temp_*                 # 排除�?temp_ 开头的文件：temp_data.csv, temp_cache

# --- 双星�?** ：匹配任意多层目�?---

**/build               # 排除任意深度�?build 目录：build/、src/build/、a/b/build/
logs/**                # 排除 logs 目录下的所有内容（�?logs 目录本身会保留）
**/test/**/*.tmp       # 排除任意 test 目录下任意深度的 .tmp 文件

# --- 问号 ? ：匹配单个字�?---

config?.yaml           # 匹配 config1.yaml、configA.yaml，不匹配 config12.yaml
?.txt                  # 匹配 a.txt�?.txt，不匹配 ab.txt

# --- 方括�?[] ：匹配指定范围的单个字符 ---

log[0-9].txt           # 匹配 log0.txt �?log9.txt
file[abc].go           # 匹配 filea.go、fileb.go、filec.go
data[!0-9].csv         # 匹配 dataA.csv，不匹配 data1.csv
```

### 1.3 完整示例

下面是一�?Go 项目�?`.gitignore` 完整示例，每行都有注释：

```gitignore
# =============================================
# Go 项目 .gitignore 示例
# =============================================

# ---------- 编译产物 ----------
# Go 编译生成的二进制文件
*.exe
*.exe~
*.dll
*.so
*.dylib

# 测试编译产物
*.test

# 覆盖率报�?
*.out

# ---------- 依赖管理 ----------
# Go workspace 文件（多模块开发时生成�?
go.work
go.work.sum

# vendor 目录（如果不想提交依赖包，取消下面的注释�?
# vendor/

# ---------- IDE / 编辑�?----------
# JetBrains (GoLand) 配置
.idea/

# VS Code 配置
.vscode/

# Vim 临时文件
*.swp
*.swo
*~

# ---------- 操作系统文件 ----------
# macOS
.DS_Store

# Windows
Thumbs.db
desktop.ini

# ---------- 敏感配置文件（最重要！） ----------
# 包含数据库密码、Redis密码、JWT密钥等真实配�?
etc/config.yaml

# 环境变量文件
.env
.env.local
.env.production

# 密钥文件
*.key
*.pem
*.p12

# ---------- 运行时产�?----------
# 上传文件目录
uploads/

# 日志文件
*.log
logs/

# 临时文件
tmp/

# ---------- Claude Code ----------
.claude/
```

### 1.4 常见问题

**Q1�?gitignore 中的规则不生效怎么办？**

如果文件已经�?Git 跟踪过（之前提交过），`.gitignore` 无法直接生效。需要先�?Git 缓存中移除：

```bash
# �?Git 跟踪中移除（不删除本地文件）
git rm --cached etc/config.yaml

# 提交这次移除操作
git commit -m "chore: remove tracked sensitive file etc/config.yaml"
```

**Q2：如何排除整个目录但保留其中某个文件�?*

```gitignore
# 排除 logs 目录下的所有内�?
logs/*

# 但保�?logs 目录下的 .gitkeep 文件（用于保持空目录�?
!logs/.gitkeep
```

**Q3：如何只排除根目录的文件，不影响子目录？**

```gitignore
# 只排除根目录�?config.yaml
/config.yaml

# 以下路径不受影响�?
# src/config.yaml （不会被排除�?
# pkg/config.yaml （不会被排除�?
```

**Q4：如何检查某个文件是否被 .gitignore 排除了？**

```bash
# 检查单个文件是否被忽略
git check-ignore -v etc/config.yaml
# 输出示例�?gitignore:30:etc/config.yaml    etc/config.yaml
# 表示 .gitignore �?30 行的规则排除了这个文�?

# 如果没有输出，说明该文件没有被任何规则排�?
```

**Q5：如何查看当前所有被忽略的文件？**

```bash
# 列出所有被忽略的文�?
git status --ignored

# 只看被忽略的文件列表（更简洁）
git status --ignored --short
```

---

## 二、敏感信息保�?

### 2.1 本项目的敏感文件

| 文件 | 包含内容 | 是否已排�?|
|------|----------|-----------|
| `etc/config.yaml` | MySQL密码、Redis密码、JWT密钥、服务器IP | 已排�?|
| `etc/config.yaml.example` | 配置模板（占位符，无真实密码�?| 已提交到仓库 |

新成�?clone 项目后的首次配置�?

```bash
# 1. 克隆项目
git clone https://github.com/tianyuanxiang/go-zero-admin-api.git
cd go-zero-admin-api

# 2. 复制配置模板
cp etc/config.yaml.example etc/config.yaml

# 3. 编辑 config.yaml，填入真实的数据库密码、Redis密码�?
# 用你喜欢的编辑器打开�?
notepad etc/config.yaml       # Windows
vim etc/config.yaml           # Linux/Mac
code etc/config.yaml          # VS Code
```

### 2.2 新增敏感文件时的操作

假设将来项目新增�?`.env` 文件�?`certs/server.key` 密钥文件�?

```bash
# 第一步：�?.gitignore 中添加规�?
echo ".env" >> .gitignore
echo "certs/*.key" >> .gitignore

# 第二步：提交 .gitignore 的变�?
git add .gitignore
git commit -m "chore: add .env and key files to gitignore"
git push origin dev
```

### 2.3 已经�?Git 跟踪的敏感文件如何移�?

如果你不小心把敏感文件提交到�?Git，需要这样处理：

```bash
# 场景：etc/config.yaml 已经被提交过了，现在要从 Git 中移�?

# 第一步：�?.gitignore 中添加排除规�?
# 确保 .gitignore 中已�?etc/config.yaml

# 第二步：�?Git 跟踪中移除（--cached 表示只删�?Git 记录，不删除本地文件�?
git rm --cached etc/config.yaml

# 第三步：提交
git add .gitignore
git commit -m "chore: remove sensitive config from git tracking"
git push origin dev

# 此时�?
# - 本地 etc/config.yaml 仍然存在，不影响你本地运�?
# - 远程仓库�?etc/config.yaml 被删�?
# - 后续提交不会再包含这个文�?
```

**注意�?* 如果敏感文件已经推送到远程仓库，历史提交中仍然包含该文件的内容。如果泄露了真实密码，最安全的做法是�?*立即修改密码**�?

### 2.4 提交前检查敏感文�?

养成好习惯，每次提交前检查是否误添加了敏感文件：

```bash
# 查看暂存区中的文件列�?
git diff --cached --name-only

# �?grep 过滤可能的敏感文�?
git diff --cached --name-only | grep -E "\.(yaml|yml|env|key|pem|p12)$"

# 如果发现误添加了敏感文件，从暂存区移除（不删除本地文件）
git reset HEAD etc/config.yaml
```

---

## 三、Git 基础操作

### 3.1 查看状�?

```bash
# 查看工作区和暂存区的状态（最常用的命令）
git status

# 输出示例�?
# On branch dev
# Changes not staged for commit:    <-- 已修改但未添加到暂存�?
#   modified:   internal/logic/auth/login_logic.go
#
# Untracked files:                  <-- 新文件，还没�?Git 跟踪
#   internal/handler/new_handler.go

# 简洁模式（更紧凑的输出�?
git status --short
# 输出示例�?
#  M internal/logic/auth/login_logic.go    （M 表示已修改）
# ?? internal/handler/new_handler.go       �?? 表示未跟踪的新文件）
```

**状态标记含义：**

| 标记 | 含义 |
|------|------|
| `M` | Modified，文件已修改 |
| `A` | Added，新文件已添加到暂存�?|
| `D` | Deleted，文件已删除 |
| `R` | Renamed，文件已重命�?|
| `??` | Untracked，新文件，未�?Git 跟踪 |

### 3.2 添加文件到暂存区

暂存区是提交前的"准备�?，只有添加到暂存区的修改才会被提交�?

```bash
# 添加单个文件
git add internal/logic/auth/login_logic.go

# 添加多个文件（空格分隔）
git add internal/logic/auth/login_logic.go internal/handler/auth/login_handler.go

# 添加某个目录下的所有修�?
git add internal/logic/

# 添加所有修改（谨慎使用，可能会添加敏感文件�?
git add .

# 推荐做法：先查看状态，再逐个添加
git status                    # 先看有哪些修�?
git add 文件1 文件2 文件3      # 只添加需要的文件
```

**从暂存区移除（不删除本地文件）：**

```bash
# 取消暂存单个文件
git reset HEAD internal/logic/auth/login_logic.go

# 取消暂存所有文�?
git reset HEAD
```

### 3.3 提交代码

```bash
# 标准提交（单行提交信息）
git commit -m "feat: add user login rate limit"

# 多行提交信息（详细描述）
git commit -m "feat: add user login rate limit

- Add rate limiter middleware for login endpoint
- Limit to 5 attempts per minute per IP
- Return 429 status code when exceeded"

# 查看刚才的提�?
git log -1
```

### 3.4 查看提交历史

```bash
# 查看完整提交历史
git log

# 查看最�?5 条提交（常用�?
git log -5

# 单行简洁模式（快速浏览历史）
git log --oneline -10
# 输出示例�?
# 7e7be37 refactor: restructure project as go-zero-admin system
# 3605b6f init: plating platform project

# 查看某个文件的修改历�?
git log --oneline internal/logic/auth/login_logic.go

# 图形化查看分支合并历�?
git log --oneline --graph --all
# 输出示例�?
# * 7e7be37 (HEAD -> dev, origin/dev, origin/master, master) refactor: ...
# * 3605b6f init: plating platform project

# 查看某次提交的详细内�?
git show 7e7be37
```

### 3.5 查看文件差异

```bash
# 查看工作区中已修改但未暂存的差异
git diff
# 输出示例�?
# -    oldCode := "before"    （红色，表示删除的行�?
# +    newCode := "after"     （绿色，表示新增的行�?

# 查看某个文件的差�?
git diff internal/logic/auth/login_logic.go

# 查看已暂存（git add 之后）的差异
git diff --cached

# 查看两个分支之间的差�?
git diff master..dev

# 只看哪些文件有变化（不看具体内容�?
git diff --name-only master..dev
```

### 3.6 撤销与回退

```bash
# --- 场景1：修改了文件但还�?git add，想恢复原样 ---
git checkout -- internal/logic/auth/login_logic.go
# 注意：这会丢失所有未暂存的修改，无法恢复

# --- 场景2：已�?git add 了，想取消暂存但保留修改 ---
git reset HEAD internal/logic/auth/login_logic.go
# 文件回到"已修改但未暂�?状态，修改内容还在

# --- 场景3：已�?git commit 了，想撤销这次提交但保留代码修�?---
git reset --soft HEAD~1
# HEAD~1 表示回退1个提�?
# --soft 表示代码修改保留在暂存区，可以重新提�?

# --- 场景4：已�?git commit 了，想完全撤销（代码也回退�?---
git reset --hard HEAD~1
# 警告�?-hard 会丢失代码修改，无法恢复，谨慎使用！

# --- 场景5：已�?push 到远程了，想撤销 ---
# 不要�?reset，用 revert 创建一�?反向提交"
git revert HEAD
# 这会创建一个新提交，内容是撤销上一次提交的修改
# 然后正常 push 即可
git push origin dev
```

---

## 四、分支管�?

### 4.1 分支基本操作

```bash
# 查看所有本地分支（* 表示当前分支�?
git branch
# 输出示例�?
# * dev
#   master

# 查看所有分支（包括远程分支�?
git branch -a
# 输出示例�?
# * dev
#   master
#   remotes/origin/dev
#   remotes/origin/master

# 创建新分�?
git branch feature/user-avatar

# 切换到已有分�?
git checkout dev
# 或者用新版命令（功能相同）
git switch dev

# 创建新分支并切换过去（一步到位）
git checkout -b feature/user-avatar
# �?
git switch -c feature/user-avatar

# 删除本地分支（确保不在该分支上）
git branch -d feature/user-avatar

# 删除远程分支
git push origin --delete feature/user-avatar

# 重命名当前分�?
git branch -m new-branch-name
```

### 4.2 本项目的分支策略

```
master  ────●────────────────●──────────────  （稳定版本，用于部署�?
             \              /
dev     ──────●────●────●──●────●────●──────  （日常开发）
```

| 分支 | 用�?| 谁来操作 |
|------|------|---------|
| `master` | 稳定版本，用于生产部�?| 只通过 dev 合并过来 |
| `dev` | 日常开�?| 开发人员直接提�?|
| `feature/xxx`（可选） | 开发大功能时从 dev 拉出 | 开发完成后合并�?dev |

---

## 五、日常开发工作流（dev 分支�?

### 5.1 标准开发流�?

```bash
# ============================================
# 第一步：确保�?dev 分支，并拉取最新代�?
# ============================================
git checkout dev
git pull origin dev

# ============================================
# 第二步：写代�?..
# （在你的编辑器中修改文件�?
# ============================================

# ============================================
# 第三步：查看修改了哪些文�?
# ============================================
git status

# ============================================
# 第四步：检查是否有敏感文件要被提交
# ============================================
git diff --cached --name-only | grep -E "\.(yaml|env|key|pem)$"
# 如果有输出，说明可能要提交敏感文件，请检�?

# ============================================
# 第五步：添加需要提交的文件（推荐逐个添加�?
# ============================================
git add internal/logic/auth/login_logic.go
git add internal/handler/auth/login_handler.go

# ============================================
# 第六步：再次确认暂存的内�?
# ============================================
git diff --cached --stat
# 输出示例�?
#  internal/logic/auth/login_logic.go   | 15 +++++++++------
#  internal/handler/auth/login_handler.go |  8 +++++---
#  2 files changed, 14 insertions(+), 9 deletions(-)

# ============================================
# 第七步：提交
# ============================================
git commit -m "feat: add login rate limit middleware"

# ============================================
# 第八步：推送到远程 dev 分支
# ============================================
git push origin dev
```

### 5.2 提交信息规范

提交信息格式：`<类型>: <简短描�?`

| 类型 | 含义 | 示例 |
|------|------|------|
| `feat` | 新功�?| `feat: add user avatar upload` |
| `fix` | 修复 Bug | `fix: correct login token expiration time` |
| `refactor` | 重构代码（不改变功能�?| `refactor: extract common validation logic` |
| `docs` | 文档变更 | `docs: update API documentation` |
| `style` | 代码格式调整（不影响逻辑�?| `style: format code with gofmt` |
| `test` | 添加或修改测�?| `test: add unit test for login logic` |
| `chore` | 构建、依赖、配置等杂项 | `chore: update go.mod dependencies` |
| `perf` | 性能优化 | `perf: optimize user list query with index` |

---

## 六、合�?dev �?master（发布流程）

### 6.1 方式一：命令行合并

```bash
# ============================================
# 第一步：确保 dev 分支代码已全部提交并推�?
# ============================================
git checkout dev
git status                  # 确保没有未提交的修改
git push origin dev         # 确保远程是最新的

# ============================================
# 第二步：切换�?master 分支，拉取最新代�?
# ============================================
git checkout master
git pull origin master

# ============================================
# 第三步：�?dev 合并�?master
# ============================================
git merge dev
# 如果没有冲突，会自动完成合并
# 如果有冲突，需要手动解决（见第八节"冲突解决"�?

# ============================================
# 第四步：推�?master 到远�?
# ============================================
git push origin master

# ============================================
# 第五步：切回 dev 分支继续开�?
# ============================================
git checkout dev
```

### 6.2 方式二：GitHub Pull Request（推荐）

Pull Request (PR) 可以在合并前进行代码审查，适合团队协作�?

```bash
# 第一步：确保 dev 分支代码已推送到远程
git checkout dev
git push origin dev
```

然后�?GitHub 网页上操作：

1. 打开仓库页面 `https://github.com/tianyuanxiang/go-zero-admin-api`
2. 点击 **"Pull requests"** 标签�?
3. 点击 **"New pull request"** 按钮
4. 设置�?
   - **base**: `master` （合并目标）
   - **compare**: `dev` （要合并的来源）
5. 填写标题和描述，点击 **"Create pull request"**
6. 审查通过后，点击 **"Merge pull request"**
7. 合并完成后，回到本地更新 master�?

```bash
git checkout master
git pull origin master
git checkout dev
```

---

## 七、远程仓库操�?

```bash
# 查看远程仓库信息
git remote -v
# 输出示例�?
# origin  https://github.com/tianyuanxiang/go-zero-admin-api.git (fetch)
# origin  https://github.com/tianyuanxiang/go-zero-admin-api.git (push)

# 从远程拉取最新代码（不自动合并）
git fetch origin

# 从远程拉取并自动合并到当前分�?
git pull origin dev

# 推送当前分支到远程
git push origin dev

# 推送并设置上游跟踪（第一次推送新分支时用�?
git push -u origin dev
# 设置后，以后只需 git push 即可，不用写 origin dev

# 修改远程仓库地址
git remote set-url origin https://github.com/tianyuanxiang/new-repo.git

# 添加第二个远程仓库（比如备份到另一个平台）
git remote add backup https://gitee.com/tianyuanxiang/go-zero-admin-api.git
git push backup dev
```

---

## 八、冲突解�?

当两个分支修改了同一个文件的同一部分时，合并时会产生冲突�?

```bash
# 合并时遇到冲�?
git merge dev
# 输出�?
# CONFLICT (content): Merge conflict in internal/logic/auth/login_logic.go
# Automatic merge failed; fix conflicts and then commit the result.
```

**冲突标记解读�?*

打开有冲突的文件，你会看到类似这样的内容�?

```go
func Login(username, password string) error {
<<<<<<< HEAD
    // master 分支的代�?
    maxRetry := 3
=======
    // dev 分支的代�?
    maxRetry := 5
>>>>>>> dev
    // ...
}
```

| 标记 | 含义 |
|------|------|
| `<<<<<<< HEAD` | 当前分支（master）的代码开�?|
| `=======` | 分隔�?|
| `>>>>>>> dev` | 合并来源分支（dev）的代码结束 |

**解决步骤�?*

```bash
# 第一步：打开冲突文件，手动选择保留哪部分代�?
# 删除冲突标记�?<<<<<< ======= >>>>>>>），保留正确的代码：
#
#   maxRetry := 5
#

# 第二步：标记冲突已解�?
git add internal/logic/auth/login_logic.go

# 第三步：完成合并提交
git commit -m "merge: resolve conflict in login_logic.go"

# 第四步：推�?
git push origin master
```

**如果想放弃本次合并：**

```bash
# 在解决冲突之前，可以取消合并
git merge --abort
# 一切回到合并之前的状�?
```

---

## 九、标签管理（版本发布�?

标签用于标记重要的版本节点，比如 v1.0.0�?

```bash
# 创建轻量标签
git tag v1.0.0

# 创建带注释的标签（推荐，包含更多信息�?
git tag -a v1.0.0 -m "release: first stable version"

# 查看所有标�?
git tag

# 查看某个标签的详细信�?
git show v1.0.0

# 推送标签到远程
git push origin v1.0.0

# 推送所有标签到远程
git push origin --tags

# 删除本地标签
git tag -d v1.0.0

# 删除远程标签
git push origin --delete v1.0.0
```

**推荐的版本号规范（语义化版本 Semantic Versioning）：**

```
v主版本号.次版本号.修订�?

v1.0.0  - 首个稳定版本
v1.1.0  - 新增功能（向下兼容）
v1.1.1  - 修复 Bug
v2.0.0  - 重大变更（不向下兼容�?
```

---

## 十、实用技�?

### 10.1 暂存当前修改（临时切换分支）

场景：你正在 dev 上开发，突然需要切�?master 修复一个紧�?Bug�?

```bash
# 暂存当前所有修改（不需要提交）
git stash
# 输出：Saved working directory and index state WIP on dev: 7e7be37 ...

# 现在可以安全切换分支
git checkout master
# 修复 Bug ...
git checkout dev

# 恢复之前暂存的修�?
git stash pop
# 代码回来了，继续开�?

# 查看暂存列表
git stash list

# 丢弃暂存（不恢复�?
git stash drop
```

### 10.2 查看谁修改了某一行代�?

```bash
# 查看文件每一行最后是谁在什么时候修改的
git blame internal/logic/auth/login_logic.go

# 输出示例�?
# 7e7be37 (tianyuanxiang 2026-04-27 10:00:00 +0800  1) package auth
# 7e7be37 (tianyuanxiang 2026-04-27 10:00:00 +0800  2)
# 3605b6f (tianyuanxiang 2026-04-20 09:30:00 +0800  3) import (

# 只看某几�?
git blame -L 10,20 internal/logic/auth/login_logic.go
```

### 10.3 搜索提交历史

```bash
# 搜索提交信息中包含关键词的提�?
git log --grep="login" --oneline

# 搜索代码中包含某个字符串的提交（谁添�?删除了这行代码）
git log -S "MaxRetryCount" --oneline

# 搜索某个作者的提交
git log --author="tianyuanxiang" --oneline -10
```

### 10.4 修改最后一次提�?

```bash
# 场景：刚提交完发现提交信息写错了
git commit --amend -m "fix: correct error message in login handler"
# 注意：只能修改最后一次提交，且该提交还没�?push 到远�?

# 场景：刚提交完发现漏了一个文�?
git add internal/handler/auth/login_handler.go
git commit --amend --no-edit
# --no-edit 表示不修改提交信息，只补充文�?
```

### 10.5 快捷别名配置

```bash
# 配置常用命令的短别名
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.lg "log --oneline --graph --all"

# 配置后可以这样使用：
git st          # 等同�?git status
git co dev      # 等同�?git checkout dev
git br          # 等同�?git branch
git ci -m "msg" # 等同�?git commit -m "msg"
git lg          # 等同�?git log --oneline --graph --all
```

### 10.6 .gitkeep 保持空目�?

Git 不会跟踪空目录。如果你想保留一个空目录（比�?`uploads/`），在里面放一个空文件�?

```bash
# 在空目录中创�?.gitkeep 文件
touch uploads/.gitkeep

# .gitignore 中保�?.gitkeep
# uploads/*
# !uploads/.gitkeep
```

---

## 附录：本项目常用命令速查�?

| 操作 | 命令 |
|------|------|
| 查看状�?| `git status` |
| 添加文件 | `git add 文件名` |
| 提交代码 | `git commit -m "类型: 描述"` |
| 推�?dev | `git push origin dev` |
| 拉取最�?| `git pull origin dev` |
| 切换分支 | `git checkout master` |
| 合并分支 | `git merge dev` |
| 查看历史 | `git log --oneline -10` |
| 查看差异 | `git diff` |
| 暂存修改 | `git stash` |
| 恢复暂存 | `git stash pop` |
| 检查忽�?| `git check-ignore -v 文件名` |
| 取消暂存 | `git reset HEAD 文件名` |

```
徐州：戏马台、博物馆、云龙湖、云龙山；烧烤不清楚，人�?
沛县：庆景冷�?老店，人�?，徐庄矿老四冷面(老店，不在县城，比较偏，但好�?，铁西老二羊杂冷面(味道可以，不出名)，老味道冷�?歌风佳苑北面的店)、沛风羊肉，鹿串烧烤、广阔烧烤、早上可以喝一杯撒汤（温州商贸城庆云撒汤）

```

