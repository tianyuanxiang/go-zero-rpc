# AI 应用开发方向 - 项目方案设计

> 本文档对应 `ai优化.md` 的需求落地，回答四件事：
> (1) 应用价值方向选什么；(2) Python 还是 Go；(3) 架构 / 数据库 / 数据流怎么设计；(4) MVP 怎么排。

---

## 一、应用价值方向调研

BOSS 自己想到的两个方向（个人云服务器登录分析、公司服务器日志分析）问题在于：
- 登录日志虽真实，但**单一数据源**支撑不起一个完整 AI 应用，做出来像一个"防火墙告警面板"，撑不起 RAG / Agent / LLMOps 这套故事。
- 公司服务器日志涉及**合规风险**，不宜在个人项目里采。

下面给出 8 个候选方向，按"故事性 × 数据可得性 × 技术覆盖面"排序：

| # | 方向 | 数据来源 | 技术覆盖 | 故事性 | 推荐度 |
|---|------|---------|---------|--------|--------|
| 1 | **个人技术知识库 + 编码助手**（沉淀笔记/Git/Issue，做对话式问答） | 自有 Markdown / PDF / Git 仓库 | RAG 全链路 + 流式对话 | 强（"你的第二大脑"） | **MVP 必做** |
| 2 | **云服务器登录安全智能分析**（阿里云 auth.log 实时入仓 + AI 总结告警） | 自有 ECS auth.log / nginx.log | 流式日志处理 + 异常检测 + LLM 摘要 | 强（"3834 次失败登录"自带钩子） | **MVP 必做** |
| 3 | **CVE 漏洞情报智能跟踪**（拉公开 CVE + 你服务器组件清单 → 个性化告警） | NVD / GitHub Advisory + 本地 SBOM | RAG + 定时 Agent + 推送 | 强（安全方向加分） | 二期扩展 |
| 4 | **开源项目智能客服**（指定一个 Go 开源库，把 README/Issue/PR 喂进去做问答） | 公开 GitHub 数据 | RAG + 多文档融合 | 中（可做 Demo） | 二期扩展 |
| 5 | **API 文档智能助手**（喂 OpenAPI/proto，给开发者用） | 自有 .api / .proto | RAG + 结构化数据 | 中 | 备选 |
| 6 | **服务器健康巡检 Agent**（受控工具调用：执行 ssh 命令查 CPU/磁盘/进程） | 自有 ECS | Agent + Tool Calling | 强（Agent 落地范例） | 二期扩展 |
| 7 | **GitHub Issue 智能分类回复**（爬一个仓库 Issue 做自动打标 + 回复草稿） | 公开 GitHub | Agent + 函数调用 | 中 | 备选 |
| 8 | **学习笔记反刍助手**（基于 RAG 主动推送"上周学的 X 该复习了"） | 自有 Markdown | RAG + 定时 + 个性化 | 弱（更像产品玩法） | 备选 |

**MVP 锁定**：方向 1（知识库 + 对话）+ 方向 2（登录安全分析）。
两个方向**协同非常好** —— pgvector、Embedding 服务、LLM 网关在两个模块间复用，第二个模块的"AI 总结告警"本质上也走 LLM 链路；既覆盖 RAG（检索增强），又覆盖日志智能分析这一面试加分项。

后续 Agent 方向选 #3（CVE 跟踪）或 #6（健康巡检 Agent），保持技术连贯性。

---

## 二、技术选型答疑：Python 还是 Go

**结论：Python（FastAPI）做 AI 引擎 + Go（go-zero）做业务网关 / 编排，通过 gRPC 解耦。**

理由：

| 维度 | Python | Go (eino/langchaingo) |
|------|--------|-------------------------|
| 生态 | LangChain / LlamaIndex / Unstructured / sentence-transformers / Pandas 全套 | 起步阶段，Embedding/Reranker/文档解析库严重缺失 |
| 模型 SDK | OpenAI / Anthropic / 智谱 / 通义 / DeepSeek SDK 第一方 | 多数靠 HTTP 自封装，断流/超时/重试踩坑多 |
| 文档解析 | unstructured / pypdf / pdfplumber / docx2txt 成熟 | 需要桥接 Python 或裸写解析器，工作量翻倍 |
| 招聘市场 | "AI 工程师"主流栈 | 偏小众，加分项但非必备 |
| 部署 | Docker 化 + Uvicorn 即可 | 单二进制更简，但你已有 Go 网关 |
| 性能 | I/O 密集型场景（RAG、LLM 调用）GIL 不影响 | 略占优但场景没瓶颈 |

**架构要点**：Go 网关层不直接调 LLM，所有 AI 推理走 Python 微服务（FastAPI）。Go 负责鉴权、限流、会话管理、对话历史持久化、SSE/WS 转发；Python 负责 RAG 检索、Embedding、LLM 调用、文档解析、Agent 执行。

这样未来你跳槽时可以同时讲两件事：
- "我用 Go 写了高并发的 AI 网关层"
- "我用 Python 实现了完整的 RAG 流水线"

---

## 三、整体架构（分层图）

```
┌────────────────────────────────────────────────────────────────────────┐
│                            用户 / 数据源                                │
│  Web 前端 (Vue3/React)   |   阿里云 ECS auth.log   |   本地知识文档    │
└──────────┬───────────────────────┬───────────────────────┬─────────────┘
           │ HTTPS / WS            │ Filebeat              │ 上传 / CLI
           ▼                       ▼                       ▼
┌────────────────────────────────────────────────────────────────────────┐
│  接入层  Go 网关 (gateway)                                             │
│  JWT 鉴权  |  Casbin 授权  |  限流  |  WS 长连接  |  SSE 流式转发      │
└──────────┬─────────────────────────────────────────────────────────────┘
           │ gRPC
           ▼
┌─────────────────────────┬──────────────────────────────────────────────┐
│  业务 RPC 层 (Go)        │                                              │
│  ┌────────────┐         ┌───────────────┐                              │
│  │ sys-rpc    │  现有    │ ai-rpc  新增  │  会话/知识库/日志事件编排  │
│  │ 用户/权限/ │         │ Go            │  调用 Python AI 引擎         │
│  │ 日志/文件  │         │               │  组装上下文、写入对话历史   │
│  └────────────┘         └──────┬────────┘                              │
└────────────────────────────────┼─────────────────────────────────────-─┘
                                 │ gRPC (HTTP/2 流式) 或 HTTP+SSE
                                 ▼
┌────────────────────────────────────────────────────────────────────────┐
│  AI 引擎层  Python FastAPI (ai-engine 新增)                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐│
│  │文档解析  │ │分块策略  │ │Embedding │ │向量检索  │ │LLM 路由      ││
│  │unstr-    │ │固定/语义/│ │bge-m3 /  │ │pgvector  │ │OpenAI/智谱/  ││
│  │uctured   │ │标题切分  │ │OpenAI    │ │+ rerank  │ │DeepSeek/本地 ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────────┘│
│  ┌──────────────────────────┐  ┌────────────────────────────────────┐ │
│  │ RAG Pipeline (LangChain) │  │ Agent (受控 Tool Calling)          │ │
│  └──────────────────────────┘  └────────────────────────────────────┘ │
└──────────┬─────────────────────────────────────────────────────────────┘
           │
           ▼
┌────────────────────────────────────────────────────────────────────────┐
│  存储层                                                                │
│  PostgreSQL14 + pgvector  |  Redis (会话缓存/限流)  |  MinIO (原始文件)│
└────────────────────────────────────────────────────────────────────────┘

二期扩展（可选，技术路线预留）：
  Kafka (日志总线) → Flink (实时清洗+异常检测+生成 embedding) → ES (全文) + Doris (OLAP) + pgvector (向量)
  Langfuse / Phoenix (LLMOps 观测)        Qdrant / Milvus (独立向量库)
```

**与现有项目的关系**：
- `gateway/` 增加 AI 相关路由 (`/api/ai/*`)，复用现有 JWT、Casbin、WS。
- `service/sys/rpc/` **不动**。
- 新增 `service/ai/rpc/` （Go gRPC，编排层）。
- 新增 `ai-engine/` （Python FastAPI，独立仓库或 monorepo 子目录均可，建议放 `service/ai/engine/`）。

---

## 四、MVP 模块详细设计

### 模块一：知识库 + AI 对话助手

**功能拆解**：

1. **知识库管理**
   - 创建/重命名/删除知识库（一个用户可多个）
   - 知识库可见性（私有/团队/公开）
   - 关联到对话时可多选

2. **文档管理**
   - 上传：md / txt / pdf / docx / html，单文件 ≤ 50MB
   - 状态机：`uploaded → parsing → chunking → embedding → ready / failed`
   - 失败可重试，进度通过 WS 推送
   - 文件存储：MinIO，DB 只存路径

3. **检索与对话**
   - 多轮对话，会话隔离
   - 流式输出（SSE，前端打字机效果）
   - Top-K 检索 + 可选 rerank
   - 引用溯源：回答里标 `[1][2]`，下方列出来源 chunk
   - 历史会话列表 + 重命名 + 删除
   - 支持"指定知识库范围"和"全局检索"两种模式

4. **可扩展点**（先预留接口，二期实现）
   - 多模态：图片/表格抽取
   - Hybrid Search（向量 + BM25）
   - Reranker 模型可切换

### 模块二：云服务器登录安全智能分析

**功能拆解**：

1. **数据接入**
   - Filebeat 采 `/var/log/auth.log` → 通过 HTTPS POST 到网关 `/api/ai/log/ingest`
   - 入库 `ai_log_event_raw`，原始保留 30 天

2. **事件富化**（后台 worker，30s 批处理一次）
   - IP 解析地理位置（MaxMind GeoLite2 离线库，免费）
   - IP 威胁评分（AbuseIPDB API，免费额度）
   - 解析 SSH 客户端串、登录方式、用户名
   - 写入 `ai_log_event_enriched`

3. **异常检测 + 告警**
   - 规则引擎（先用规则，简单可解释）：
     - 5 分钟内同 IP 失败 ≥ 10 次 → `brute_force`
     - 高威胁分 IP 登录尝试 → `malicious_ip`
     - 罕见地理位置成功登录 → `unusual_location`
   - 命中 → `ai_log_alert`

4. **AI 智能总结**（关键 AI 价值）
   - 每日 09:00 定时（复用现有 `service/job`）：
     - 拉取过去 24h 告警 + 富化事件
     - LLM Prompt 总结："共 N 次失败登录，主要来自 X 国 Y 个 IP 段，使用工具 Z..."
     - 写入 `ai_log_alert_summary`，前端仪表盘展示
   - 阈值触发（如 1 小时内告警 > 50）：实时调 LLM 生成"事件简报"，WS 推送

5. **对话式查询**（与模块一融合）
   - 用户："今天有多少次恶意登录？"
   - Agent → 工具调用 `query_log_alerts(date=today)` → 返回结构化数据 → LLM 组装回答

---

## 五、数据库设计初稿

沿用现有风格：PG14、表名 `ai_` 前缀、`BIGSERIAL` 主键、`created_at/updated_at/deleted_at` 三件套、`snake_case` 字段。

**前置**：`CREATE EXTENSION IF NOT EXISTS vector;`

```sql
-- ============================================================
-- 模块一：知识库 + 对话
-- ============================================================

-- 知识库
CREATE TABLE ai_knowledge_base (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL,                       -- 关联 sys_user.id
  name         VARCHAR(128) NOT NULL,
  description  VARCHAR(512) DEFAULT '',
  visibility   VARCHAR(16)  NOT NULL DEFAULT 'private', -- private/team/public
  embed_model  VARCHAR(64)  NOT NULL DEFAULT 'bge-m3', -- 该库使用的 embedding 模型
  doc_count    INTEGER      NOT NULL DEFAULT 0,
  chunk_count  INTEGER      NOT NULL DEFAULT 0,
  created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at   TIMESTAMP    NULL
);
CREATE INDEX ai_kb_idx_user ON ai_knowledge_base (user_id) WHERE deleted_at IS NULL;

-- 文档
CREATE TABLE ai_document (
  id             BIGSERIAL PRIMARY KEY,
  kb_id          BIGINT NOT NULL,
  user_id        BIGINT NOT NULL,
  file_name      VARCHAR(256) NOT NULL,
  file_type      VARCHAR(16)  NOT NULL,                 -- md/pdf/docx/txt/html
  file_size      BIGINT       NOT NULL,
  storage_path   VARCHAR(512) NOT NULL,                 -- MinIO 路径
  status         VARCHAR(32)  NOT NULL DEFAULT 'uploaded', -- uploaded/parsing/chunking/embedding/ready/failed
  error_msg      VARCHAR(1024) DEFAULT '',
  chunk_count    INTEGER      NOT NULL DEFAULT 0,
  metadata       JSONB        NOT NULL DEFAULT '{}',    -- 标题/作者/页数等
  created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at     TIMESTAMP    NULL
);
CREATE INDEX ai_doc_idx_kb     ON ai_document (kb_id) WHERE deleted_at IS NULL;
CREATE INDEX ai_doc_idx_status ON ai_document (status);

-- 文档分块（含向量）
CREATE TABLE ai_document_chunk (
  id             BIGSERIAL PRIMARY KEY,
  doc_id         BIGINT NOT NULL,
  kb_id          BIGINT NOT NULL,                       -- 冗余，便于按 KB 过滤
  chunk_index    INTEGER NOT NULL,                      -- 在文档内的顺序
  content        TEXT NOT NULL,
  content_hash   VARCHAR(64) NOT NULL,                  -- 去重用 sha256
  token_count    INTEGER NOT NULL DEFAULT 0,
  embedding      VECTOR(1024),                          -- bge-m3 输出 1024 维；切模型时建议建分表
  metadata       JSONB NOT NULL DEFAULT '{}',           -- 标题/页码/位置
  created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ai_chunk_idx_doc  ON ai_document_chunk (doc_id);
CREATE INDEX ai_chunk_idx_kb   ON ai_document_chunk (kb_id);
-- 向量索引：HNSW 召回质量高，IVFFlat 占用小；MVP 用 HNSW
CREATE INDEX ai_chunk_idx_vec  ON ai_document_chunk USING hnsw (embedding vector_cosine_ops);

-- 对话会话
CREATE TABLE ai_conversation (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL,
  title        VARCHAR(256) NOT NULL DEFAULT '新会话',  -- 第一轮后用 LLM 自动命名
  kb_ids       BIGINT[] NOT NULL DEFAULT '{}',         -- 本会话挂载的知识库
  model        VARCHAR(64) NOT NULL,                    -- 使用的 LLM 模型
  system_prompt TEXT DEFAULT '',
  msg_count    INTEGER NOT NULL DEFAULT 0,
  created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at   TIMESTAMP NULL
);
CREATE INDEX ai_conv_idx_user ON ai_conversation (user_id) WHERE deleted_at IS NULL;

-- 对话消息
CREATE TABLE ai_message (
  id              BIGSERIAL PRIMARY KEY,
  conversation_id BIGINT NOT NULL,
  role            VARCHAR(16) NOT NULL,                 -- user/assistant/system/tool
  content         TEXT NOT NULL,
  references      JSONB NOT NULL DEFAULT '[]',          -- 检索到的 chunk_id 列表
  prompt_tokens   INTEGER DEFAULT 0,
  completion_tokens INTEGER DEFAULT 0,
  latency_ms      INTEGER DEFAULT 0,
  llm_call_id     BIGINT NULL,                          -- 关联 ai_llm_call_log
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ai_msg_idx_conv ON ai_message (conversation_id, id);

-- ============================================================
-- 模块二：登录安全
-- ============================================================

-- 原始日志事件
CREATE TABLE ai_log_event_raw (
  id           BIGSERIAL PRIMARY KEY,
  source       VARCHAR(32) NOT NULL,                    -- aliyun-ecs-1, aliyun-ecs-2...
  log_type     VARCHAR(16) NOT NULL,                    -- ssh/sudo/nginx
  raw_line     TEXT NOT NULL,
  occurred_at  TIMESTAMP NOT NULL,
  ingested_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  enriched     BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX ai_log_raw_idx_time     ON ai_log_event_raw (occurred_at);
CREATE INDEX ai_log_raw_idx_enriched ON ai_log_event_raw (enriched, occurred_at) WHERE enriched = FALSE;

-- 富化事件
CREATE TABLE ai_log_event_enriched (
  id            BIGSERIAL PRIMARY KEY,
  raw_id        BIGINT NOT NULL,
  source        VARCHAR(32) NOT NULL,
  occurred_at   TIMESTAMP NOT NULL,
  event_type    VARCHAR(32) NOT NULL,                   -- login_success/login_fail/sudo/...
  username      VARCHAR(64),
  src_ip        VARCHAR(64),
  src_country   VARCHAR(64),
  src_city      VARCHAR(128),
  src_isp       VARCHAR(128),
  threat_score  INTEGER DEFAULT 0,                      -- 0-100 来自 AbuseIPDB
  client_tool   VARCHAR(128),                           -- OpenSSH_8.x / libssh / hydra ...
  auth_method   VARCHAR(32),                            -- password/publickey
  metadata      JSONB NOT NULL DEFAULT '{}',
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ai_log_enr_idx_time   ON ai_log_event_enriched (occurred_at);
CREATE INDEX ai_log_enr_idx_ip     ON ai_log_event_enriched (src_ip);
CREATE INDEX ai_log_enr_idx_type   ON ai_log_event_enriched (event_type, occurred_at);

-- 告警
CREATE TABLE ai_log_alert (
  id           BIGSERIAL PRIMARY KEY,
  alert_type   VARCHAR(32) NOT NULL,                    -- brute_force/malicious_ip/unusual_location
  severity     VARCHAR(16) NOT NULL DEFAULT 'medium',   -- low/medium/high/critical
  src_ip       VARCHAR(64),
  username     VARCHAR(64),
  event_count  INTEGER NOT NULL DEFAULT 1,
  first_seen   TIMESTAMP NOT NULL,
  last_seen    TIMESTAMP NOT NULL,
  status       VARCHAR(16) NOT NULL DEFAULT 'open',     -- open/ack/closed
  ai_summary   TEXT,                                    -- LLM 生成的事件解释
  details      JSONB NOT NULL DEFAULT '{}',
  created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ai_alert_idx_status ON ai_log_alert (status, created_at);
CREATE INDEX ai_alert_idx_type   ON ai_log_alert (alert_type, created_at);

-- 每日汇总
CREATE TABLE ai_log_alert_summary (
  id            BIGSERIAL PRIMARY KEY,
  summary_date  DATE NOT NULL,
  total_events  INTEGER NOT NULL,
  fail_count    INTEGER NOT NULL,
  alert_count   INTEGER NOT NULL,
  top_ips       JSONB NOT NULL,                         -- [{ip, country, count}]
  ai_report     TEXT NOT NULL,                          -- LLM 总结的自然语言报告
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX ai_alert_sum_uk_date ON ai_log_alert_summary (summary_date);

-- ============================================================
-- 通用：LLM 调用日志（LLMOps 雏形）/ 模型配置
-- ============================================================

CREATE TABLE ai_llm_call_log (
  id              BIGSERIAL PRIMARY KEY,
  user_id         BIGINT,
  scene           VARCHAR(32) NOT NULL,                 -- chat/summary/agent
  conversation_id BIGINT,
  provider        VARCHAR(32) NOT NULL,                 -- openai/zhipu/deepseek/local
  model           VARCHAR(64) NOT NULL,
  prompt          TEXT NOT NULL,
  completion      TEXT,
  prompt_tokens   INTEGER DEFAULT 0,
  completion_tokens INTEGER DEFAULT 0,
  total_cost      DECIMAL(10,6) DEFAULT 0,              -- 估算成本（美元）
  latency_ms      INTEGER NOT NULL,
  status          VARCHAR(16) NOT NULL,                 -- success/error/timeout
  error_msg       VARCHAR(1024),
  metadata        JSONB NOT NULL DEFAULT '{}',
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ai_llm_log_idx_user  ON ai_llm_call_log (user_id, created_at);
CREATE INDEX ai_llm_log_idx_scene ON ai_llm_call_log (scene, created_at);

CREATE TABLE ai_model_config (
  id            BIGSERIAL PRIMARY KEY,
  provider      VARCHAR(32) NOT NULL,
  model_name    VARCHAR(64) NOT NULL,
  model_type    VARCHAR(16) NOT NULL,                   -- chat/embedding/rerank
  api_base      VARCHAR(256),
  api_key_ref   VARCHAR(64),                            -- 不直接存 key，存到环境变量名/Vault key
  enabled       BOOLEAN NOT NULL DEFAULT TRUE,
  is_default    BOOLEAN NOT NULL DEFAULT FALSE,
  config_json   JSONB NOT NULL DEFAULT '{}',
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX ai_model_uk ON ai_model_config (provider, model_name, model_type);
```

---

## 六、数据流转图

### 流程 A：知识库文档入库（写）

```
[前端上传 PDF]
     │
     ▼ multipart/form-data
[Go gateway]  鉴权 + 大小校验
     │
     ▼ gRPC UploadDocument
[ai-rpc Go]
     │ 1. 写 ai_document (status=uploaded)
     │ 2. 文件落盘 MinIO
     │ 3. 通过 Redis Stream / 直接 gRPC 投任务
     ▼
[ai-engine Python]  异步 worker
     │ status=parsing  → unstructured 解析为纯文本
     │ status=chunking → 按 800 token + 100 重叠切分
     │ status=embedding→ 批量调 bge-m3 → 向量
     │ 写入 ai_document_chunk (含向量)
     │ status=ready
     ▼
[WS 推送进度] → 前端
```

### 流程 B：对话问答（读，流式）

```
[用户提问 "我去年那篇 RAG 笔记里讲了什么"]
     │
     ▼ POST /api/ai/chat (SSE)
[Go gateway]  鉴权
     │
     ▼ gRPC StreamChat
[ai-rpc Go]
     │ 1. 加载 conversation + 最近 N 条 message
     │ 2. 调 ai-engine /retrieve  (带 kb_ids + query)
     ▼
[ai-engine] embedding(query) → pgvector top-10 → rerank → top-3
     │
     ▼ 返回 chunks
[ai-rpc Go]
     │ 3. 拼 prompt (system + history + retrieved + user)
     │ 4. 调 ai-engine /chat/stream
     ▼
[ai-engine] LLM 流式调用 → token 流
     │
     ▼ 流式逐 token 返回
[ai-rpc Go] → [Go gateway SSE] → [前端打字机效果]
     │
     │ 流结束后异步：
     ▼
[ai-rpc Go]  写 ai_message + ai_llm_call_log
```

### 流程 C：登录日志智能分析

```
[阿里云 ECS] /var/log/auth.log
     │
     ▼ Filebeat (5s 一批)
[Go gateway POST /api/ai/log/ingest]   API Key 鉴权
     │
     ▼
[ai-rpc Go]  批量写 ai_log_event_raw
     │
     │ 后台 enrich worker (30s tick)
     ▼
[ai-rpc Go]  扫 enriched=false 的 raw 行
     │  调 ai-engine /enrich
     ▼
[ai-engine]  正则解析 + GeoIP + AbuseIPDB
     │  写 ai_log_event_enriched
     │
     │ 规则引擎扫描窗口
     ▼
[ai-rpc Go]  命中规则 → 写 ai_log_alert
     │  WS 推送实时告警卡片到前端
     │
     │ 每日 09:00 (service/job 定时任务)
     ▼
[ai-rpc Go]  聚合昨日数据 → 调 ai-engine /summarize
     │  写 ai_log_alert_summary
     ▼
[前端仪表盘]  地图热力 / 告警时间线 / AI 报告
```

---

## 七、接口设计要点

新增到 `gateway/api/desc/` 下，建议拆 3 个 .api：

```
ai_kb.api      知识库 + 文档     POST/GET/PUT/DELETE /api/ai/kb/...
ai_chat.api    对话               POST /api/ai/chat (SSE), 历史会话 CRUD
ai_log.api     日志安全           POST /api/ai/log/ingest, GET 告警/汇总
```

`ai-engine` (Python FastAPI) 内部接口（仅 ai-rpc 调用，不暴露公网）：

```
POST /v1/parse           文档解析 → 纯文本
POST /v1/chunk           文本分块
POST /v1/embed           批量 embedding  (input: List[str], output: List[Vector])
POST /v1/retrieve        kb_id + query → top-k chunks
POST /v1/rerank          chunks + query → 排序后 chunks
POST /v1/chat            非流式
POST /v1/chat/stream     SSE 流式
POST /v1/enrich          原始日志行 → 富化字段
POST /v1/summarize       结构化数据 + 模板 → 自然语言总结
POST /v1/agent/run       (二期) Agent 执行
```

Go 与 Python 之间建议用 **HTTP + JSON**（不是 gRPC），原因：
- 流式用 SSE 比 gRPC streaming 跨语言麻烦少；
- Python 端开发调试快；
- 性能不是瓶颈（瓶颈在 LLM）。

---

## 八、目录结构调整建议

```
go-zero-rpc/
├── common/                      # 现有，不动
├── gateway/                     # 现有
│   └── api/desc/
│       ├── sys/...              # 现有
│       └── ai/                  # 新增
│           ├── ai_kb.api
│           ├── ai_chat.api
│           └── ai_log.api
├── service/
│   ├── sys/rpc/                 # 现有，不动
│   ├── job/                     # 现有，加一个每日总结 job
│   └── ai/                      # 新增
│       ├── rpc/                 # Go 编排层 (go-zero gRPC)
│       │   ├── ai.proto
│       │   ├── etc/ai.yaml
│       │   └── internal/...
│       └── engine/              # Python AI 引擎 (FastAPI)
│           ├── pyproject.toml
│           ├── Dockerfile
│           ├── app/
│           │   ├── main.py
│           │   ├── config.py
│           │   ├── api/         # FastAPI 路由
│           │   ├── services/    # rag/embedding/llm/agent
│           │   ├── parsers/     # 各类文档解析
│           │   ├── prompts/     # prompt 模板
│           │   └── db/          # asyncpg + pgvector 客户端
│           └── tests/
├── deploy/
│   ├── system_mysql.sql
│   ├── system_pgsql.sql
│   └── ai.sql                   # 新增，本文第五节的建表 SQL
└── docx/
    ├── ai优化.md
    └── ai优化_方案设计.md       # 本文
```

---

## 九、开发路线图（MVP，建议 4 个 Sprint，每周一个）

| Sprint | 主题 | 验收标准 |
|---|---|---|
| **S0 准备**（半周） | 环境 + 骨架 | pgvector 装好；ai-engine 起得来；ai-rpc 起得来；deploy/ai.sql 跑通；调通 OpenAI / DeepSeek / 智谱任一 LLM |
| **S1 知识库写入链路** | 文档上传→分块→embedding→入库 | 上传 PDF，状态机走完到 `ready`；pgvector 中能查到向量 |
| **S2 对话问答链路** | 检索 + 流式对话 | 前端能看到流式回答，回答末尾带引用 chunk；对话历史持久化 |
| **S3 登录日志接入** | Filebeat → 入库 → 富化 → 告警 | Filebeat 配好；告警面板能看到实时事件 + GeoIP；规则告警生效 |
| **S4 AI 总结 + 收尾** | 日报 + LLMOps 日志 + 联调 | 09:00 定时跑出昨日报告；ai_llm_call_log 完整记录；写一份 README |

**单人节奏估算**：每周 8-10 小时投入，4 周完成 MVP；如果只挑一个模块（建议先做模块一），2 周可成。

---

## 十、风险点与应对

| 风险 | 应对 |
|---|---|
| **LLM API 费用失控** | 1) 默认用 DeepSeek（便宜）；2) ai_llm_call_log 实时算账；3) 单用户每日额度限制 |
| **embedding 模型切换导致历史向量失效** | 1) `ai_knowledge_base.embed_model` 字段记录；2) 切换模型时新建 KB 而非原地改；3) 维度不同建议分表 `ai_chunk_1024` / `ai_chunk_768` |
| **PDF 解析质量差** | 1) MVP 用 unstructured + 兜底 pypdf；2) 失败的文档保留原文供人工标注；3) 二期接 OCR |
| **Filebeat 拉到敏感日志** | 1) Filebeat 端配置只采 auth.log；2) ingest 入口做字段白名单过滤；3) 个人项目仅采自己的服务器，不要碰公司机器 |
| **prompt 注入 / 越权** | 1) system prompt 写死职责；2) 用户输入截断长度；3) 检索结果做用户隔离过滤 (kb 只查自己 user_id 的) |
| **pgvector 性能不够** | MVP 量级（百万 chunk 内）pgvector 完全够；超过则切 Qdrant，DAO 层抽象好即可 |

---

## 十一、二期扩展路线（不实现，只预留接口）

1. **Kafka + Flink 替换"后台 worker"**：日志吞吐到达万级 QPS 时切换；接口形态不变。
2. **Agent 模块**：受控工具调用（查日志/查告警/查文档），用 LangGraph。
3. **AI 配置中心**：从 ai_model_config 加可视化界面，运行时切换。
4. **LLMOps 接 Langfuse**：把 ai_llm_call_log 双写一份到 Langfuse。
5. **多模态**：图片表格抽取 + CLIP embedding。
6. **CVE 智能跟踪**（方向 #3）：作为模块二的横向扩展。

---

## 附：核心问题答案速查

| 问题 | 答案 |
|---|---|
| 应用价值方向？ | **MVP**：个人知识库 RAG（方向1）+ 阿里云登录安全分析（方向2）。两者数据稳定可得、共用 pgvector + Embedding + LLM 链路、覆盖 RAG 与日志智能分析两大热门场景。 |
| Python 还是 Go？ | **Python (FastAPI) + Go (go-zero) 混合**：Python 做 AI 引擎（生态完整），Go 做网关/编排（复用现有项目），HTTP+SSE 解耦。 |
| 现有代码怎么动？ | sys-rpc 不动；新增 service/ai/rpc (Go) 和 service/ai/engine (Python)；网关加 ai 路由；deploy/ai.sql 新增建表。 |
| 大数据组件什么时候上？ | MVP 不上。预留接口位置（worker 抽象 → Kafka consumer），后期热替换。 |
