# cubox-cli

[License: MIT](https://opensource.org/licenses/MIT)
[Go Version](https://go.dev/)
[npm version](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

Cubox 官方 CLI。收藏、搜索、阅读，并借助 AI 使用你读过的内容。让你的个人阅读记忆，真正被用起来。

[安装](#安装) · [认证](#认证) · [命令](#命令) · [AI Agent](#ai-agent-快速开始) · [示例](#示例) · [开发](#开发)

## 功能


| 类别  | 能力                                 |
| --- | ---------------------------------- |
| 收藏夹 | 列出和浏览文件夹                         |
| 标签  | 列出和浏览标签层级；重命名、批量删除、合并标签              |
| 卡片  | 按收藏夹、标签、星标/已读/标注/归档状态、关键词、时间范围过滤和搜索卡片 |
| RAG  | 自然语言语义搜索（基于意图的智能检索）                 |
| 详情  | 查看卡片全文（Markdown）、标注、AI 洞察（摘要 + 问答） |
| 保存  | 保存网页，支持标题/描述，批量 JSON 输入            |
| 更新  | 星标/取消星标、已读/未读、移动收藏夹、管理标签              |
| 归档  | 批量归档收藏卡片，或恢复（取消归档）到指定收藏夹           |
| 删除  | 按 ID 删除收藏卡片，支持 dry-run 预览          |
| 标注  | 列出和搜索所有卡片的标注，包含高亮和高亮笔记                     |


## 安装

### 前置要求

- Node.js 16+（用于 `npm`/`npx`）
- Go 1.23+（仅从源码构建时需要）

### 快速开始（人类用户）

> **AI Agent 请跳转：** [AI Agent 快速开始](#ai-agent-快速开始)

#### 安装

**方式一：通过 npm（推荐）**

```bash
# 安装 CLI
npm install -g cubox-cli

# 安装 CLI Skill（AI Agent 必需）
npx skills add OLCUBO/cubox-cli -g -y
```

**方式二：从源码构建**

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

1. **登录 Cubox 网页版** — 根据您的 Cubox 所在地区，打开对应的扩展设置页面：
   - .pro 用户：[https://cubox.pro/web/settings/extensions](https://cubox.pro/web/settings/extensions)
   - .cc 国际用户：[https://cubox.cc/web/settings/extensions](https://cubox.cc/web/settings/extensions)
2. **启用 API 扩展** — 在扩展中心页面找到「API 扩展」，启用并复制链接（如 `https://cubox.pro/c/api/save/abcd12345`）。
3. **粘贴链接** — CLI 会自动提取 token 和服务器信息。

#### 开始使用

```bash
# 列出收藏夹
cubox-cli folder list

# 列出标签
cubox-cli tag list

# 浏览最近的卡片
cubox-cli card list --limit 5 -o pretty

# 阅读一篇文章
cubox-cli card detail --id CARD_ID
```

### AI Agent 快速开始

> 以下步骤适用于 AI Agent 帮助用户设置 cubox-cli。

**步骤 1：安装**

```bash
# 安装 CLI
npm install -g cubox-cli

# 安装 CLI Skill（必需）
npx skills add OLCUBO/cubox-cli -g -y
```

**步骤 2：认证**

引导用户获取 API 链接：

1. 请根据您的 Cubox 所在地区，登录对应的网页版：
   - .pro 用户：`https://cubox.pro/web/settings/extensions`
   - .cc 国际用户：`https://cubox.cc/web/settings/extensions`
2. 在扩展中心页面找到「API 扩展」，启用并复制链接（如 `https://cubox.pro/c/api/save/abcd12345`）。

Agent **不要**让用户把 token 粘贴到对话里，也不要把 token 拼接进命令行 argv。请指引用户自行在终端运行以下任意一种：

```bash
# 交互式（推荐给人类用户）
cubox-cli auth login

# Agent / CI 无落盘模式 —— token 仅存在于 shell 环境变量中
export CUBOX_SERVER=cubox.pro
export CUBOX_TOKEN=... # 也可以是完整的 API 链接 URL
cubox-cli folder list

# 非交互式持久登录 —— token 通过 stdin 传入，不进入 argv / ps / history
printf '%s' "$TOKEN" | cubox-cli auth login --server cubox.pro --token-stdin
```

旧式 `--token TOKEN` 参数仍然可用，但会把 token 泄漏到 shell history 和 `ps`，仅建议在受控环境使用。

**步骤 3：验证**

```bash
cubox-cli auth status
```

**步骤 4：使用**

```bash
cubox-cli folder list
cubox-cli card list --limit 10
```

## 认证


| 命令                                                                               | 说明                                                  |
| ------------------------------------------------------------------------------- | --------------------------------------------------- |
| `cubox-cli auth login`                                                          | 交互式登录（选择服务器 + 输入 token）                             |
| `printf '%s' "$TOKEN" \| cubox-cli auth login --server cubox.pro --token-stdin` | 从 stdin 读取 token 的非交互式登录（推荐 Agent 使用）                |
| `CUBOX_SERVER=cubox.pro CUBOX_TOKEN=... cubox-cli ...`                          | 纯环境变量模式，不落盘（推荐 CI / 沙箱场景）                           |
| `cubox-cli auth login --server cubox.pro --token TOKEN`                         | 旧式 argv 写法，会泄漏到 shell history / ps，尽量避免             |
| `cubox-cli auth status`                                                         | 显示当前服务器、脱敏 token、连接测试                               |
| `cubox-cli auth logout`                                                         | 删除已保存的凭证                                            |


凭证保存在 `~/.config/cubox-cli/config.json`。设置 `CUBOX_TOKEN` / `CUBOX_SERVER` 环境变量会覆盖配置文件；当配置文件不存在时，单独使用这两个环境变量即可登录。

## 命令

### 输出格式

所有命令都接受 `-o` / `--output` 参数。大多数数据类命令会按该参数输出：


| 参数          | 说明                   |
| ----------- | -------------------- |
| `-o json`   | 紧凑 JSON（默认，适合 Agent） |
| `-o pretty` | 格式化 JSON             |
| `-o text`   | 人类可读的文本/树形输出         |


注意：`save`、`update` 和部分 `auth` 成功路径目前即使选择 `-o json` 也会输出纯文本。

### `cubox-cli version`

显示当前安装的 CLI 版本。

```bash
cubox-cli version
```

### `cubox-cli folder list`

列出所有收藏夹。

```bash
cubox-cli folder list
cubox-cli folder list -o text
```

**JSON 输出字段：** `id`, `nested_name`, `name`, `parent_id`, `uncategorized`

### `cubox-cli tag list`

列出所有标签。

```bash
cubox-cli tag list
cubox-cli tag list -o text
```

**JSON 输出字段：** `id`, `nested_name`, `name`, `parent_id`

### `cubox-cli tag update`

按 ID 重命名标签。新名称仅影响叶子段，嵌套的子标签会自动跟随到新路径下。

```bash
cubox-cli tag update --id TAG_ID --new-name NEW_NAME
```

| 参数                | 说明                                       |
| ----------------- | ---------------------------------------- |
| `--id ID`         | 要重命名的标签 ID（必填）                          |
| `--new-name NAME` | 新的叶子名称（必填；不能包含 `/`）                     |

### `cubox-cli tag delete`

按 ID 批量删除标签。被删除标签关联的卡片本身保留，仅解除标签关联。

```bash
cubox-cli tag delete --id TAG_ID[,ID2,...]
```

| 参数            | 说明                              |
| ------------- | ------------------------------- |
| `--id ID,...` | 要删除的标签 ID（逗号分隔，必填）            |

### `cubox-cli tag merge`

将一个或多个源标签合并到目标标签。源标签关联的卡片会重新挂到目标标签下，随后源标签被删除。

```bash
cubox-cli tag merge --source SRC_ID[,ID2,...] --target TARGET_ID
```

| 参数                | 说明                                |
| ----------------- | --------------------------------- |
| `--source ID,...` | 要合并的源标签 ID（逗号分隔，必填）             |
| `--target ID`     | 合并到的目标标签 ID（必填）                  |

### `cubox-cli card list`

过滤和搜索收藏卡片。支持关键词搜索（使用页码分页）和浏览模式（使用游标分页）。

```bash
cubox-cli card list [flags]
```


| 参数                  | 说明                                           |
| ------------------- | -------------------------------------------- |
| `--folder ID,...`   | 按收藏夹 ID 过滤                                   |
| `--tag ID,...`      | 按标签 ID 过滤                                    |
| `--starred`         | 仅星标卡片                                        |
| `--read`            | 仅已读卡片                                        |
| `--unread`          | 仅未读卡片                                        |
| `--annotated`       | 仅有标注的卡片                                      |
| `--archived`        | 仅归档卡片（默认仅返回未归档卡片）                              |
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

保存一个或多个网页为书签。支持三种输入方式。

```bash
# 简单模式 — URL 参数
cubox-cli save https://example.com
cubox-cli save https://a.com https://b.com --folder "阅读清单"

# 单个带元数据
cubox-cli save https://example.com --title "我的页面" --desc "有趣的文章"

# 批量 JSON 输入
cubox-cli save --json '[{"url":"https://a.com","title":"标题A"},{"url":"https://b.com"}]' --tag 技术,AI/LLM
```

| 参数               | 说明                                                             |
| ---------------- | -------------------------------------------------------------- |
| `--title TEXT`   | 页面标题（单 URL 模式）                                          |
| `--desc TEXT`    | 页面描述（单 URL 模式）                                          |
| `--json JSON`    | 批量卡片 JSON 数组 `[{"url","title","description"}]`             |
| `--folder NAME`  | 目标收藏夹名称（如 `"父级/子级"`）                                |
| `--tag NAME,...` | 标签名称（逗号分隔，支持嵌套如 `"父级/子级"`）                     |

### `cubox-cli update`

更新卡片属性。

```bash
cubox-cli update --id CARD_ID [flags]
```


| 参数                     | 说明                                                      |
| ---------------------- | ------------------------------------------------------- |
| `--star` / `--unstar`  | 星标/取消星标                                               |
| `--read` / `--unread`  | 已读/未读                                                  |
| `--folder NAME`        | 按名称移动到收藏夹（如 `"父级/子级"`；`""` = 未分类）           |
| `--tag NAME,...`       | 按名称替换全部标签（逗号分隔，支持嵌套如 `"父级/子级"`）           |
| `--add-tag NAME,...`   | 在原有标签基础上新增标签                                     |
| `--remove-tag NAME,...`| 仅移除指定的标签                                             |
| `--title TEXT`         | 更新标题                                                   |
| `--description TEXT`   | 更新描述                                                   |


> 归档/取消归档为批量操作，已拆分为独立命令——见下方 [`cubox-cli archive`](#cubox-cli-archive)。

### `cubox-cli archive`

按 ID 批量归档收藏卡片。归档后的卡片不会出现在默认的 `card list` 中（使用 `card list --archived` 查看归档列表）。

```bash
cubox-cli archive --id CARD_ID[,ID2,...]
```

| 参数            | 说明                              |
| ------------- | ------------------------------- |
| `--id ID,...` | 要归档的卡片 ID（逗号分隔，必填）            |

### `cubox-cli unarchive`

通过将归档卡片移动到一个非归档收藏夹，批量恢复归档的收藏卡片。**目标收藏夹为必填**。

```bash
cubox-cli unarchive --id CARD_ID[,ID2,...] --folder NAME
```

| 参数              | 说明                                                                  |
| --------------- | ------------------------------------------------------------------- |
| `--id ID,...`   | 要恢复的卡片 ID（逗号分隔，必填）                                              |
| `--folder NAME` | 目标收藏夹名称（必填；`""` 表示未分类；嵌套写法 `"父级/子级"`）                       |

收藏夹名称会先通过 `folder list` 在客户端解析为 `folder_id`；如果名称无法匹配现有收藏夹，命令会以可读错误终止。

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

### `cubox-cli annotation list`

列出和搜索所有卡片的高亮标注。

```bash
cubox-cli annotation list [flags]
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
cubox-cli folder list -o text
cubox-cli card list --folder 7230156249357091393 --limit 10
```

### 阅读文章并查看 AI 洞察

```bash
cubox-cli card detail --id 7247925101516031380 -o pretty
```

### 保存链接并加星标

```bash
cubox-cli save https://example.com --title "示例网站"
cubox-cli update --id CARD_ID --star --read
```

### 归档与恢复卡片

```bash
# 批量归档
cubox-cli archive --id 7444025677600260245,7443973659296793971

# 浏览归档列表
cubox-cli card list --archived --limit 10

# 批量恢复（移动到非归档收藏夹）
cubox-cli unarchive --id 7444025677600260245,7443973659296793971 --folder "阅读清单"
```

### 整理标签（重命名 / 删除 / 合并）

```bash
# 重命名标签（仅修改叶子；子标签自动跟随）
cubox-cli tag update --id 7295070793040398540 --new-name 链接

# 批量删除标签（卡片本身保留）
cubox-cli tag delete --id 7444025677600260245,7443973659296793971

# 合并标签（源标签的卡片转移到目标标签，源标签随后被删除）
cubox-cli tag merge --source 7342187912403881105,7342187917722258501 --target 7247925099053977508
```

### 删除卡片（带 dry-run）

```bash
cubox-cli delete --id 7435692934957108160 --dry-run -o pretty
cubox-cli delete --id 7435692934957108160
```

### 导出所有标注

```bash
cubox-cli annotation list --all -o pretty
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
    folder.go             # folder list
    tag.go                # tag list/update/delete/merge
    card.go               # card list/detail/rag
    save.go               # save 保存网页
    update.go             # update 更新卡片
    archive.go            # 批量归档 / 恢复归档
    delete.go             # delete 删除卡片（支持 dry-run）
    annotation.go         # annotation 标注列表
    version.go            # version
  internal/
    client/               # HTTP 客户端 + API 类型
    config/               # 配置文件管理
    timefmt/              # 灵活时间解析
    update/               # npm 更新检查
  scripts/                # npm 分发包装
  skills/cubox/           # AI Agent Skill
  .github/workflows/      # CI/CD
```

## 许可证

[MIT](LICENSE)
