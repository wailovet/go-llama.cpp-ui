package main

import "strings"

type ChatContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory []ChatContent

type ModelStatus string

const (
	ModelStatusNotReady ModelStatus = "not_ready"
	ModelStatusReady    ModelStatus = "ready"
)

// 中文标点符号转英文标点符号
func ChinesePunctuationToEnglishPunctuation(text string) string {
	text = strings.ReplaceAll(text, "，", ",")
	text = strings.ReplaceAll(text, "。", ".")
	text = strings.ReplaceAll(text, "！", "!")
	text = strings.ReplaceAll(text, "？", "?")
	text = strings.ReplaceAll(text, "；", ";")
	text = strings.ReplaceAll(text, "：", ":")
	text = strings.ReplaceAll(text, "“", "\"")
	text = strings.ReplaceAll(text, "”", "\"")
	text = strings.ReplaceAll(text, "‘", "'")
	text = strings.ReplaceAll(text, "’", "'")
	text = strings.ReplaceAll(text, "（", "(")
	text = strings.ReplaceAll(text, "）", ")")
	text = strings.ReplaceAll(text, "《", "<")
	text = strings.ReplaceAll(text, "》", ">")
	text = strings.ReplaceAll(text, "【", "[")
	text = strings.ReplaceAll(text, "】", "]")
	text = strings.ReplaceAll(text, "、", ",")
	text = strings.ReplaceAll(text, "—", "-")
	text = strings.ReplaceAll(text, "——", "-")
	text = strings.ReplaceAll(text, "…", "...")
	text = strings.ReplaceAll(text, "·", ".")
	return text
}
