package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lengzhao/conf"
	"github.com/lengzhao/gptui/history"
	"github.com/sashabaranov/go-openai"
)

type Chat struct {
	client       *openai.Client
	systemPrompt string
	history      history.History
	req          openai.ChatCompletionRequest
}

func newGPTClient() *openai.Client {
	var c openai.ClientConfig
	if os.Getenv("gptType") == "azure" {
		if os.Getenv("azureApiKey") == "" || os.Getenv("azureEndpoint") == "" {
			log.Println("lost config:azureApiKey or azureEndpoint")
			return nil
		}
		c = openai.DefaultAzureConfig(os.Getenv("azureApiKey"), os.Getenv("azureEndpoint"))
		// fmt.Println("use azure gpt")
	} else {
		if os.Getenv("gptType") != "openai" {
			log.Println("unknow config:gptType")
			return nil
		}
		if os.Getenv("OPENAI_API_KEY") == "" {
			log.Println("lost config:OPENAI_API_KEY")
			return nil
		}
		c = openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	}
	c.HTTPClient = &http.Client{
		Timeout: time.Duration(conf.GetInt("HTTP_TIMEOUT", 20)) * time.Second,
	}
	return openai.NewClientWithConfig(c)
}

func New() *Chat {
	client := newGPTClient()
	if client == nil {
		return nil
	}
	return NewWithClient(client)
}

func NewWithClient(client *openai.Client) *Chat {
	var out Chat
	out.client = client
	out.systemPrompt = os.Getenv("prompt")
	out.req = openai.ChatCompletionRequest{
		Model:       conf.Get("model", openai.GPT3Dot5Turbo),
		MaxTokens:   conf.GetInt("MaxTokens", 2000),
		Temperature: float32(conf.GetFloat("Temperature", 1)),
		TopP:        float32(conf.GetFloat("TopP", 0.9)),
	}
	return &out
}

func (c *Chat) Send(input string) (string, error) {
	c.req.Messages = make([]openai.ChatCompletionMessage, 0)
	if c.systemPrompt != "" {
		c.req.Messages = append(c.req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: c.systemPrompt,
		})
	}
	hItem := c.history.Get(conf.GetInt("HistoryLimit", 2) * 2)
	for _, it := range hItem {
		c.req.Messages = append(c.req.Messages, openai.ChatCompletionMessage{
			Role:    it.Role,
			Content: it.Text,
		})
	}
	c.req.Messages = append(c.req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		c.req,
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}
	if len(resp.Choices) < 1 {
		return "", fmt.Errorf("not response, choices.length = 0")
	}
	out := resp.Choices[0].Message.Content
	c.history.Add("user", input)
	c.history.Add(openai.ChatMessageRoleAssistant, out)

	return out, nil
}

func (c *Chat) Reset() {
	c.history = history.History{}
}

func (c *Chat) SetSystemPrompt(prompt string) {
	fmt.Println("set system prompt:", prompt)
	c.systemPrompt = prompt
}
