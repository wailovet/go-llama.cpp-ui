//go:build cuda

package main

import llama "github.com/wailovet/go-llama.cpp-winbin"

func init() {
	llama.InstallCuda()
}
