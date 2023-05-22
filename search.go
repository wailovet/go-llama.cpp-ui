package main

import (
	"log"

	"github.com/wailovet/gosearch"
	"github.com/wailovet/nuwa"
)

var msearch *gosearch.BingSearch
var segBlock = gosearch.DefaultSegBlock

func searchServiceStartUp() {
	if !nuwa.Helper().PathExists("bert.bin") || !nuwa.Helper().PathExists("bert-tokenizer.json") {
		log.Println("bert.bin or bert-tokenizer.json not exists")
		return
	}
	gosearch.Install(
		"bert.bin",
		"bert-tokenizer.json",
	)

	msearch = gosearch.NewBingSearch()
	msearch.SetWebdriver(mwebdriver)
}

func search(query string, blockQuery string, day ...int) string {
	if msearch == nil {
		return "Search function is not turned on"
	}
	threshold := float32(3)
	text := msearch.Search(query, threshold, 5, day...)
	if text == "" && len(day) > 0 && day[0] == 0 {
		text = msearch.Search(query, threshold, 5, 5)
	}

	if text == "" {
		return "no information available"
	}
	// return string([]rune(text)[:600])
	// log.Println("搜索结果:", text)
	return segBlock.Search(text, blockQuery, threshold, 500) //限制1000个字
}
