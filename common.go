package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

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

func modelType(filename string) string {
	//open a file
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	// take the first byte,type is u32
	var b uint32
	err = binary.Read(f, binary.LittleEndian, &b)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return ""
	}

	var c uint32
	err = binary.Read(f, binary.LittleEndian, &c)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return ""
	}

	switch b {
	case 0x746A6767:
		return "llama"
	case 0x67676A74:
		return "llama"
	case 0x666D6767:
		return "rwkv"
	case 0x67676D66:
		return "rwkv"
	default:
		return fmt.Sprintf("unknown:%x", b)
	}
}

type Prompts struct {
	Instruct        string `json:"instruct"`
	AssistantPrefix string `json:"assistant_prefix"` //助手前缀
	UserPrefix      string `json:"user_prefix"`      //用户前缀
}

type NLP interface {
	ModelFile() string
	StartUp(modelfile string) error
	Free()
	Predict(p Prompts, his ChatHistory, opts *PredictOption) (string, error)
	IsReady() bool
}

type PredictOption struct {
	TopK        int
	Repeat      int
	BatchSize   int
	Penalty     float64
	Temperature float64
	TopP        float64
	Tokens      int
	MaxTokens   int
	Threads     int
	StreamFn    func(outputText string) (stop bool)
}
