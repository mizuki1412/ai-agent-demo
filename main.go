package main

import (
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/aikit"
	"github.com/mizuki1412/go-core-kit/v2/service/aikit/schema"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/spf13/cobra"
	"log"
	"time"
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
			start := time.Now()
			res, usage := client.Request([]schema.Message{
				{
					Role:    schema.User,
					Content: "说下你是谁？",
				},
			})
			log.Println("推理过程:  " + res.ReasoningContent)
			log.Println("结果:  " + res.Content)
			if usage != nil {
				log.Println("消耗:  " + jsonkit.ToString(usage))
			}
			log.Println("耗时:  " + time.Since(start).String())
		},
	}
	r.Flags().String("ai.key", "", "")
	cli.RootCMD(r)
	cli.Execute()
}
