package config

import (
	"os"
	"testing"
)

// TestLoadConfigValueSet tests the LoadConfigValue function for a properly set environment variable.
func TestLoadConfigValueSet(t *testing.T) {
	// Set an environment variable for testing.
	expectedValue := "test-value"
	os.Setenv("TEST_ENV_VAR", expectedValue)
	defer os.Unsetenv("TEST_ENV_VAR") // Clean up after the test.

	// Call the function with the known environment variable.
	if val := LoadConfigValue("TEST_ENV_VAR"); val != expectedValue {
		t.Errorf("LoadConfigValue returned %s, want %s", val, expectedValue)
	}
}

// TestLoadConfigValueUnset tests the LoadConfigValue function for an unset environment variable.
// This test expects a log.Fatal to be called, which we simulate here by recovering from a panic.
func TestLoadConfigValueUnset(t *testing.T) {
	// Ensure the environment variable is unset.
	os.Unsetenv("TEST_ENV_VAR")

	// Use a deferred function to recover from the expected panic due to log.Fatal inside LoadConfigValue.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to unset environment variable, but no panic occurred")
		}
	}()

	// This should cause a panic as the environment variable is not set.
	LoadConfigValue("TEST_ENV_VAR")
}

// TestNewConfig tests the NewConfig function to ensure it loads all configurations correctly.
func TestNewConfig(t *testing.T) {
	// Setup environment
	os.Setenv("GEMINI_API_KEY", "fake-gemini-key")
	os.Setenv("OPENAI_API_KEY", "fake-openai-key")
	os.Setenv("GEMINI_MODEL_NAME", "gemini-1.5-pro-latest")
	os.Setenv("OPENAI_MODEL_NAME", "gpt-4o-mini")
	os.Setenv("OPENAI_API_BASE", "https://api.fakeurl.com")

	// Expected to run without causing a panic
	config := NewConfig()

	// Verify the configuration
	if config.GeminiAPIKey != "fake-gemini-key" {
		t.Errorf("Incorrect Gemini API Key, got: %s, want: %s", config.GeminiAPIKey, "fake-gemini-key")
	}
	if config.OpenaiAPIKey != "fake-openai-key" {
		t.Errorf("Incorrect OpenAI API Key, got: %s, want: %s", config.OpenaiAPIKey, "fake-openai-key")
	}

	// Clean up environment
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("GEMINI_MODEL_NAME")
	os.Unsetenv("OPENAI_MODEL_NAME")
	os.Unsetenv("OPENAI_API_BASE")
}
