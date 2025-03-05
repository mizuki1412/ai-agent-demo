package aikit

import (
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"mizuki/project/ai-agent-demo/framekit"
	"mizuki/project/ai-agent-demo/httpkit"
)

// ChatModelConfig api连接配置信息
type ChatModelConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	// 模型名称
	Model string `json:"model"`

	MaxTokens int `json:"max_tokens,omitempty"`
	Timeout   int `json:"timeout"` // seconds
}

type ChatModelClient struct {
	Config ChatModelConfig
}

func NewChatModelClient(config ChatModelConfig) *ChatModelClient {
	if config.APIKey == "" {
		panic(exception.New("api key is nil"))
	}
	if config.BaseURL == "" {
		panic(exception.New("api baseUrl is nil"))
	}
	if config.Model == "" {
		panic(exception.New("model is nil"))
	}
	return &ChatModelClient{
		Config: config,
	}
}

func (client *ChatModelClient) Request(messages []schema.Message) {
	var all []byte
	httpkit.Request(httpkit.Req{
		Url:         client.Config.BaseURL,
		Method:      "post",
		ContentType: "application/json",
		Header: map[string]string{
			"Authorization": "Bearer " + client.Config.APIKey,
		},
		JsonData: map[string]any{
			"model":    client.Config.Model,
			"messages": messages,
			"stream":   true,
			"stream_options": map[string]any{
				"include_usage": true,
			},
		},
		Stream: true,
		StreamHandler: func(data []byte) {
			all = append(all, data...)
		},
	})
	fmt.Println(string(all))
}

func newApiResDecoder() *framekit.Decoder {
	return framekit.NewDecoder(1024, func(bytes []byte) []byte {
		// 百炼的格式： data: {} ; data: [DONE]
	})
}
