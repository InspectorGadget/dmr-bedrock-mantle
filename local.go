package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	localBaseURL      = "http://localhost:12434/engines/v1"
	localModelDefault = "ai/gemma3"
)

type LocalSummarizer struct {
	client openai.Client
	model  string
	stream bool
}

func NewLocalSummarizer(model string, stream bool) *LocalSummarizer {
	if model == "" {
		model = localModelDefault
	}
	client := openai.NewClient(
		option.WithBaseURL(localBaseURL),
		option.WithAPIKey("not-required"),
	)
	return &LocalSummarizer{client: client, model: model, stream: stream}
}

func (s *LocalSummarizer) Name() string  { return "Docker Model Runner" }
func (s *LocalSummarizer) Model() string { return s.model }

func (s *LocalSummarizer) Summarize(ctx context.Context, content string) error {
	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModel(s.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(content),
		},
	}

	if s.stream {
		return s.summarizeStreaming(ctx, params)
	}
	return s.summarizeBatch(ctx, params)
}

func (s *LocalSummarizer) summarizeStreaming(ctx context.Context, params openai.ChatCompletionNewParams) error {
	stream := s.client.Chat.Completions.NewStreaming(ctx, params)
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
	return stream.Err()
}

func (s *LocalSummarizer) summarizeBatch(ctx context.Context, params openai.ChatCompletionNewParams) error {
	resp, err := s.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}
	if len(resp.Choices) > 0 {
		fmt.Println(resp.Choices[0].Message.Content)
	}
	return nil
}
