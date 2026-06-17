package main

import "context"

// Summarizer is the single interface both backends implement.
type Summarizer interface {
	Summarize(ctx context.Context, content string) error
	Name() string
	Model() string
}

// systemPrompt is shared by both backends — same prompt, same output shape.
const systemPrompt = `You are a DevOps documentation assistant.
Summarize the provided configuration file in plain English.
Structure your response as:

**What this does** — 2-3 sentences a non-technical stakeholder can understand.

**Key components** — bullet list: component name and its role.

**Things to note** — any notable choices, gotchas, or security considerations.

Be concise. Avoid jargon where possible.`
