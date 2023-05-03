//go:build !cuda

package main

import (
	llama "github.com/wailovet/go-llama.cpp-winbin"
	"github.com/wailovet/go-rwkv.cpp-winbin/gorwkv"
)

func init() {
	llama.Install()  //释放llama.cpp相关dll文件
	gorwkv.Install() //释放gorwkv.cpp相关dll文件
}
