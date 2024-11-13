package api

import (
	"bilibili_subtitle/internal/api/gemini"
	"bilibili_subtitle/internal/api/openai"
	"bilibili_subtitle/internal/config"
	"context"
	"errors"
	"log"
)

type SubtitleAnalyzer interface {
	AnalyzeSubtitles(ctx context.Context, prompt, text string) (string, error)
}

func NewSubtitleAnalyzerClient(clientChoice string, cfg *config.Config) (SubtitleAnalyzer, error) {
	switch clientChoice {
	case "gemini":
		return gemini.NewGeminiClient(cfg)
	case "openai":
		return openai.NewOpenAIClient(cfg)
	default:
		return nil, errors.New("unknown client choice")
	}
}

func AnalyzeWithFallback(ctx context.Context, clientChoice string, cfg *config.Config, parsedText string) (string, error) {
	client, err := NewSubtitleAnalyzerClient(clientChoice, cfg)
	if err != nil {
		return "", err
	}

	result, err := client.AnalyzeSubtitles(ctx, cfg.Prompt, parsedText)
	if err != nil || result == "" {
		log.Printf("Primary client (%s) failed or returned empty result: %v. Switching client.\n", clientChoice, err)
		fallbackChoice := getFallbackClientChoice(clientChoice)
		fallbackClient, err := NewSubtitleAnalyzerClient(fallbackChoice, cfg)
		if err != nil {
			return "", err
		}
		result, err = fallbackClient.AnalyzeSubtitles(ctx, cfg.Prompt, parsedText)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func getFallbackClientChoice(clientChoice string) string {
	if clientChoice == "gemini" {
		return "openai"
	}
	return "gemini"
}
