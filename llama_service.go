package main

import (
	"fmt"
	"log"

	llama "github.com/wailovet/go-llama.cpp-winbin"
)

type Llama struct {
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
	l.model, err = llama.New(modelfile, llama.SetContext(2048), llama.SetParts(-1), llama.SetNGPULayers(0))
	if err != nil {
		return err
	}
	l.modelfile = modelfile
	return nil
}

func (l *Llama) Free() {
	if l.model != nil {
		l.model.Free()
		l.model = nil
		l.modelfile = ""
	}
}

func (l *Llama) Predict(p Prompts, his ChatHistory, opts *PredictOption) (string, error) {
	if l.ModelFile() == "" {
		return "", fmt.Errorf("model_not_loaded")
	}
	text := fmt.Sprintln(p.Instruct)

	if his == nil || len(his) == 0 {

	} else {
		if his[len(his)-1].Role == "assistant" {
			if len(his)-1 > 0 {
				his = his[:len(his)-1]
			} else {
				return "", fmt.Errorf("history is nil")
			}
		}

		for _, v := range his {
			if v.Role == "assistant" {
				text += fmt.Sprintln(p.AssistantPrefix, v.Content)
			} else {
				text += fmt.Sprintln(p.UserPrefix, v.Content)
			}
		}
		text += p.AssistantPrefix
	}

	log.Println(text)

	var mopts []llama.PredictOption
	if opts != nil {
		if opts.BatchSize > 0 {
			mopts = append(mopts, llama.SetBatchSize(opts.BatchSize))
		}
		if opts.Penalty > 0 {
			mopts = append(mopts, llama.SetPenalty(opts.Penalty))
		}
		if opts.Temperature > 0 {
			mopts = append(mopts, llama.SetTemperature(opts.Temperature))
		}
		if opts.TopP > 0 {
			mopts = append(mopts, llama.SetTopP(opts.TopP))
		}
		if opts.Tokens > 0 {
			mopts = append(mopts, llama.SetTokens(opts.Tokens))
		}
		if opts.Threads > 0 {
			mopts = append(mopts, llama.SetThreads(opts.Threads))
		}
		if opts.StreamFn != nil {
			mopts = append(mopts, llama.SetStreamFn(opts.StreamFn))
		}

	}

	return l.model.Predict(text, mopts...)
}

func (l *Llama) IsReady() bool {
	return l.model != nil
}
