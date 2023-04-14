package main

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
