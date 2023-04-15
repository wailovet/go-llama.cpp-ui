var host = `http://localhost:36182`;
var ws_host = `ws://localhost:36182`;
var app;

//json
initVue();

// Generate a pseudo-GUID by concatenating random hexadecimal.
function guid() {
  function S4() {
    return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
  }
  return (
    S4() +
    S4() +
    "-" +
    S4() +
    "-" +
    S4() +
    "-" +
    S4() +
    "-" +
    S4() +
    S4() +
    S4()
  );
}

function connectWebsocket() {
  app.session_id = guid();
  app.chatWebSocket = new WebSocket(
    `${ws_host}/chat/event?session_id=${app.session_id}`
  );
  app.chatWebSocket.onopen = function () {
    console.log("websocket 已连接");
  };
  app.chatWebSocket.onmessage = function (evt) {
    var received_msg = evt.data;
    var data = JSON.parse(received_msg);
    var type = data?.type;
    if (type == "chat") {
      var content = data?.content;
      // 过滤 \ufffd
      content = content.replace(/\ufffd/g, "");

      console.log("websocket 收到消息", content);
      sourceChatList = JSON.parse(JSON.stringify(app.sourceChatList));
      sourceChatList[sourceChatList.length - 1].content = content;
      app.sourceChatList = sourceChatList;
      app.$forceUpdate();
    }
  };
  app.chatWebSocket.onclose = function () {
    console.log("websocket 已关闭");
    setTimeout(() => {
      connectWebsocket();
    }, 1000);
  };
}

//@note loadChatHistory
async function loadChatHistory() {
  app.history = [];
  var res = await storageGetValue("chat_history");
  let data;
  if (!res) {
    data = [];
  } else {
    data = JSON.parse(res);
  }

  data.sort((a, b) => {
    return b.update_time - a.update_time;
  });

  app.history = data;
}

function initVue() {
  app = new Vue({
    el: "#app",
    data: function () {
      return {
        chatWebSocket: null,
        is_stop_generate: false,
        uuid: guid(),
        session_id: "",
        nlp_model: "",
        nlp_model_base_path: "",
        nlp_model_list: [],
        history: [],
        sourceChatList: [],
        chatList: [],
        inputText: ``,
        title: "我有什么可以帮助你的？",
        lock: false,
        setting_visible: false,
        search_title: "",
        active_screen: "chat",
        chat_editor_data: {
          index: -1,
          content: "",
        },
        predict_config: {
          temperature: 0.65,
          top_k: 10000,
          top_p: 0.7,
          penalty: 1,
          repeat: 64,
          tokens: 128,
          threads: 2,
          instruct: "### instruction: 你是一个友好的助手,你会详细的对用户的问题进行回复",
          user_prefix: "### user:",
          assistant_prefix: "### assistant:",
          stop_words: "##",
        },
      };
    },
    created: function () {},
    mounted: async function () {
      await timeout(200);
      await loadPage();
      await timeout(2000);
      //
    },
    watch: {
      sourceChatList: async function (val) {
        await timeout(100);
        this.chatList = preprocessing(app.sourceChatList);
        this.toEnd();
      },
      nlp_model_base_path: function (val) {
        storageSetValue("config.nlp_model_base_path", val);
      },
      nlp_model: function (val) {
        storageSetValue("config.nlp_model", val);
      },
      predict_config: {
        handler: function (val) {
          storageSetValue("config.predict_config", JSON.stringify(val));
        },
        deep: true,
      },
    },
    methods: {
      activeScreen: async function (screen) {
        switch (screen) {
          case "chat":
            app.active_screen = screen;
            await loadPage();
            break;
          case "setting":
            app.active_screen = screen;
            break;
        }
      },
      getOptions: function (item) {
        return [
          {
            type: "li",
            text: "删除",
            callback: function () {
              app.deleteUUID(item["uuid"]);
            },
          },
        ];
      },
      getChatRecordOptions: function (index, item) {
        let opts = [
          {
            type: "li",
            text: "编辑",
            callback: async function () {
              app.active_screen = "chat_editor";
              app.chat_editor_data = {
                index: index,
                content: app.sourceChatList[index].content,
              };
            },
          },
        ];

        if (item.role == "user") {
          opts.push({
            type: "li",
            text: "将以上记录克隆到新的窗口",
            disabled: index - 1 < 1,
            callback: async function () {
              app.uuid = guid();
              app.sourceChatList = app.sourceChatList.slice(0, index);
            },
          });
        }
        return opts;
      },
      chatEditorSave: async function () {
        let index = app.chat_editor_data.index;
        let content = app.chat_editor_data.content;
        app.sourceChatList[index].content = content;
        app.active_screen = "chat";

        await saveCurrentChatHistory();

        await loadChatHistory();
        await loadChatContent();
      },
      toEnd() {
        setTimeout(() => {
          document.querySelector(".all_chat_list_container").scrollTop =
            document.querySelector(".all_chat_list_container").scrollHeight;

          PR.prettyPrint();
        }, 100);
      },
      reload: function () {
        location.reload();
      },
      inputTextChange: function () {
        const textarea = document.querySelector("#input_text");
        textarea.style.height = "auto";
        textarea.style.height = textarea.scrollHeight - 5 + "px";
      },
      clearChatList: async function (e) {
        if (app.lock) { 
          return;
        }
        app.uuid = guid(); 
        await loadChatHistory(); 
        await loadChatContent();
      },
      loadPage: async function () {
        loadPage();
      },
      loadModel: async function () {
        await loadModel();
      },
      sendText: async function () {
        if (app.lock) {
          return;
        }
        app.lock = true;

        await sendText();
      },
      stopGenerate: async function () {
        app.chatWebSocket.close();
        app.lock = false;
        app.is_stop_generate = true;
        await loadChatContent();
      },
      setUUID: async function (uuid) {
        if (app.lock) {
          app.$message.error("请先停止当前对话");
          return;
        }
        app.uuid = uuid;
        await loadChatContent();
      },
      deleteUUID: async function (uuid) {
        var data = app.history;
        var index = data.findIndex((item) => {
          return item.uuid == uuid;
        });
        if (index > -1) {
          data.splice(index, 1);
        }
        await storageSetValue("chat_history", JSON.stringify(data));
        await loadChatHistory();
        await loadChatContent();
      },
    },
  });

  connectWebsocket();
}

async function loadPage() {
  await loadConfig();
  await checkInitModel();
  await loadChatHistory();
  await loadChatContent();
}

async function loadConfig() {
  app.nlp_model_base_path = await storageGetValue("config.nlp_model_base_path");
  app.nlp_model = await storageGetValue("config.nlp_model");
  let predict_config = await storageGetValue("config.predict_config");
  if (predict_config) {
    app.predict_config = JSON.parse(predict_config);
  }
}

async function loadModel() {
  let res = await axios.post(`${host}/model/reload`, {
    base_path: app.nlp_model_base_path,
    filename: app.nlp_model,
  });
  res = res.data;
  if (res.code == 0) {
    app.$message.success("加载成功");
    app.activeScreen("chat");
  } else {
    app.$message.error("加载失败");
  }
}

async function checkInitModel() {
  let res = await axios.get(`${host}/model/status`);
  res = res.data;
  if (!res.data || res.data == "") {
    let model_list = await axios.post(`${host}/model/list`, {
      base_path: app.nlp_model_base_path,
    });
    model_list = model_list.data;
    app.nlp_model_list = [];
    for (let item of model_list.data) {
      app.nlp_model_list.push({
        label: item,
        value: item,
      });
    }

    app.activeScreen("setting");
    return false;
  }
  app.nlp_model = res.data;
}

async function loadChatContent() {
  for (let item of app.history) {
    if (item.uuid == app.uuid) {
      app.sourceChatList = item.content;
      app.title = item.title;
      return;
    }
  }

  app.title = "我有什么可以帮助你的？";

  app.sourceChatList = [];

  app.$forceUpdate();
}

var appAlert;
async function sendText() {
  app.is_stop_generate = false;
  var content = app.inputText;
  if (content) {
    bakSourceChatList = JSON.parse(JSON.stringify(app.sourceChatList));

    app.sourceChatList.push({
      role: "user",
      content: content,
    });
    app.sourceChatList.push({
      role: "assistant",
      content: "",
    });

    app.toEnd();

    app.inputText = "";
    setTimeout(() => {
      app.inputTextChange();
    }, 100);

    var res = null;

    var post_data = JSON.parse(JSON.stringify(app.predict_config));
    post_data["session_id"] = app.session_id;
    post_data["content"] = app.sourceChatList;
    console.log(post_data);
    res = await axios.post(`${host}/chat/send`, post_data);
    res = res.data;

    if (res.code > 0) {
      app.$message.error(res.message);
      app.sourceChatList = bakSourceChatList;
      app.chatList[
        app.chatList.length - 1
      ].content = `<p style="color:red;">${res.message}</p>`;
      return;
    }

    let out_text = "";
    if (!app.is_stop_generate) {
      out_text = res?.data;
    } else {
      out_text = app.sourceChatList[app.sourceChatList.length - 1].content;
      if (out_text == "") {
        app.sourceChatList = bakSourceChatList;
        return;
      }
    }

    //去除前后空格和换行
    out_text = out_text.replace(/^\s+|\s+$/g, "");

    app.sourceChatList[app.sourceChatList.length - 1].content = out_text;
    app.$forceUpdate();

    await saveCurrentChatHistory();

    await loadChatContent();
  }

  app.lock = false;
}

async function saveCurrentChatHistory() {
  let history = app.history;
  let isExist = false;
  for (let i = 0; i < history.length; i++) {
    if (history[i].uuid == app.uuid) {
      isExist = true;
      history[i].content = app.sourceChatList;
      history[i].title = history[i].content[0].content;
      history[i].update_time = new Date().getTime();
      break;
    }
  }
  if (!isExist) {
    history.push({
      uuid: app.uuid,
      title: app.sourceChatList[0].content,
      content: app.sourceChatList,
      update_time: new Date().getTime(),
    });
  }

  // 排序
  history.sort((a, b) => {
    return b.update_time - a.update_time;
  });

  storageSetValue("chat_history", JSON.stringify(history));

  await loadChatContent();
}

function preprocessing(data) {
  data = JSON.parse(JSON.stringify(data));
  for (let i = 0; i < data.length; i++) {
    if (data[i].role == "assistant") {
      if (data[i].content) {
        data[i].content = marked.parse(data[i].content, {
          langPrefix: "language-",
        });
        data[i].content = data[i].content.replace(
          /<pre><code>/g,
          "<pre><code class='prettyprint'>"
        );
        data[i].content = data[i].content.replace(
          /<pre><code class="language/g,
          '<pre><code class="prettyprint language'
        );
        data[i].content = data[i].content.replace(
          /<a href/g,
          "<a target='_blank' href"
        );
      }
    }
  }

  return data;
}

async function storageSetValue(key, value) {
  await axios.post(`${host}/storage/set`, {
    key: key,
    value: value,
  });
}

async function storageGetValue(key) {
  let res = await axios.post(`${host}/storage/get`, {
    key: key,
  });
  res = res.data;
  return res.data;
}
