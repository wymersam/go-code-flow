package api

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

var FuncSummaries = make(map[string]string)

func GetFunctionSummary(code string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant. Summarise the following Go function in one or two sentences without repeating the function name.",
				},
				{
					Role:    "user",
					Content: code,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
