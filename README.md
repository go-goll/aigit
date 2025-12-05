# aigit

AI-powered git commit message generator and code reviewer.

## Features

- **AI Commit Messages**: Automatically generate meaningful commit messages based on your staged changes
- **Code Review**: Analyze code changes for potential bugs, security issues, and code quality problems
- **Multi-Provider Support**: Works with OpenAI (GPT-4), Claude, and Google Gemini
- **Bilingual**: Supports both English and Chinese output

## Installation

```bash
go install github.com/go-goll/aigit@latest
```

## Quick Start

### 1. Configure your AI provider

```bash
aigit config
```

Follow the prompts to:
- Select your AI provider (OpenAI/Claude/Google)
- Enter your API key
- Choose model (optional, uses default)
- Select language (English/Chinese)

### 2. Generate commit message

```bash
# Stage your changes first
git add .

# Generate commit message
aigit commit

# Or stage all and commit in one command
aigit commit -a

# Auto-commit without confirmation
aigit commit -y
```

### 3. Review code for bugs

```bash
# Review all changes (staged + unstaged)
aigit review

# Review only staged changes
aigit review -s
```

## Commands

| Command | Description |
|---------|-------------|
| `aigit config` | Configure AI provider and settings |
| `aigit commit` | Generate commit message for staged changes |
| `aigit review` | Review code changes for potential bugs |

### Commit Flags

| Flag | Description |
|------|-------------|
| `-a, --all` | Stage all changes before commit |
| `-y, --yes` | Auto-commit without confirmation |

### Review Flags

| Flag | Description |
|------|-------------|
| `-s, --staged` | Review only staged changes |

## Configuration

Configuration is stored in `~/.aigit/config.json`:

```json
{
  "provider": "openai",
  "api_key": "your-api-key",
  "model": "gpt-4o",
  "language": "en",
  "base_url": ""
}
```

### Supported Providers

| Provider | Default Model | API Key Source |
|----------|---------------|----------------|
| OpenAI | gpt-4o | [OpenAI Platform](https://platform.openai.com/api-keys) |
| Claude | claude-sonnet-4-20250514 | [Anthropic Console](https://console.anthropic.com/) |
| Google | gemini-1.5-pro | [Google AI Studio](https://aistudio.google.com/app/apikey) |

### Custom Base URL

For users who need to use a proxy or custom endpoint:

```json
{
  "provider": "openai",
  "api_key": "your-api-key",
  "base_url": "https://your-proxy.com/v1"
}
```

## Examples

### Generate commit message

```
$ aigit commit
Staged files:
  • internal/auth/handler.go
  • internal/auth/middleware.go

Generating commit message...

--- Generated Commit Message ---
feat(auth): add JWT authentication middleware

- Implement JWT token validation
- Add user context extraction
- Support token refresh mechanism
--------------------------------

Commit with this message? [Y/n/e(dit)]: y
✓ Committed successfully!
```

### Review code

```
$ aigit review -s
Reviewing staged changes...

=== Code Review Results ===
**HIGH**: Potential SQL injection in `internal/db/query.go:45`
- User input is directly concatenated into SQL query
- Fix: Use parameterized queries

**MEDIUM**: Missing error handling in `internal/api/handler.go:78`
- HTTP response not checked for errors
- Fix: Add error handling for response.Body.Close()

**LOW**: Unused variable in `internal/utils/helper.go:23`
- Variable `temp` is declared but never used
===========================
```

## License

MIT License
