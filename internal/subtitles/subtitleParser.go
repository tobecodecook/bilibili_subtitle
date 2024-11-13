package subtitles

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SubtitleParser 是一个通用的字幕解析器接口
type SubtitleParser interface {
	Parse(filePath string, fileData []byte) (string, error)
}

// SRTSubtitleParser 实现了 SubtitleParser 接口，专门处理 SRT 或 TXT 格式
type SRTSubtitleParser struct{}

type SubtitleContent struct {
	From     float64 `json:"from"`
	To       float64 `json:"to"`
	Sid      int     `json:"sid"`
	Location int     `json:"location"`
	Content  string  `json:"content"`
	Music    float64 `json:"music"`
}

// 旧 JSON 格式的字幕解析器
type OldJSONSubtitleParser struct{}

type OldSubtitleFormat []SubtitleContent

// 新 JSON 格式的字幕解析器
type NewJSONSubtitleParser struct{}

type NewSubtitleFormat struct {
	FontSize        float64           `json:"font_size"`
	FontColor       string            `json:"font_color"`
	BackgroundAlpha float64           `json:"background_alpha"`
	BackgroundColor string            `json:"background_color"`
	Stroke          string            `json:"Stroke"`
	Type            string            `json:"type"`
	Lang            string            `json:"lang"`
	Version         string            `json:"version"`
	Body            []SubtitleContent `json:"body"`
}

// Parse 解析 SRT/TXT 格式字幕文件
func (p *SRTSubtitleParser) Parse(filePath string, fileData []byte) (string, error) {
	// 使用 bufio.Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(bytes.NewReader(fileData))
	var paragraph strings.Builder

	// 逐行处理文件
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过时间戳行、序号行和空行
		if line == "" || strings.Contains(line, "-->") || allDigits(line) {
			continue
		}

		paragraph.WriteString(line + ", ")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading subtitle file: %v", err)
	}

	return paragraph.String(), nil
}

// Parse 解析旧 JSON 格式字幕文件
func (p *OldJSONSubtitleParser) Parse(filePath string, fileData []byte) (string, error) {
	// 尝试解析为旧格式的 JSON
	var subtitles []SubtitleContent
	if err := json.Unmarshal(fileData, &subtitles); err != nil {
		return "", fmt.Errorf("error decoding JSON in file '%s': %v", filePath, err)
	}

	// 使用 strings.Builder 高效拼接字幕内容
	var paragraph strings.Builder
	for _, subtitle := range subtitles {
		paragraph.WriteString(subtitle.Content + ", ")
	}

	// 返回拼接后的字幕文本
	return paragraph.String(), nil
}

// Parse 解析新 JSON 格式字幕文件
func (p *NewJSONSubtitleParser) Parse(filePath string, fileData []byte) (string, error) {
	// 尝试解析为新的 JSON 格式
	var format NewSubtitleFormat
	if err := json.Unmarshal(fileData, &format); err != nil {
		return "", fmt.Errorf("error decoding JSON in file '%s': %v", filePath, err)
	}

	// 使用 strings.Builder 高效拼接字幕内容
	var paragraph strings.Builder
	for _, subtitle := range format.Body {
		paragraph.WriteString(subtitle.Content + ",")
	}

	// 返回拼接后的字幕文本
	return paragraph.String(), nil
}

// NewSubtitleParser 创建一个合适的字幕解析器
func NewSubtitleParser(filePath string) (SubtitleParser, []byte, error) {
	// 根据文件后缀名来判断选择哪个解析器
	if strings.HasSuffix(filePath, ".srt") || strings.HasSuffix(filePath, ".txt") {
		// 读取文件数据
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening subtitle file: %v", err)
		}
		// 返回 SRT 解析器并附带文件内容
		return &SRTSubtitleParser{}, fileData, nil

	} else if strings.HasSuffix(filePath, ".json") {
		// 尝试一次性读取整个文件
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening JSON subtitle file: %v", err)
		}

		// 尝试解码为旧的 JSON 格式
		var oldSubtitles []SubtitleContent
		if err := json.Unmarshal(fileData, &oldSubtitles); err == nil {
			// 如果解码成功为旧格式，则返回旧格式解析器
			return &OldJSONSubtitleParser{}, fileData, nil
		}

		// 尝试解码为新的 JSON 格式
		var newSubtitles NewSubtitleFormat
		if err := json.Unmarshal(fileData, &newSubtitles); err != nil {
			return nil, nil, fmt.Errorf("error decoding JSON: %v", err)
		}

		// 如果解码成功，说明是新的 JSON 格式
		return &NewJSONSubtitleParser{}, fileData, nil
	}

	return nil, nil, fmt.Errorf("unsupported subtitle file format")
}

// ParseSubtitleFile 封装了从文件路径到解析后的文本的所有操作
func ParseSubtitleFile(filePath string) (string, error) {
	// 使用工厂函数获取适当的字幕解析器以及文件内容
	parser, fileData, err := NewSubtitleParser(filePath)
	if err != nil {
		return "", err
	}

	// 使用解析器解析字幕文件内容
	parsedText, err := parser.Parse(filePath, fileData)
	if err != nil {
		return "", err
	}

	return parsedText, nil
}
