# aigit

[English](./README.md)

AI 驱动的 Git 提交信息生成器和代码审查工具。

## 功能特性

- **AI 生成提交信息**：根据暂存的代码变更自动生成有意义的 commit message
- **代码审查**：分析代码变更，识别潜在的 bug、安全问题和代码质量问题
- **多 AI 服务商支持**：支持 OpenAI、Claude、Google Gemini 和 OpenRouter
- **中英双语**：支持中文和英文输出
- **Git Hooks**：提交前自动审查代码

## 安装

```bash
go install github.com/go-goll/aigit@latest
```

## 快速开始

### 1. 配置 AI 服务商

```bash
aigit config
```

按提示操作：
- 选择 AI 服务商（OpenAI/Claude/Google/OpenRouter）
- 输入 API Key
- 选择模型（可选，使用默认值）
- 选择语言（英文/中文）

**OpenRouter 快速配置：**
```bash
aigit config provider openrouter
aigit config api_key sk-or-xxx
aigit config model anthropic/claude-sonnet-4-20250514
```

### 2. 生成提交信息

```bash
# 先暂存变更
git add .

# 生成提交信息
aigit commit

# 或者暂存所有变更并提交
aigit commit -a

# 自动提交，无需确认
aigit commit -y
```

### 3. 代码审查

```bash
# 审查所有变更（已暂存 + 未暂存）
aigit review

# 仅审查已暂存的变更
aigit review -s
```

### 4. 安装 Git Hooks（可选）

```bash
# 安装 pre-commit hook，提交前自动审查
aigit hooks install

# 卸载
aigit hooks uninstall
```

## 命令列表

| 命令 | 说明 |
|------|------|
| `aigit config` | 配置 AI 服务商和设置 |
| `aigit commit` | 为暂存的变更生成提交信息 |
| `aigit review` | 审查代码变更，查找潜在问题 |
| `aigit hooks install` | 安装 pre-commit hook |
| `aigit hooks uninstall` | 卸载 pre-commit hook |

### Commit 参数

| 参数 | 说明 |
|------|------|
| `-a, --all` | 提交前暂存所有变更 |
| `-y, --yes` | 自动提交，无需确认 |

### Review 参数

| 参数 | 说明 |
|------|------|
| `-s, --staged` | 仅审查已暂存的变更 |

## 配置

配置文件存储在 `~/.aigit/config.json`：

```json
{
  "provider": "openrouter",
  "api_key": "your-api-key",
  "model": "anthropic/claude-sonnet-4-20250514",
  "language": "zh",
  "base_url": ""
}
```

### 支持的服务商

| 服务商 | 默认模型 | API Key 获取 |
|--------|----------|--------------|
| OpenAI | gpt-4o | [OpenAI Platform](https://platform.openai.com/api-keys) |
| Claude | claude-sonnet-4-20250514 | [Anthropic Console](https://console.anthropic.com/) |
| Google | gemini-1.5-pro | [Google AI Studio](https://aistudio.google.com/app/apikey) |
| OpenRouter | anthropic/claude-sonnet-4-20250514 | [OpenRouter](https://openrouter.ai/keys) |

### 自定义 Base URL

如需使用代理或自定义端点：

```json
{
  "provider": "openai",
  "api_key": "your-api-key",
  "base_url": "https://your-proxy.com/v1"
}
```

## 使用示例

### 生成提交信息

```
$ aigit commit
Staged files:
  • internal/auth/handler.go
  • internal/auth/middleware.go

Generating commit message...

--- Generated Commit Message ---
feat(auth): 添加 JWT 认证中间件

- 实现 JWT token 验证
- 添加用户上下文提取
- 支持 token 刷新机制
--------------------------------

Commit with this message? [Y/n/e(dit)]: y
✓ Committed successfully!
```

### 代码审查

```
$ aigit review -s
Reviewing staged changes...

=== Code Review Results ===
**高**：`internal/db/query.go:45` 存在 SQL 注入风险
- 用户输入直接拼接到 SQL 查询中
- 修复：使用参数化查询

**中**：`internal/api/handler.go:78` 缺少错误处理
- HTTP 响应未检查错误
- 修复：添加 response.Body.Close() 的错误处理

**低**：`internal/utils/helper.go:23` 存在未使用的变量
- 变量 `temp` 已声明但未使用
===========================
```

## 许可证

MIT License
