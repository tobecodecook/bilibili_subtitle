package summarization

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveSubtitleToFile 保存原始文本和生成文本到指定文件
func SaveSubtitleToFile(filePath, parsedText, result string) error {
	// 获取文件名和文件后缀
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)

	// 生成文件名
	originalFileName := strings.TrimSuffix(fileName, ext) + "original.md"
	analysisResultFileName := strings.TrimSuffix(fileName, ext) + "analysis.md"

	// 生成文件路径
	originalFilePath := filepath.Join(filepath.Dir(filePath), originalFileName)
	analysisResultFilePath := filepath.Join(filepath.Dir(filePath), analysisResultFileName)

	// 写入原始文本到 original.md 文件
	err := writeTextToFile(originalFilePath, parsedText)
	if err != nil {
		return fmt.Errorf("error writing original text to file %s: %w", originalFilePath, err)
	}

	// 写入原始文本和生成文本到 analysis.md 文件
	err = writeAnalysisToFile(analysisResultFilePath, parsedText, result)
	if err != nil {
		return fmt.Errorf("error writing analysis result to file %s: %w", analysisResultFilePath, err)
	}

	return nil
}

// writeTextToFile 写入文本到文件
func writeTextToFile(filePath, text string) error {
	// 创建并打开文件，如果文件已存在则覆盖
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	// 使用缓冲写入器提高写入性能
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入文本内容
	_, err = writer.WriteString(text)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", filePath, err)
	}

	// 确保所有缓冲区的内容被写入文件
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file %s: %w", filePath, err)
	}

	return nil
}

// writeAnalysisToFile 写入分析结果文件
func writeAnalysisToFile(filePath, parsedText, result string) error {
	// 创建并打开文件，如果文件已存在则覆盖
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	// 使用缓冲写入器提高写入性能
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入原始文本标题和内容
	_, err = writer.WriteString("## 原始文本：\n\n")
	if err != nil {
		return fmt.Errorf("error writing original text header to file %s: %w", filePath, err)
	}

	_, err = writer.WriteString(parsedText)
	if err != nil {
		return fmt.Errorf("error writing original text to file %s: %w", filePath, err)
	}

	// 写入生成文本标题
	_, err = writer.WriteString("\n\n## 生成文本：\n\n")
	if err != nil {
		return fmt.Errorf("error writing generated text header to file %s: %w", filePath, err)
	}

	// 写入生成的文本
	_, err = writer.WriteString(result)
	if err != nil {
		return fmt.Errorf("error writing generated text to file %s: %w", filePath, err)
	}

	// 确保所有缓冲区的内容被写入文件
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing file %s: %w", filePath, err)
	}

	return nil
}
