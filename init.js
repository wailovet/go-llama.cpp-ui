var __POST = "%s";
var __ADDR = "http://127.0.0.1:" + __POST;

window.localStorage.getItem = function (key) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  xhr.open("GET", `${__ADDR}/localStorage/getItem?key=${key}`, false);
  xhr.send();
  var data = JSON.parse(xhr.responseText);
  return data?.data;
};

window.localStorage.setItem = function (key, value) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  value = encodeURIComponent(value);
  xhr.open(`${__ADDR}/localStorage/setItem?key=${key}&value=${value}`, false);
  xhr.send();
};

window.localStorage.removeItem = function (key) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  xhr.open("GET", `${__ADDR}/localStorage/removeItem?key=${key}`, false);
  xhr.send();
};

var nlp_model = null;

function checkModel() {
  while (true) {
    if (!nlp_model || nlp_model == "null") {
      nlp_model = prompt("请输入模型路径", localStorage.getItem("nlp_model"));
      localStorage.setItem("nlp_model", nlp_model);
    } else {
      console.log("模型路径", localStorage.getItem("nlp_model"));
      break;
    }
  }
}

var _fetch = window.fetch;
window.fetch = function (url, options) {
  if (url.indexOf("/v1/chat/completions") > -1) {
    checkModel();
    if (options.body) {
      var body = JSON.parse(options.body);
      body["model"] = nlp_model;
      options.body = JSON.stringify(body);
    }
  }
  return _fetch(url, options);
};

ah.proxy({
  //请求发起前进入
  onRequest: (config, handler) => {
    if (config.url.indexOf("/v1/chat/completions") > -1) {
      checkModel();
      if (config.body) {
        var body = JSON.parse(config.body);
        body["model"] = nlp_model;
        config.body = JSON.stringify(body);
      }
    }

    handler.next(config);
  },
  //请求发生错误时进入，比如超时；注意，不包括http状态码错误，如404仍然会认为请求成功
  onError: (err, handler) => {
    console.log(err.type);
    handler.next(err);
  },
  //请求成功后进入
  onResponse: (response, handler) => {
    console.log(response.response);
    handler.next(response);
  },
});


checkModel();