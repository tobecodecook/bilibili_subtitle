package gemini

import (
	"bilibili_subtitle/internal/config" // Make sure to import the correct path
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
	"strings"
	"time"
)

type GeminiClient struct {
	Config   *config.GeminiModelConfig
	aiClient *genai.Client
	sem      *semaphore.Weighted // 用于并发控制
}

func NewGeminiClient(cfg *config.Config) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini AI client: %w", err)
	}
	sem := semaphore.NewWeighted(cfg.GeminiModelConfig.MaxConcurrentRequests) // 设置最大并发请求数
	return &GeminiClient{
		Config:   &cfg.GeminiModelConfig,
		aiClient: client,
		sem:      sem,
	}, nil
}

// Helper function to split text into manageable parts
func splitTextToFitModel(prompt, text string, maxTokens int32) ([]string, error) {
	// Tokenize the text to ensure splitting respects word boundaries
	words := strings.Fields(text)
	var parts []string
	currentPart := prompt
	currentLen := len(strings.Fields(prompt)) // Start with the length of the prompt in tokens

	for _, word := range words {
		wordTokens := len(strings.Fields(word)) // Calculate once
		if int32(currentLen+wordTokens+1) > maxTokens {
			parts = append(parts, currentPart)
			currentPart = prompt
			currentLen = len(strings.Fields(prompt))
		}
		currentPart += " " + word
		currentLen += wordTokens + 1 // Include the space
	}

	if currentPart != prompt { // Add the last part if it's not just the prompt
		parts = append(parts, currentPart)
	}

	return parts, nil
}

func (c *GeminiClient) AnalyzeSubtitles(ctx context.Context, prompt, text string) (string, error) {
	parts, err := splitTextToFitModel(prompt, text, c.Config.MaxTokens)
	if err != nil {
		return "", fmt.Errorf("failed to split text: %w", err)
	}
	// 获取信号量
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer c.sem.Release(1)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second) // 设置超时时间
	defer cancel()

	var resultBuilder strings.Builder
	for _, part := range parts {
		model := c.aiClient.GenerativeModel(c.Config.ModelName)
		model.SetTemperature(c.Config.Temperature)
		model.SetTopP(c.Config.TopP)
		model.SetTopK(c.Config.TopK)
		model.SetMaxOutputTokens(c.Config.MaxTokens)
		resp, err := model.GenerateContent(ctx, genai.Text(part))
		if err != nil {
			return "", fmt.Errorf("failed to generate content for part: %w", err)
		}
		resultBuilder.WriteString(toStringResponse(resp))
	}
	if resultBuilder.Len() == 0 {
		return "", fmt.Errorf("no content generated by model %s", c.Config.ModelName)
	}
	return resultBuilder.String(), nil
}

func toStringResponse(resp *genai.GenerateContentResponse) string {
	var builder strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partString := fmt.Sprintf("%s\n\n", part)
				if partString != "" {
					builder.WriteString(partString)
				}
			}
		}
	}
	return builder.String()
}
