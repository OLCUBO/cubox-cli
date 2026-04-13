# cubox-cli

[License: MIT](https://opensource.org/licenses/MIT)
[Go Version](https://go.dev/)
[npm version](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

[Cubox](https://cubox.pro) 官方命令行工具 — 为人类用户和 AI Agent 而设计。在终端中管理书签、浏览收藏夹、阅读保存的内容。

[安装](#安装) · [认证](#认证) · [命令](#命令) · [AI Agent](#ai-agent-快速开始) · [示例](#示例) · [开发](#开发)

## 功能


| 类别  | 能力                                 |
| --- | ---------------------------------- |
| 收藏夹 | 列出和浏览书签文件夹                         |
| 标签  | 列出和浏览标签层级                          |
| 卡片  | 按收藏夹、标签、星标/已读/标注状态、关键词、时间范围过滤和搜索卡片 |
| RAG  | 自然语言语义搜索（基于意图的智能检索）                 |
| 详情  | 查看卡片全文（Markdown）、标注、AI 洞察（摘要 + 问答） |
| 保存  | 保存网页链接为书签                          |
| 更新  | 星标/取消星标、已读/未读、归档、移动收藏夹、添加标签        |
| 删除  | 按 ID 删除收藏卡片，支持 dry-run 预览          |
| 标注  | 列出和搜索所有卡片的高亮标注                     |


## 安装

### 前置要求

- Node.js 16+（用于 `npm`/`npx`）
- Go 1.23+（仅从源码构建时需要）

### 快速开始（人类用户）

> **AI Agent 请跳转：** [AI Agent 快速开始](#ai-agent-快速开始)

#### 安装

**方式一 — 通过 npm（推荐）：**

```bash
# 安装 CLI
npm install -g cubox-cli

# 安装 CLI Skill（AI Agent 必需）
npx skills add OLCUBO/cubox-cli -g -y
```

**方式二 — 从源码构建：**

```bash
git clone https://github.com/OLCUBO/cubox-cli.git
cd cubox-cli
make install

# 安装 CLI Skill（AI Agent 必需）
npx skills add OLCUBO/cubox-cli -g -y
```

#### 登录

```bash
cubox-cli auth login
```

CLI 会引导你完成以下步骤：

1. **选择服务器** — `cubox.pro`（中国）或 `cubox.cc`（国际版）
2. **获取 API Key** — CLI 会提示打开对应网址：
  - 中国版：[https://cubox.pro/web/settings/extensions](https://cubox.pro/web/settings/extensions)
  - 国际版：[https://cubox.cc/web/settings/extensions](https://cubox.cc/web/settings/extensions)
3. **粘贴 API 链接** — 复制完整链接（如 `https://cubox.pro/c/api/save/abcdefg`）并粘贴。CLI 会自动提取 token。

#### 开始使用

```bash
# 列出收藏夹
cubox-cli group list

# 列出标签
cubox-cli tag list

# 浏览最近的卡片
cubox-cli card list --limit 5 -o pretty

# 阅读一篇文章
cubox-cli card detail --id CARD_ID
```

### AI Agent 快速开始

> 以下步骤适用于 AI Agent 帮助用户设置 cubox-cli。

**步骤 1 — 安装**

```bash
# 安装 CLI
npm install -g cubox-cli

# 安装 CLI Skill（必需）
npx skills add OLCUBO/cubox-cli -g -y
```

**步骤 2 — 认证**

询问用户使用哪个服务器（`cubox.pro` 或 `cubox.cc`），然后引导用户打开扩展设置页面并复制 API 链接：

- 中国版：`https://cubox.pro/web/settings/extensions`
- 国际版：`https://cubox.cc/web/settings/extensions`

用户提供 API 链接或 token 后执行：

```bash
cubox-cli auth login --server cubox.pro --token TOKEN
```

`--token` 参数接受完整的 API 链接 URL 或纯 token 字符串。

**步骤 3 — 验证**

```bash
cubox-cli auth status
```

**步骤 4 — 使用**

```bash
cubox-cli group list
cubox-cli card list --limit 10
```

## 认证


| 命令                                                      | 说明                      |
| ------------------------------------------------------- | ----------------------- |
| `cubox-cli auth login`                                  | 交互式登录（选择服务器 + 输入 token） |
| `cubox-cli auth login --server cubox.pro --token TOKEN` | 非交互式登录（适用于 Agent）       |
| `cubox-cli auth status`                                 | 显示当前服务器、脱敏 token、连接测试   |
| `cubox-cli auth logout`                                 | 删除已保存的凭证                |


凭证保存在 `~/.config/cubox-cli/config.json`。

## 命令

### 输出格式

所有命令支持 `-o` / `--output` 参数：


| 参数          | 说明                   |
| ----------- | -------------------- |
| `-o json`   | 紧凑 JSON（默认，适合 Agent） |
| `-o pretty` | 格式化 JSON             |
| `-o text`   | 人类可读的文本/树形输出         |


### `cubox-cli group list`

列出所有收藏夹（文件夹）。

```bash
cubox-cli group list
cubox-cli group list -o text
```

**JSON 输出字段：** `id`, `nested_name`, `name`, `parent_id`, `uncategorized`

### `cubox-cli tag list`

列出所有标签。

```bash
cubox-cli tag list
cubox-cli tag list -o text
```

**JSON 输出字段：** `id`, `nested_name`, `name`, `parent_id`

### `cubox-cli card list`

过滤和搜索收藏卡片。支持关键词搜索（使用页码分页）和浏览模式（使用游标分页）。

```bash
cubox-cli card list [flags]
```


| 参数                  | 说明                                           |
| ------------------- | -------------------------------------------- |
| `--group ID,...`    | 按收藏夹 ID 过滤                                   |
| `--tag ID,...`      | 按标签 ID 过滤                                    |
| `--starred`         | 仅星标卡片                                        |
| `--read`            | 仅已读卡片                                        |
| `--unread`          | 仅未读卡片                                        |
| `--annotated`       | 仅有标注的卡片                                      |
| `--keyword TEXT`    | 关键词搜索                                        |
| `--start-time TIME` | 按收藏开始时间过滤（如 `2026-01-01T00:00:00:000+08:00`） |
| `--end-time TIME`   | 按收藏结束时间过滤                                    |
| `--limit N`         | 每页数量（默认 50）                                  |
| `--last-id CARD_ID` | 浏览模式的游标分页（非搜索）                               |
| `--page N`          | 搜索模式的页码（从 1 开始，配合 `--keyword` 使用）            |
| `--all`             | 自动翻页获取全部结果                                   |


**分页规则：** 使用 `--keyword` 搜索时，用 `--page` 翻页；不搜索时，用 `--last-id` 游标分页。

### `cubox-cli card detail --id ID`

获取卡片完整详情，包含文章全文（Markdown）、作者、标注和 AI 洞察（摘要 + 问答）。

```bash
cubox-cli card detail --id 7247925101516031380
cubox-cli card detail --id 7247925101516031380 -o pretty
```

使用 `-o text` 仅输出 Markdown 内容。

### `cubox-cli card rag --query TEXT`

通过 RAG（检索增强生成）进行自然语言语义搜索。与关键词搜索不同，RAG 理解查询意图，即使精确词汇不匹配也能返回概念相关的卡片。

```bash
cubox-cli card rag --query "Java实现数据库图片上传功能"
cubox-cli card rag --query "如何构建带认证的 REST API" -o pretty
```

| 参数            | 说明                   |
| ------------- | -------------------- |
| `--query TEXT` | 自然语言查询文本（必填） |

**何时用 RAG vs 关键词搜索：**

- **`card list --keyword`** — 精确词汇、已知标题、域名、短语
- **`card rag --query`** — 提问、主题探索、概念性或模糊查询

### `cubox-cli save`

保存一个或多个网页链接为书签。

```bash
cubox-cli save https://example.com
cubox-cli save https://a.com https://b.com --group GROUP_ID
cubox-cli save https://example.com --tag TAG_ID1,TAG_ID2
```

### `cubox-cli update`

更新卡片属性。

```bash
cubox-cli update --id CARD_ID [flags]
```


| 参数                     | 说明      |
| ---------------------- | ------- |
| `--star` / `--unstar`  | 星标/取消星标 |
| `--read` / `--unread`  | 已读/未读   |
| `--archive`            | 归档      |
| `--group GROUP_ID`     | 移动到收藏夹  |
| `--add-tag TAG_ID,...` | 添加标签    |


### `cubox-cli delete`

删除一个或多个收藏卡片。支持 `--dry-run` 预览将要删除的内容。

```bash
cubox-cli delete --id CARD_ID [flags]
```


| 参数            | 说明                 |
| ------------- | ------------------ |
| `--id ID,...` | 要删除的卡片 ID（逗号分隔，必填） |
| `--dry-run`   | 预览将要删除的卡片，不实际执行删除  |


**Dry Run：** 建议始终先使用 `--dry-run` 预览将要删除的卡片。删除 ≤ 3 张卡片时会预览标题和 URL；批量删除更多卡片时仅显示数量，避免逐个请求详情。

```bash
# 预览
cubox-cli delete --id 7435692934957108160,7435691601617225646 --dry-run

# 确认后执行
cubox-cli delete --id 7435692934957108160,7435691601617225646
```

### `cubox-cli mark list`

列出和搜索所有卡片的高亮标注。

```bash
cubox-cli mark list [flags]
```


| 参数                  | 说明                                                |
| ------------------- | ------------------------------------------------- |
| `--color COLOR,...` | 按颜色过滤：`Yellow`, `Green`, `Blue`, `Pink`, `Purple` |
| `--keyword TEXT`    | 搜索标注                                              |
| `--start-time TIME` | 按开始时间过滤                                           |
| `--end-time TIME`   | 按结束时间过滤                                           |
| `--limit N`         | 每页数量（默认 50）                                       |
| `--last-id ID`      | 游标分页                                              |
| `--all`             | 自动翻页获取全部结果                                        |


## 示例

### 搜索文章（关键词）

```bash
cubox-cli card list --keyword "机器学习" --page 1 -o pretty
```

### 语义搜索（RAG）

```bash
cubox-cli card rag --query "如何用 Java 实现数据库图片上传并在前端展示" -o pretty
```

### 浏览特定收藏夹的卡片

```bash
cubox-cli group list -o text
cubox-cli card list --group 7230156249357091393 --limit 10
```

### 阅读文章并查看 AI 洞察

```bash
cubox-cli card detail --id 7247925101516031380 -o pretty
```

### 保存链接并加星标

```bash
cubox-cli save https://example.com
cubox-cli update --id CARD_ID --star --read
```

### 删除卡片（带 dry-run）

```bash
cubox-cli delete --id 7435692934957108160 --dry-run -o pretty
cubox-cli delete --id 7435692934957108160
```

### 导出所有标注

```bash
cubox-cli mark list --all -o pretty
```

### 游标分页（浏览模式）

```bash
cubox-cli card list --limit 5
# 使用最后一张卡片的 ID 获取下一页
cubox-cli card list --limit 5 --last-id 7433152100604841820
```

### 搜索分页

```bash
cubox-cli card list --keyword "AI" --page 1
cubox-cli card list --keyword "AI" --page 2
```

## 开发

### 从源码构建

```bash
git clone https://github.com/OLCUBO/cubox-cli.git
cd cubox-cli
make build        # 为当前平台构建
make build-all    # 交叉编译所有平台
make release      # 创建发布包
```

### 项目结构

```
cubox-cli/
  main.go                 # 入口
  cmd/                    # cobra 命令
    root.go               # 根命令，--output 参数
    auth.go               # auth login/status/logout
    group.go              # group list
    tag.go                # tag list
    card.go               # card list, card detail
    save.go               # save 保存链接
    update.go             # update 更新卡片
    delete.go             # delete 删除卡片（支持 dry-run）
    mark.go               # mark 标注列表
    version.go            # version
  internal/
    client/               # HTTP 客户端 + API 类型
    config/               # 配置文件管理
  scripts/                # npm 分发包装
  skills/cubox/           # AI Agent Skill
  .github/workflows/      # CI/CD
```

## 许可证

[MIT](LICENSE)