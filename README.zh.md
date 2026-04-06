# cubox-cli

[License: MIT](https://opensource.org/licenses/MIT)
[Go Version](https://go.dev/)
[npm version](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

[Cubox](https://cubox.pro) 官方命令行工具 — 为人类用户和 AI Agent 而设计。在终端中管理书签、浏览收藏夹、阅读保存的内容。

[安装](#安装) · [认证](#认证) · [命令](#命令) · [AI Agent](#ai-agent-快速开始) · [示例](#示例) · [开发](#开发)

## 功能


| 类别  | 能力                             |
| --- | ------------------------------ |
| 收藏夹 | 列出和浏览书签文件夹                     |
| 标签  | 列出和浏览标签层级                      |
| 卡片  | 按收藏夹、类型、标签、星标/已读/标注状态过滤卡片；游标分页 |
| 内容  | 以 Markdown 格式阅读文章全文            |


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
cubox-cli card content --id CARD_ID
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

过滤和列出收藏卡片。

```bash
cubox-cli card list [flags]
```


| 参数                 | 说明                                                                    |
| ------------------ | --------------------------------------------------------------------- |
| `--group ID,...`   | 按收藏夹 ID 过滤                                                            |
| `--type TYPE,...`  | 按类型过滤：`Article`, `Snippet`, `Memo`, `Image`, `Audio`, `Video`, `File` |
| `--tag ID,...`     | 按标签 ID 过滤                                                             |
| `--starred`        | 仅星标卡片                                                                 |
| `--read`           | 仅已读卡片                                                                 |
| `--unread`         | 仅未读卡片                                                                 |
| `--annotated`      | 仅有标注的卡片                                                               |
| `--limit N`        | 每页数量（默认 50）                                                           |
| `--cursor ID,TIME` | 分页游标，从上一张卡片继续                                                         |
| `--all`            | 自动翻页获取全部结果                                                            |


**JSON 输出字段：** `id`, `title`, `description`, `article_title`, `domain`, `type`, `tags`, `url`, `cubox_url`, `words_count`, `create_time`, `update_time`, `highlights`

### `cubox-cli card content --id ID`

获取文章全文（Markdown 格式）。

```bash
cubox-cli card content --id 7247925101516031380
```

默认输出原始 Markdown。使用 `-o pretty` 获取 JSON 包装格式。

## 示例

### 列出所有星标文章

```bash
cubox-cli card list --starred --type Article -o pretty
```

### 浏览特定收藏夹的卡片

```bash
# 查找收藏夹 ID
cubox-cli group list -o text

# 列出该收藏夹的卡片
cubox-cli card list --group 7230156249357091393 --limit 10
```

### 阅读保存的文章

```bash
cubox-cli card content --id 7247925101516031380
```

### 获取所有带标注的卡片

```bash
cubox-cli card list --annotated --all -o pretty
```

### 分页

```bash
# 第一页
cubox-cli card list --limit 5

# 使用最后一张卡片的 ID 和 update_time 获取下一页
cubox-cli card list --limit 5 --cursor "7247925102807877551,2024-12-04T16:23:01:347+08:00"
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
    card.go               # card list, card content
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