//go:generate rsrc -ico resource/icon.ico -manifest resource/go-llama.cpp-ui.exe.manifest -o main.syso
package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/wailovet/go-llama.cpp-winbin"
	"github.com/wailovet/gowebview2"
	"github.com/wailovet/nuwa"
)

//go:embed src
var fsbin embed.FS

var lm = Llama{}

const debug = false

const PORT = "36182" //内部端口

var wsConnMap = map[string]*websocket.Conn{} //ws连接池
var wsConnMapLock sync.RWMutex               //ws连接池锁

var adminHome, _ = os.UserHomeDir() //用户目录

func init() {
	llama.Install() //释放llama.cpp相关dll文件
}

func modelLoad(filename string) error {
	err := lm.StartUp(filename)
	if err != nil {
		return err
	}
	return nil
}

func serviceStartUp() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}

			time.Sleep(time.Second * 1)
			go serviceStartUp()
		}()

		nuwa.Config().Host = "127.0.0.1"
		nuwa.Config().Port = PORT

		nuwa.Http().HandleFunc("/model/status", func(ctx nuwa.HttpContext) {
			if lm.ModelFile() == "" {
				ctx.DisplayByData("")
				return
			}
			ctx.DisplayByData(filepath.Base(lm.ModelFile()))
		})

		nuwa.Http().HandleFunc("/model/reload", func(ctx nuwa.HttpContext) {
			basePath := ctx.REQUEST["base_path"]
			if basePath == "" {
				basePath, _ = nuwa.Helper().GetCurrentPath()
			}
			filename := ctx.ParamRequired("filename")
			filename = filepath.Join(basePath, filename)
			err := modelLoad(filename)
			ctx.CheckErrDisplayByError(err)
			ctx.DisplayByData(filename)
		})

		nuwa.Http().HandleFunc("/model/list", func(ctx nuwa.HttpContext) {
			basePath := ctx.REQUEST["base_path"]
			if basePath == "" {
				basePath, _ = nuwa.Helper().GetCurrentPath()
			}

			lists, err := os.ReadDir(basePath)
			ctx.CheckErrDisplayByError(err)

			fl := []string{}
			for _, v := range lists {
				if v.IsDir() {
					continue
				}
				if !strings.HasSuffix(v.Name(), ".bin") {
					continue
				}
				fl = append(fl, v.Name())
			}
			ctx.DisplayByData(fl)
		})

		nuwa.Http().HandleFunc("/chat/send", func(ctx nuwa.HttpContext) {
			sessionId := ctx.ParamRequired("session_id")
			repeat := cast.ToInt(ctx.REQUEST["repeat"])
			penalty := cast.ToFloat64(ctx.REQUEST["penalty"])
			temperature := cast.ToFloat64(ctx.REQUEST["temperature"])
			topP := cast.ToFloat64(ctx.REQUEST["top_p"])
			topK := cast.ToInt(ctx.REQUEST["top_k"])
			tokens := cast.ToInt(ctx.REQUEST["tokens"])
			threads := cast.ToInt(ctx.REQUEST["threads"])
			stop_words := ctx.REQUEST["stop_words"]
			if repeat == 0 {
				repeat = 8
			}
			if penalty == 0 {
				penalty = 1.5
			}
			if temperature == 0 {
				temperature = 0.7
			}
			if topP == 0 {
				topP = 0.7
			}
			if tokens == 0 {
				tokens = 1024
			}
			if threads == 0 {
				threads = 2
			}
			if stop_words == "" {
				stop_words = "##"
			}

			content := gjson.Get(ctx.BODY, "content").Raw
			var his ChatHistory
			json.Unmarshal([]byte(content), &his)
			ret := ""

			prompts := Prompts{
				Instruct:        ctx.ParamRequired("instruct"),
				AssistantPrefix: ctx.ParamRequired("assistant_prefix"),
				UserPrefix:      ctx.ParamRequired("user_prefix"),
			}
			_, err := lm.Predict(prompts, his, llama.SetTopK(topK), llama.SetRepeat(repeat), llama.SetPenalty(penalty), llama.SetTemperature(temperature), llama.SetTopP(topP), llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetStreamFn(func(outputText string) (stop bool) {

				if strings.HasSuffix(outputText, stop_words) {
					return true
				}
				ret = outputText
				if conn, ok := wsConnMap[sessionId]; ok {
					err := conn.WriteMessage(websocket.TextMessage, []byte(nuwa.Helper().JsonEncode(map[string]interface{}{
						"type":    "chat",
						"content": ret,
					})))
					if err != nil {
						log.Println(err)
						return true
					}
				}

				return false
			}))
			ctx.CheckErrDisplayByError(err)
			ctx.DisplayByData(ret)
		})

		nuwa.Http().HandleFunc("/storage/set", func(ctx nuwa.HttpContext) {
			key := ctx.ParamRequired("key")
			value := ctx.ParamRequired("value")
			err := nuwa.Bolt().SetRaw(key, value)
			ctx.CheckErrDisplayByError(err)
			ctx.DisplayByData("ok")
		})

		nuwa.Http().HandleFunc("/storage/get", func(ctx nuwa.HttpContext) {
			key := ctx.ParamRequired("key")
			value := nuwa.Bolt().GetRaw(key)
			ctx.DisplayByData(value)
		})

		nuwa.Http().HandleFunc("/storage/delete", func(ctx nuwa.HttpContext) {
			key := ctx.ParamRequired("key")
			err := nuwa.Bolt().Delete(key)
			ctx.CheckErrDisplayByError(err)
			ctx.DisplayByData("ok")
		})

		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		nuwa.Http().GetChiRouter().HandleFunc("/chat/event", func(w http.ResponseWriter, r *http.Request) {

			sessionId := r.URL.Query().Get("session_id")
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			defer conn.Close()
			wsConnMapLock.Lock()
			wsConnMap[sessionId] = conn
			wsConnMapLock.Unlock()

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}
			}
		})

		log.Println("Starting server on port " + PORT)
		err := nuwa.Http().Run()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
}

func main() {
	serviceStartUp() //启动服务

	var fc *gowebview2.AppMode
	fc, err := gowebview2.NewAppModeWithMemory(fsbin, "src")
	if err != nil {
		panic(err)
	}

	nuwa.DefaultBoltDbPath = filepath.Join(adminHome, "go-llama.cpp-ui.data") //设置bolt文件路径

	fc.Run(map[string]string{
		"width":  "1200",
		"height": "900",
		"debug":  cast.ToString(debug),
	})
}