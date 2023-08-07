package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

const token = "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"

func main() {
	// 1. 设置配置 & 代理
	config := openai.DefaultConfig(token)
	proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	// 2.创建客户端
	client := openai.NewClientWithConfig(config)

	// 3.语音转文字
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "audio.wav",
	}
	resp1, err := client.CreateTranscription(context.Background(), req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp1.Text)

	// 4.聊天
	resp2, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "将下面这句话翻译成英文: " + resp1.Text,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp2.Choices[0].Message.Content)
}
