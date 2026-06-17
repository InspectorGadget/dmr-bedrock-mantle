# dmr-bedrock-mantle

> Demo code for the talk **"Docker Model Runner Meets Amazon Bedrock: Local-to-Cloud AI Workflows"**
> Presented at the Docker KL × AWS User Group Malaysia joint meetup.

Go CLI that summarizes config files — Docker Compose, CloudFormation, Kubernetes manifests, whatever — using AI. One flag switches between a **local model via Docker Model Runner** and a **cloud model via Amazon Bedrock**. Same code. Same prompt. One flag.

---

## How it works

```
┌─────────────────────────────────────────────────────────┐
│                     summarizer CLI                       │
│                                                          │
│   -backend local          │   -backend bedrock           │
│   -model  ai/gemma3       │   -model  claude-3-haiku     │
│   -stream true/false      │   -stream true/false         │
└───────────┬───────────────┴──────────────┬───────────────┘
            │                              │
            ▼                              ▼
┌───────────────────────┐    ┌─────────────────────────────┐
│   Docker Model Runner  │    │       Amazon Bedrock         │
│  (Docker Desktop)      │    │   Converse API (streaming)   │
│                        │    │                              │
│  OpenAI-compatible API │    │  anthropic.claude-3-haiku    │
│  localhost:12434       │    │  (or any enabled model)      │
│                        │    │                              │
│  ai/gemma3 (default)   │    │  ap-southeast-1              │
└───────────────────────┘    └─────────────────────────────┘
```

Both backends implement `Summarizer` from `backend.go` — that's the whole abstraction.

---

## Prerequisites

| Requirement          | Notes                                |
| -------------------- | ------------------------------------ |
| Go 1.23+             | `go version` to check                |
| Docker Desktop 4.40+ | With **Docker Model Runner** enabled |
| AWS account          | With Bedrock model access granted    |
| AWS credentials      | `aws configure` or env vars          |

### Enable Docker Model Runner

Docker Desktop → Settings → Features in development → Enable Docker Model Runner

### Enable Bedrock model access

AWS Console → Amazon Bedrock → Model access → Enable **Anthropic Claude 3 Haiku**

---

## Setup

```bash
# 1. Clone the repo
git clone https://github.com/InspectorGadget/dmr-bedrock-mantle
cd dmr-bedrock-mantle

# 2. Pull Go dependencies
make tidy

# 3. Pull the default local model (2.5 GB, one-time)
make pull-model
```

---

## Running the demo

### Core demo — same file, two backends

```bash
# Local: Gemma 3 via Docker Model Runner (no cloud, no cost)
make run-local

# Cloud: Claude Haiku via Amazon Bedrock
make run-bedrock
```

### CloudFormation template

```bash
# Show it works on any config format, not just Compose
make demo-cfn-local
make demo-cfn-bedrock
```

### Streaming vs batch mode

```bash
# Default: tokens stream in as they're generated
make run-local

# Batch: wait for the full response, then print
make demo-no-stream-local
make demo-no-stream-bedrock
```

### Live model swap

```bash
# Pull additional models first
docker model pull ai/phi4
docker model pull ai/llama3.2

# Then swap with a flag — no code change
make demo-phi4
make demo-llama

# On Bedrock — pass any model ID directly
go run . -backend bedrock -model anthropic.claude-3-5-sonnet-20241022-v2:0 -file sample-compose.yml
```

### Custom file

```bash
# Summarize any file you have
make run-custom FILE=your-k8s-deployment.yaml
make run-custom-bedrock FILE=your-terraform.tf
```

### Direct CLI usage

```bash
go run . \
  -backend bedrock \
  -file    sample-cfn.yaml \
  -model   anthropic.claude-3-5-sonnet-20241022-v2:0 \
  -stream  true
```

---

## Project structure

```
.
├── main.go            # CLI entry point — flags, banner, timing
├── backend.go         # Summarizer interface + shared system prompt
├── local.go           # Docker Model Runner backend (OpenAI-compatible API)
├── bedrock.go         # Amazon Bedrock backend (Converse API)
├── sample-compose.yml # Demo input: multi-service Docker Compose stack
├── sample-cfn.yaml    # Demo input: ECS + Aurora + ElastiCache CloudFormation
├── Makefile           # All demo targets
└── go.mod
```

### Key design choice

`local.go` and `bedrock.go` both implement this interface from `backend.go`:

```go
type Summarizer interface {
    Summarize(ctx context.Context, content string) error
    Name() string
    Model() string
}
```

Switching backends is a flag, not a code change. This is the point of the talk.

---

## Flags

| Flag       | Default              | Description                                                                                               |
| ---------- | -------------------- | --------------------------------------------------------------------------------------------------------- |
| `-backend` | `local`              | `local` (Docker Model Runner) or `bedrock` (Amazon Bedrock)                                               |
| `-file`    | `sample-compose.yml` | Path to the file to summarize                                                                             |
| `-model`   | _(backend default)_  | Override the model. Local default: `ai/gemma3`. Bedrock default: `anthropic.claude-3-haiku-20240307-v1:0` |
| `-stream`  | `true`               | `true` streams tokens as generated. `false` waits for full response.                                      |

---

## AWS region note

This demo was built and tested in `ap-southeast-1` (Singapore). Amazon Nova models require cross-region inference profiles and may not be available in all regions — Claude 3 Haiku is used as the Bedrock default because it's natively available in `ap-southeast-1` without additional setup.

To use Nova models from a US region:

```bash
go run . -backend bedrock -model us.amazon.nova-lite-v1:0
```

---

## License

MIT
