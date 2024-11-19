package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

import (
	"context"
	"github.com/google/generative-ai-go/genai"
)
import "google.golang.org/api/option"

type Subtitle struct {
	Body []struct {
		Content string `json:"content"`
	} `json:"body"`
}

//var prompt1 = "Please analyze the content of video subtitles where the sentences are separated by commas. Provide a detailed summary and in-depth analysis of these subtitles, synthesizing the main themes and ideas in Chinese."

var prompt1 = "The following is the content of the video subtitles. Speaking intervals are separated by commas,Please provide a comprehensive analysis of the specified content in Chinese:"

func main() {
	err := os.Setenv("HTTP_PROXY", os.Getenv("HTTP_PROXY"))
	if err != nil {
		return
	}
	err = os.Setenv("HTTPS_PROXY", os.Getenv("HTTPS_PROXY"))
	if err != nil {
		return
	}
	// 打开文件选择对话框
	filePath, err := dialog.File().Filter("JSON,SRT and txt files", "json", "srt", "txt").Load()
	if err != nil {
		fmt.Println("Failed to open file dialog:", err)
		return
	}

	// 检查并获取文件名与扩展名
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)
	var newFilePath string
	var combinedText string

	// 根据文件类型执行不同的处理
	switch fileExt {
	case ".json":
		newFilePath, err = handleJSON(filePath, fileName, &combinedText)
		if err != nil {
			fmt.Println(err)
			return
		}
	case ".srt":
		newFilePath, err = handleSRT(filePath, fileName, &combinedText)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("Please select a JSON or SRT file.")
	}

	summarizeText(combinedText, newFilePath)

	dir := filepath.Dir(newFilePath)
	err = openDirectory(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func printResponse(combinedText string, resp *genai.GenerateContentResponse, newFilePath string) error {
	// Open the file for writing (creating if necessary)
	file, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Write the entire response to the file using a buffered writer for performance
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			log.Printf("Error flushing file: %v", err)
		}
	}(writer)

	_, err = writer.WriteString("## 原始文本：\n\n")
	if err != nil {
		return fmt.Errorf("error writing original text: %w", err)
	}

	_, err = writer.WriteString(combinedText)
	if err != nil {
		return fmt.Errorf("error writing original text: %w", err)
	}

	_, err = writer.WriteString("\n\n## 生成文本：\n\n")
	if err != nil {
		return fmt.Errorf("error writing generated text header: %w", err)
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partString := fmt.Sprintf("%s\n\n", part) // Add two newlines for Markdown formatting

				_, err = writer.WriteString(partString)
				if err != nil {
					return fmt.Errorf("error writing candidate part: %w", err)
				}
			}
		}
	}

	// Handle any potential errors during file flushing
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file: %w", err)
	}

	return nil
}

func combineJsonSubtitles(data Subtitle) string {
	var texts []string
	for _, item := range data.Body {
		texts = append(texts, item.Content)
	}
	return strings.Join(texts, ",")
}

func summarizeText(text string, newFilePath string) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	defer func(client *genai.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt1+text))
	if err != nil {
		log.Fatal(err)
	}

	if err := printResponse(text, resp, newFilePath); err != nil {
		log.Fatal(err)
	}
}

// handleJSON 处理 JSON 文件，转换为 Markdown
func handleJSON(filePath, fileName string, combinedText *string) (string, error) {
	newFileName := strings.TrimSuffix(fileName, ".json") + ".md"
	newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)
	// 读取和解析文件内容
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", fmt.Errorf("error reading file: %v", err)
	}

	var subtitles Subtitle
	if err := json.Unmarshal(jsonData, &subtitles); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	*combinedText = combineJsonSubtitles(subtitles)

	// 创建并写入新的 Markdown 文件
	originalFilename := strings.TrimSuffix(fileName, ".json") + "original.md"
	originalFilePath := filepath.Join(filepath.Dir(filePath), originalFilename)
	// 创建并写入新的 Markdown 文件
	err = os.WriteFile(originalFilePath, []byte(*combinedText), 0644)
	if err != nil {
		fmt.Println("Error writing Markdown file:", err)
		return "", fmt.Errorf("error writing Markdown file: %v", err)
	}

	return newFilePath, nil
}

// handleSRT 处理 SRT 文件，转换为 markdown
func handleSRT(filePath, fileName string, combinedText *string) (string, error) {
	fileData, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer func(fileData *os.File) {
		err := fileData.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(fileData)

	scanner := bufio.NewScanner(fileData)
	var paragraph strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		// 检查行是否为空，或者是序号，或者是时间码
		if line == "" || strings.Contains(line, "-->") || allDigits(line) {
			continue
		}
		paragraph.WriteString(line + " ")
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	*combinedText = paragraph.String()

	newFileName := strings.TrimSuffix(fileName, ".srt") + ".md"
	newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)

	originalFilename := strings.TrimSuffix(fileName, ".srt") + "original.md"
	originalFilePath := filepath.Join(filepath.Dir(filePath), originalFilename)
	// 创建并写入新的 Markdown 文件
	err = os.WriteFile(originalFilePath, []byte(*combinedText), 0644)
	if err != nil {
		fmt.Println("Error writing Markdown file:", err)
		return "", fmt.Errorf("error writing Markdown file: %v", err)
	}

	return newFilePath, nil
}

// allDigits 检查字符串是否只包含数字
func allDigits(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("0123456789", r) {
			return false
		}
	}
	return true
}

func openDirectory(path string) error {
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
