package openai

import (
	"bilibili_subtitle/internal/config"
	"context"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
	"strings"
	"time"
)

type OpenaiClient struct {
	Config       *config.OpenaiModelConfig
	openaiClient *openai.Client
	httpClient   *http.Client
}

func NewOpenAIClient(cfg *config.Config) (*OpenaiClient, error) {
	config := openai.DefaultConfig(cfg.OpenaiAPIKey)
	config.BaseURL = cfg.OpenaiModelConfig.Endpoint
	client := openai.NewClientWithConfig(config)

	return &OpenaiClient{
		Config:       &cfg.OpenaiModelConfig,
		openaiClient: client,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.OpenaiModelConfig.Timeout) * time.Second,
		},
	}, nil
}

func (c *OpenaiClient) AnalyzeSubtitles(ctx context.Context, prompt, text string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second) // Ensure the overall operation respects this timeout
	defer cancel()

	parts, err := splitTextIntoParts(text, prompt, c.Config)
	if err != nil {
		fmt.Printf("Error splitting text: %v\n", err)
		return "", err
	}

	var fullResponse string

	// Start the conversation with the initial prompt
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: prompt},
	}

	for _, part := range parts {
		// Append user message for each part
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: part,
		})

		// Create a chat completion for the current set of messages
		resp, err := c.openaiClient.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       c.Config.ModelName,
				MaxTokens:   c.Config.MaxTokens,
				TopP:        c.Config.TopP,
				Temperature: c.Config.Temperature,
				Messages:    messages,
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return "", err
		}

		// Collect the response and prepare for the next iteration
		lastMessage := resp.Choices[0].Message.Content
		fullResponse += lastMessage + " " // Append the response and a space to separate responses

		// Reset messages to just contain the last part of the dialogue to maintain context without redundancy
		messages = []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: lastMessage},
		}
	}

	return strings.TrimSpace(fullResponse), nil
}

// 辅助函数：根据最大tokens拆分长文本
func splitTextIntoParts(text, prompt string, Config *config.OpenaiModelConfig) ([]string, error) {
	tkm, err := tiktoken.EncodingForModel(Config.ModelName)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return nil, err
	}

	// Encode the prompt to understand its token count
	promptTokens := tkm.Encode(prompt, nil, nil)
	promptTokenCount := len(promptTokens)

	words := strings.Fields(text)
	var parts []string
	var currentPart []string

	for _, word := range words {
		tempPart := append(currentPart, word)
		tempEncoded := tkm.Encode(strings.Join(tempPart, " "), nil, nil)
		totalTokens := len(tempEncoded) + promptTokenCount
		if totalTokens > Config.MaxTokens {
			if len(currentPart) == 0 {
				// Handle edge case where a single word is too large after a prompt
				currentPart = append(currentPart, word) // Force add the word to ensure progress
				parts = append(parts, strings.Join(currentPart, " "))
				currentPart = []string{}
			} else {
				parts = append(parts, strings.Join(currentPart, " "))
				currentPart = []string{word} // Start new part with the current word
			}
		} else {
			currentPart = tempPart
		}
	}
	if len(currentPart) > 0 {
		parts = append(parts, strings.Join(currentPart, " "))
	}

	return parts, nil
}
