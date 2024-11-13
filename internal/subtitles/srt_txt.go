package subtitles

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//1
//0:0:0,28 --> 0:0:2,14
//平常打混双的都知道
//
//2
//0:0:2,14 --> 0:0:6,78
//最怕就是女后男前女生被按在后场动弹不得

// HandleSRT processes the SRT subtitle file and returns the combined text
func HandleSRT(filePath string) (string, string, error) {
	// Open the SRT file and read line by line
	fileData, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("error opening SRT file: %v", err)
	}
	defer fileData.Close()

	scanner := bufio.NewScanner(fileData)
	var paragraph strings.Builder

	// Process each line in the file
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines, sequence numbers, and timestamps
		if line == "" || strings.Contains(line, "-->") || allDigits(line) {
			continue
		}
		paragraph.WriteString(line + " ")
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading SRT file: %v", err)
	}

	// Return the combined text and the new file path
	combinedText := paragraph.String()
	newFileName := strings.TrimSuffix(filePath, ".srt") + ".md"
	return newFileName, combinedText, nil
}

// allDigits checks if the string contains only digits
func allDigits(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("0123456789", r) {
			return false
		}
	}
	return true
}
