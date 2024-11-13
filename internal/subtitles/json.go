package subtitles

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Subtitle structure to parse JSON subtitle files
type Subtitle struct {
	Body []struct {
		Content string `json:"content"`
	} `json:"body"`
}

// HandleJSON processes the JSON subtitle file and returns the combined text
func HandleJSON(filePath string) (string, string, error) {
	// Open the file and read its contents
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("error reading JSON file: %v", err)
	}

	// Parse JSON content
	var subtitles Subtitle
	if err := json.Unmarshal(jsonData, &subtitles); err != nil {
		return "", "", fmt.Errorf("error parsing JSON: %v", err)
	}

	// Combine text content
	combinedText := combineJsonSubtitles(subtitles)

	// Output file path and combined text
	newFileName := strings.TrimSuffix(filePath, ".json") + ".md"
	return newFileName, combinedText, nil
}

// combineJsonSubtitles combines the content of the JSON subtitles into a single string
func combineJsonSubtitles(data Subtitle) string {
	var texts []string
	for _, item := range data.Body {
		texts = append(texts, item.Content)
	}
	return strings.Join(texts, ", ")
}
