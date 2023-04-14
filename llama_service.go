package main

import (
	"fmt"
	"log"

	llama "github.com/wailovet/go-llama.cpp-winbin"
)

type Llama struct {
	// contains filtered or unexported fields
	model     *llama.LLama
	modelfile string
}

func (l *Llama) ModelFile() string {
	return l.modelfile
}

func (l *Llama) StartUp(modelfile string) error {
	if l.model != nil {
		l.model.Free()
	}
	var err error
	l.model, err = llama.New(modelfile, llama.SetContext(2048), llama.SetParts(-1))
	if err != nil {
		return err
	}
	l.modelfile = modelfile
	return nil
}

type Prompts struct {
	Instruct        string `json:"instruct"`
	AssistantPrefix string `json:"assistant_prefix"` //助手前缀
	UserPrefix      string `json:"user_prefix"`      //用户前缀
}

func (l *Llama) Predict(p Prompts, his ChatHistory, opts ...llama.PredictOption) (string, error) {
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
			text += fmt.Sprintln(p.AssistantPrefix, v.Content)
		} else {
			text += fmt.Sprintln(p.UserPrefix, v.Content)
		}
	}
	text += p.AssistantPrefix
	log.Println(text)
	return l.model.Predict(text, opts...)
}
