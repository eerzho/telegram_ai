# telegram_ai

Backend service for Telegram web client that integrates AI capabilities using Firebase Genkit and OpenAI. Provides streaming responses via SSE.

## Requirements

- **taskfile**: [Install](https://taskfile.dev) or `brew install go-task` (macOS)
- **golangci-lint**: [Install](https://golangci-lint.run) or `brew install golangci-lint` (macOS)

## Quick Start

Create `.env` file from template:
```bash
cp .env.example .env
```

Run containers:
```bash
docker compose up -d
```

Run the dev server:
```bash
task http:dev
```

To see all available commands:
```bash
task
```
