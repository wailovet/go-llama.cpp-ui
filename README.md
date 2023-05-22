# go-llama.cpp-ui 
为llama.cpp编写的UI操作界面,在win上可以快速体验llama.cpp的功能

## 更新
+ 20230523: 更新llama.cpp到最新版本,修复了一些bug,新增搜索模式
+ 20230503: 新增rwkv模型支持
+ 20230428: 优化cuda版本,使用大prompt时有明显加速
+ 20230427: 当相同目录下存在app文件夹使,使用app文件夹下的UI进行启动
+ 20230422: 新增翻译模式,发送时翻译为英语,接收时翻译为中文,可以用于一些中文支持较差的模型

## 演示
![演示](demo.gif)

### 翻译模式演示
https://www.bilibili.com/video/BV1Xg4y1j7iJ

### 搜索模式
在开始的指令框中填入含{$search}的指令,程序会自动搜索并返回结果
例如:
```
You are a helpful AI Assistant. 
We have provided an existing answer: 
{$search}
We have the opportunity to refine the existing answer(only if needed) with some more context below.
```
其中{$search}会被搜索结果替换


## 快速开始 
 
#### 下载主程序:
https://gitee.com/angry_cr/go-llama.cpp-ui/releases

#### 模型下载地址:
+ https://huggingface.co/Mabbs/chinese-Alpaca-lora-7b-ggml 
+ https://huggingface.co/eachadea/ggml-vicuna-13b-1.1
+ 其他llama.cpp支持的模型
 
与模型放置同一目录下,运行`go-llama.cpp-ui.exe`即可
 
## 外部UI支持
可将自定义的UI放置在`app`文件夹下,程序会自动加载`app`文件夹下的index.html文件
release中附带了一个编译自`https://github.com/ztjhz/BetterChatGPT`的版本,可以直接使用

  