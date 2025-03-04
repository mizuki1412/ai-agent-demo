package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mizuki/project/ai-agent-demo/aikit"
)

func main() {
	r := &cobra.Command{
		Use: "main",
		Run: func(cmd *cobra.Command, args []string) {

			//ctx := context.Background()
			apiKey := configkit.GetString("ai.key")
			if apiKey == "" {
				log.Fatal("ai api key is not set")
			}

			client := aikit.NewChatModelClient(aikit.ChatModelConfig{
				APIKey:  apiKey,
				Model:   "deepseek-r1-distill-llama-70b",
				BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
			})
			client.Request([]schema.Message{
				{
					Role:    schema.User,
					Content: "说下你是谁？",
				},
			})

			//fmt.Println("\n=== Basic Chat ===")
			//basicChat(ctx, cm)

			//fmt.Println("\n=== Streaming Chat ===")
			//streamingChat(ctx, cm)
			//
			//fmt.Println("\n=== Prefix ===")
			//prefixChat(ctx, cm)
		},
	}
	r.Flags().String("ai.key", "", "")
	cli.RootCMD(r)
	cli.Execute()
}

func basicChat(ctx context.Context, cm model.ChatModel) {
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: "说下你是谁？",
		},
	}

	resp, err := cm.Generate(ctx, messages)
	if err != nil {
		log.Printf("Generate error: %v", err)
		return
	}

	reasoning, ok := deepseek.GetReasoningContent(resp)
	if !ok {
		fmt.Printf("Unexpected: non-reasoning")
	} else {
		fmt.Printf("Resoning Content: %s\n", reasoning)
	}
	fmt.Printf("Assistant: %s\n", resp.Content)
	if resp.ResponseMeta != nil && resp.ResponseMeta.Usage != nil {
		fmt.Printf("Tokens used: %d (prompt) + %d (completion) = %d (total)\n",
			resp.ResponseMeta.Usage.PromptTokens,
			resp.ResponseMeta.Usage.CompletionTokens,
			resp.ResponseMeta.Usage.TotalTokens)
	}
}

func streamingChat(ctx context.Context, cm model.ChatModel) {
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: "说下你是谁？",
		},
	}

	stream, err := cm.Stream(ctx, messages)
	if err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Stream receive error: %v", err)
			return
		}
		if reasoning, ok := deepseek.GetReasoningContent(resp); ok {
			fmt.Printf("Resoning Content: %s\n", reasoning)
		}
		if len(resp.Content) > 0 {
			fmt.Printf("Content: %s\n", resp.Content)
		}
		if resp.ResponseMeta != nil && resp.ResponseMeta.Usage != nil {
			fmt.Printf("Tokens used: %d (prompt) + %d (completion) = %d (total)\n",
				resp.ResponseMeta.Usage.PromptTokens,
				resp.ResponseMeta.Usage.CompletionTokens,
				resp.ResponseMeta.Usage.TotalTokens)
		}
	}
}

func prefixChat(ctx context.Context, cm model.ChatModel) {
	messages := []*schema.Message{
		schema.UserMessage("Please write quick sort code"),
		schema.AssistantMessage("```python\n", nil),
	}
	deepseek.SetPrefix(messages[1])

	result, err := cm.Generate(ctx, messages)
	if err != nil {
		log.Printf("Generate error: %v", err)
	}

	reasoningContent, ok := deepseek.GetReasoningContent(result)
	if !ok {
		fmt.Printf("No reasoning content")
	} else {
		fmt.Printf("Reasoning: %v\n", reasoningContent)
	}
	fmt.Printf("Content: %v\n", result)
}
