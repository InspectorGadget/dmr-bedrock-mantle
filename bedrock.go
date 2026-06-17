package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// Anthropic Claude 3 Haiku on Bedrock — fast, available natively in ap-southeast-1.
const bedrockModelDefault = "anthropic.claude-3-haiku-20240307-v1:0"

// Auth uses the standard AWS credential chain: env vars, ~/.aws/credentials, IAM role.
type BedrockSummarizer struct {
	client *bedrockruntime.Client
	model  string
	stream bool
}

func NewBedrockSummarizer(ctx context.Context, model string, stream bool) (*BedrockSummarizer, error) {
	if model == "" {
		model = bedrockModelDefault
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &BedrockSummarizer{
		client: bedrockruntime.NewFromConfig(cfg),
		model:  model,
		stream: stream,
	}, nil
}

func (s *BedrockSummarizer) Name() string  { return "Amazon Bedrock" }
func (s *BedrockSummarizer) Model() string { return s.model }

func (s *BedrockSummarizer) Summarize(ctx context.Context, content string) error {
	if s.stream {
		return s.summarizeStreaming(ctx, content)
	}
	return s.summarizeBatch(ctx, content)
}

func (s *BedrockSummarizer) converseInput(content string) ([]types.SystemContentBlock, []types.Message) {
	system := []types.SystemContentBlock{
		&types.SystemContentBlockMemberText{Value: systemPrompt},
	}
	messages := []types.Message{
		{
			Role: types.ConversationRoleUser,
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{Value: content},
			},
		},
	}
	return system, messages
}

func (s *BedrockSummarizer) summarizeStreaming(ctx context.Context, content string) error {
	system, messages := s.converseInput(content)
	output, err := s.client.ConverseStream(ctx, &bedrockruntime.ConverseStreamInput{
		ModelId:  aws.String(s.model),
		System:   system,
		Messages: messages,
	})
	if err != nil {
		return fmt.Errorf("bedrock converse stream: %w", err)
	}

	stream := output.GetStream()
	defer stream.Close()

	for event := range stream.Events() {
		switch v := event.(type) {
		case *types.ConverseStreamOutputMemberContentBlockDelta:
			if delta, ok := v.Value.Delta.(*types.ContentBlockDeltaMemberText); ok {
				fmt.Print(delta.Value)
			}
		}
	}
	fmt.Println()
	return stream.Err()
}

func (s *BedrockSummarizer) summarizeBatch(ctx context.Context, content string) error {
	system, messages := s.converseInput(content)
	output, err := s.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId:  aws.String(s.model),
		System:   system,
		Messages: messages,
	})
	if err != nil {
		return fmt.Errorf("bedrock converse: %w", err)
	}

	if msg, ok := output.Output.(*types.ConverseOutputMemberMessage); ok {
		for _, block := range msg.Value.Content {
			if text, ok := block.(*types.ContentBlockMemberText); ok {
				fmt.Println(text.Value)
			}
		}
	}
	return nil
}
