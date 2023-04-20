//go:build !cuda

package main

import llama "github.com/wailovet/go-llama.cpp-winbin"

func init() {
	llama.Install() //释放llama.cpp相关dll文件
}
