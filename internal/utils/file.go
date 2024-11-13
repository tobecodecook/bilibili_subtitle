package utils

import (
	"os"
)

// WriteTextToFile writes the specified content to a file
func WriteTextToFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
