package main

import (
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/spf13/cobra"
	"log"
	"mizuki/project/ai-agent-demo/aikit"
	"mizuki/project/ai-agent-demo/aikit/schema"
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
			log.Println("over")
		},
	}
	r.Flags().String("ai.key", "", "")
	cli.RootCMD(r)
	cli.Execute()
}
