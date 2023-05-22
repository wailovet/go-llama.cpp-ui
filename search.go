package main

import (
	"github.com/wailovet/gosearch"
)

var msearch *gosearch.BingSearch
var segBlock = gosearch.DefaultSegBlock

func searchServiceStartUp() {
	gosearch.Install(
		"bert.bin",
		"bert-tokenizer.json",
	)

	msearch = gosearch.NewBingSearch()
	msearch.SetWebdriver(mwebdriver)
}

func search(query string, blockQuery string, day ...int) string {
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
