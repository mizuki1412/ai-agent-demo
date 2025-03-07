package aikit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"log"
	"mizuki/project/ai-agent-demo/aikit/schema"
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
	decoder := newApiResDecoder()
	req := schema.RequestBody{
		Model:         client.Config.Model,
		Messages:      messages,
		Stream:        true,
		StreamOptions: schema.StreamOption{IncludeUsage: true},
	}
	overChan := make(chan bool)
	httpkit.Request(httpkit.Req{
		Url:         client.Config.BaseURL,
		Method:      "post",
		ContentType: "application/json",
		Header: map[string]string{
			"Authorization": "Bearer " + client.Config.APIKey,
		},
		JsonData: req,
		Stream:   true,
		StreamHandler: func(data []byte) {
			decoder.Put(data)
		},
	})
	decoder.Recv(func(bytes []byte) {
		println(string(bytes))
		res := &schema.ResponseBody{}
		jsonkit.ParseObj(string(bytes), res)
		// todo

		overChan <- true
	})
	<-overChan
	log.Println("finish")
}

func newApiResDecoder() *framekit.Decoder {
	return framekit.NewDecoder(1024, func(bytes []byte) ([]byte, []byte, bool) {
		// 百炼的格式： data: {} ; data: [DONE]
		i := 0
		// 找到json字符串起点
		beginFlag := 0
		// 存放json首尾标记符
		jsonFlags := make([]byte, 0, 10)
		for {
			if beginFlag == 0 {
				// 寻找data:
				if len(bytes) <= i+6+6 {
					break
				}
				// 结束
				if string(bytes[i:i+6+6]) == "data: [DONE]" {
					return bytes, nil, true
				}
				if string(bytes[i:i+6]) == "data: " {
					i += 6
					beginFlag = i
					continue
				}
			} else {
				if i >= len(bytes) {
					break
				}
				switch bytes[i] {
				case '[', '{':
					// 排除作为内容含义的[{
					if len(jsonFlags) > 0 && jsonFlags[len(jsonFlags)-1] != '"' {
						jsonFlags = append(jsonFlags, bytes[i])
					}
				case '"':
					if len(jsonFlags) > 0 && jsonFlags[len(jsonFlags)-1] != '"' {
						jsonFlags = append(jsonFlags, bytes[i])
					} else if len(jsonFlags) > 0 && jsonFlags[len(jsonFlags)-1] == '"' {
						jsonFlags = jsonFlags[:len(jsonFlags)-1]
					}
				case ']':
					if len(jsonFlags) > 0 && jsonFlags[len(jsonFlags)-1] == '[' {
						jsonFlags = jsonFlags[:len(jsonFlags)-1]
					}
				case '}':
					if len(jsonFlags) > 0 && jsonFlags[len(jsonFlags)-1] == '{' {
						jsonFlags = jsonFlags[:len(jsonFlags)-1]
					}
					// 识别是否结束
					if len(jsonFlags) == 0 {
						if i+1 == len(bytes) {
							return []byte{}, bytes[beginFlag : i+1], false
						}
						return bytes[i:], bytes[beginFlag : i+1], false
					}
				}
			}
			i++
		}
		return bytes, nil, false
	})
}
