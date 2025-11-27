# telegram-ai

Backend service for Telegram web client that integrates AI capabilities using Firebase Genkit and OpenAI. Provides streaming responses via SSE.

## Requirements

- **Taskfile**: [Install](https://taskfile.dev) or `brew install go-task` (macOS)
- **golangci-lint**: [Install](https://golangci-lint.run) or `brew install golangci-lint` (macOS)
- **hadolint**: [Install](https://hadolint.com) or `brew install hadolint` (macOS)

## Quick Start

1. Create `.env` file from template:
```bash
cp .env.example .env
```

2. Run the application:

**For local development:**
```bash
task http:dev
```

**For production build:**
```bash
task http:build
task http:run
```

**With Docker:**
```bash
task docker:http:build
task docker:http:run
```

3. Verify the server is running on `http://localhost`:
```bash
curl http://localhost/_hc
```
