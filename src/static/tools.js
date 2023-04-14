function timeout(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

 
var keyEvent = {};

document.addEventListener("keydown", function (event) {
  const keyCode = event.keyCode;
  const ctrlKey = event.ctrlKey || event.metaKey; // 判断是否按下Ctrl键

  if (ctrlKey) {
    var oklist = "ACVXZ";
    console.log(String.fromCharCode(keyCode));
    if (oklist.indexOf(String.fromCharCode(keyCode)) > -1) {
      // 判断是否按下Ctrl+C组合键
      return;
    } else {
      if (keyEvent[String.fromCharCode(keyCode)]) {
        keyEvent[String.fromCharCode(keyCode)]();
      }
      event.preventDefault(); // 阻止默认行为
    }
  } else {
    //禁止F12 F5 F11 F4 F3 F2 F1
    if (keyCode <= 123 && keyCode >= 112) {
      event.preventDefault(); // 阻止默认行为
    }
  }
});
