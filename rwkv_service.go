package main

import (
	"fmt"
	"log"

	"github.com/wailovet/go-rwkv.cpp-winbin/gorwkv"
)

type RWKV struct {
	model     *gorwkv.RWKV
	modelfile string
}

func (l *RWKV) ModelFile() string {
	return l.modelfile
}

func (l *RWKV) StartUp(modelfile string) error {
	if l.model != nil {
		l.model.Free()
	}
	var err error
	// cpuNum := runtime.NumCPU()
	l.model, err = gorwkv.NewRWKV(modelfile, 2, true)
	if err != nil {
		return err
	}
	l.modelfile = modelfile
	return nil
}

func (l *RWKV) Free() {
	if l.model != nil {
		l.model.Free()
		l.model = nil
		l.modelfile = ""
	}
}

func (l *RWKV) Predict(p Prompts, his ChatHistory, opts *PredictOption) (string, error) {
	if l.ModelFile() == "" {
		return "", fmt.Errorf("model_not_loaded")
	}
	if his == nil {
		return "", fmt.Errorf("history is nil")
	}

	if his[len(his)-1].Role == "assistant" {
		if len(his)-1 > 0 {
			his = his[:len(his)-1]
		} else {
			return "", fmt.Errorf("history is nil")
		}
	}

	text := fmt.Sprintln(p.Instruct)
	for _, v := range his {
		if v.Role == "assistant" {
			text += p.AssistantPrefix + v.Content + "\n\n"
		} else {
			text += p.UserPrefix + v.Content + "\n\n"
		}
	}
	text += p.AssistantPrefix
	log.Println(text)

	err := l.model.GenerateWithCache(
		text,
		gorwkv.RWKV_Config{
			MaxSeqLength:      uint32(opts.MaxTokens),
			MaxTokens:         uint32(opts.MaxTokens),
			TopP:              float32(opts.TopP),
			TopK:              float32(opts.TopK),
			Temperature:       float32(opts.Temperature),
			NoRepeatNgramSize: uint32(opts.Repeat),
			Stream:            opts.StreamFn,
		},
	)
	return text, err
}

func (l *RWKV) IsReady() bool {
	return l.model != nil
}
