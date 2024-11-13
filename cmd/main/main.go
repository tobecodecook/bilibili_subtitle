package main

import (
	"bilibili_subtitle/internal/api"
	"bilibili_subtitle/internal/config"
	"bilibili_subtitle/internal/subtitles"
	"bilibili_subtitle/internal/summarization"
	"bilibili_subtitle/internal/utils"
	"context"
	"fmt"
	"github.com/sqweek/dialog"
	"log"
	"path/filepath"
)

func main() {
	//Set proxy (from utils)
	if err := utils.SetProxy(); err != nil {
		log.Fatal("Failed to set proxy:", err)
	}

	filePath, err := openFileDialog()
	handleError(err, "Failed to open file dialog")

	if filePath == "" {
		log.Fatal("No file selected.")
		return
	}

	cfg := config.NewConfig()
	clientChoice := "gemini"
	//clientChoice := "openai"

	err = processSubtitles(filePath, clientChoice, cfg)
	handleError(err, "Error processing subtitles")

	dir := filepath.Dir(filePath)
	err = utils.OpenDirectory(dir)
	handleError(err, fmt.Sprintf("Failed to open directory %s", dir))

}

func openFileDialog() (string, error) {
	return dialog.File().Filter("JSON ,SRT and txt files", "json", "srt", "txt").Load()
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func processSubtitles(filePath string, clientChoice string, cfg *config.Config) error {
	// 解析字幕文件
	parsedText, err := subtitles.ParseSubtitleFile(filePath)
	if err != nil {
		return err
	}

	// 执行字幕分析
	ctx := context.Background()
	result, err := api.AnalyzeWithFallback(ctx, clientChoice, cfg, parsedText)
	if err != nil {
		return err
	}

	// 保存分析结果
	err = summarization.SaveSubtitleToFile(filePath, parsedText, result)
	if err != nil {
		return err
	}

	return nil
}
