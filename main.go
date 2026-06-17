package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	backend := flag.String("backend", "local", `AI backend: "local" (Docker Model Runner) or "bedrock" (Amazon Bedrock)`)
	file := flag.String("file", "sample-compose.yml", "Path to the file to summarize")
	model := flag.String("model", "", "Override the default model (e.g. ai/phi4, anthropic.claude-3-5-haiku-20241022-v1:0)")
	stream := flag.Bool("stream", true, "Stream tokens as they arrive (set false to wait for full response)")
	flag.Parse()

	content, err := os.ReadFile(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %q: %v\n", *file, err)
		os.Exit(1)
	}

	ctx := context.Background()

	var s Summarizer
	switch *backend {
	case "local":
		s = NewLocalSummarizer(*model, *stream)
	case "bedrock":
		s, err = NewBedrockSummarizer(ctx, *model, *stream)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bedrock init: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown backend %q — use \"local\" or \"bedrock\"\n", *backend)
		os.Exit(1)
	}

	streamLabel := "streaming"
	if !*stream {
		streamLabel = "batch"
	}

	fmt.Printf("\n🤖  Backend : %s\n", s.Name())
	fmt.Printf("🧠  Model   : %s\n", s.Model())
	fmt.Printf("📄  File    : %s\n", *file)
	fmt.Printf("⚡  Mode    : %s\n", streamLabel)
	fmt.Println(divider())

	start := time.Now()

	if err := s.Summarize(ctx, string(content)); err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(divider())
	fmt.Printf("⏱   Done in %.2fs\n\n", time.Since(start).Seconds())
}

func divider() string {
	return "─────────────────────────────────────────"
}
