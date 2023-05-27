package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
	"github.com/wailovet/nuwa"
)

type ExternalApiService struct {
	Url string
}

// Free implements NLP
func (e *ExternalApiService) Free() {

}

// IsReady implements NLP
func (e *ExternalApiService) IsReady() bool {
	return true
}

// ModelFile implements NLP
func (e *ExternalApiService) ModelFile() string {
	return e.Url
}

// Predict implements NLP
func (e *ExternalApiService) Predict(p Prompts, his ChatHistory, opts *PredictOption) (string, error) {
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

	req := gorequest.New()

	if !strings.HasSuffix(e.Url, "/") {
		e.Url += "/"
	}
	sendUrl := e.Url + "completion"

	postJson := map[string]interface{}{
		"prompt":      text,
		"batch_size":  64,
		"as_loop":     true,
		"n_keep":      -1,
		"interactive": true,
		"stop":        []string{"\n### Human:"},
	}
	if opts != nil {
		if opts.BatchSize > 0 {
			postJson["batch_size"] = opts.BatchSize
		}
		if opts.Penalty > 0 {
			postJson["penalty"] = opts.Penalty
		}
		if opts.Temperature > 0 {
			postJson["temperature"] = opts.Temperature
		}
		if opts.TopP > 0 {
			postJson["top_p"] = opts.TopP
		}
		if opts.TopK > 0 {
			postJson["top_k"] = opts.TopK
		}
		if opts.Tokens > 0 {
			postJson["n_predict"] = opts.Tokens
		}
		if opts.Threads > 0 {
			postJson["threads"] = opts.Threads
		}
		if opts.Stop != nil {
			postJson["stop"] = opts.Stop
		}
	} else {
		return "", fmt.Errorf("opts is nil")
	}

	req.Post(sendUrl).Send(nuwa.Helper().JsonEncode(postJson)).End()

	nextTokenUrl := e.Url + "next-token"

	message := ""
	isStop := false
	for {
		var errs []error
		var result string
		if isStop {
			_, result, errs = req.Get(nextTokenUrl + "?stop=true").End()
		} else {

			_, result, errs = req.Get(nextTokenUrl).End()
		}
		if errs != nil {
			return "", errs[0]
		}
		// message += result.data.content;
		// if (result.data.stop) {
		//     console.log("Completed");
		//     // make sure to add the completion to the prompt.
		//     prompt += `### Assistant: ${message}`;
		//     break;
		// }
		message += gjson.Get(result, "content").String()
		if gjson.Get(result, "stop").Bool() {
			break
		}

		log.Println("web.result:", result)

		if !isStop {
			isStop = opts.StreamFn(message)
		}
	}

	return message, nil
}

// StartUp implements NLP
func (e *ExternalApiService) StartUp(modelfile string) error {
	e.Url = modelfile
	return nil
}
