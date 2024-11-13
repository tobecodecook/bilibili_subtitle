package config

import (
	"log"
	"os"
)

// Config holds all the API and model configurations for the project.
type Config struct {
	GeminiAPIKey      string
	OpenaiAPIKey      string
	GeminiModelConfig GeminiModelConfig
	OpenaiModelConfig OpenaiModelConfig
	Prompt            string
}

// GeminiModelConfig holds the configuration for the Gemini AI model.
type GeminiModelConfig struct {
	ModelName             string  // Name of the Gemini model
	Temperature           float32 // Temperature setting for creativity
	TopP                  float32 // TopP for controlling randomness
	MaxTokens             int32   // Max number of tokens to generate
	TopK                  int32   // TopK for controlling the diversity
	MaxConcurrentRequests int64
}

// OpenaiModelConfig holds the configuration for the OpenAI model.
type OpenaiModelConfig struct {
	ModelName   string  // Name of the OpenAI model (e.g., GPT-3, GPT-4, etc.)
	Temperature float32 // Temperature setting for creativity
	TopP        float32 // TopP for controlling randomness
	MaxTokens   int     // Max number of tokens to generate
	Timeout     int     // Timeout for requests to the OpenAI server
	Endpoint    string  // OpenAPI server URL
}

func LoadConfigValue(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Environment variable %s is not set.", key)
	}
	return val
}

func NewConfig() *Config {
	// Prompt is the text template to be used by the generative AI model.
	//Prompt1 := "The following is the content of the subtitles. Speaking intervals are separated by commas, please provide a comprehensive analysis in Chinese:"
	Prompt2 := "Here is a transcript of video subtitles, with speaking intervals separated by commas. Please conduct a thorough analysis of the themes, content, and any cultural nuances present in these subtitles. Summarize the key points and provide insights into the dialogue dynamics. All analysis and summary should be presented clearly in Chinese."
	return &Config{
		GeminiAPIKey: LoadConfigValue("gemini_api"),
		GeminiModelConfig: GeminiModelConfig{
			ModelName:             "gemini-1.5-pro-latest",
			Temperature:           0.9,
			TopP:                  0.5,
			MaxTokens:             8192,
			TopK:                  20,
			MaxConcurrentRequests: 4,
		},
		OpenaiAPIKey: LoadConfigValue("OPENAI_API_KEY"),
		OpenaiModelConfig: OpenaiModelConfig{
			ModelName:   "gpt-4o-mini", // Default OpenAI model (could be dynamically set)
			Temperature: 0.7,
			TopP:        0.9,
			MaxTokens:   16384,
			Timeout:     30,                                 // Timeout in seconds
			Endpoint:    LoadConfigValue("OPENAI_API_BASE"), // Default OpenAI endpoint
		},
		Prompt: Prompt2,
	}
}
