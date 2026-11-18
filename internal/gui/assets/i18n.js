// Every string the user sees goes through t(). The English text is the key
// and the fallback; each other language has its table here.
(function () {
  const zh = {
    "Workspace": "工作区",
    "Home": "主页",
    "Settings": "设置",
    "Welcome to ALGOX": "欢迎使用 ALGOX",
    "This is the sample home page. Add your own pages beside it.": "这是示例主页，可以在旁边加入自己的页面。",
    "Version": "版本",
    "Opened in": "打开方式",
    "Desktop app": "桌面应用",
    "Browser": "浏览器",
    "Browser access": "浏览器访问",
    "On": "已开启",
    "Off": "已关闭",
    "Getting started": "开始",
    "Add a page": "添加页面",
    "Put a script in assets/pages/ that calls App.page(), then list it in index.html.": "在 assets/pages/ 新建脚本并调用 App.page()，再在 index.html 引入。",
    "Add an API": "添加接口",
    "Add a route in internal/gui/api.go and call it with App.api().": "在 internal/gui/api.go 注册路由，前端用 App.api() 调用。",
    "General": "通用",
    "Theme": "主题",
    "Light, dark, or the system's": "浅色、深色或跟随系统",
    "System": "跟随系统",
    "Light": "浅色",
    "Dark": "深色",
    "Language": "语言",
    "The language of the app's text": "界面文字的语言",
    "Saved": "已保存",
    "Access": "访问",
    "Who can open this page in a browser while the app runs": "应用运行时，谁可以在浏览器中打开此页面",
    "This computer": "仅本机",
    "Local network": "局域网",
    "Port": "端口",
    "From 1024 to 65535": "1024 到 65535",
    "Apply": "应用",
    "Links": "链接",
    "Each link carries the key. Anyone with it can open ALGOX.": "链接中带有密钥，拿到链接的人都能打开 ALGOX。",
    "Copy": "复制",
    "Copied": "已复制",
    "New key": "更换密钥",
    "Browsers signed in with the old key are signed out.": "使用旧密钥登录的浏览器将被登出。",
    "Serving": "正在提供",
    "Not served": "未提供",
    "Not encrypted: other computers reach it over plain http.": "未加密：其他电脑通过明文 http 访问。",
    "Browser access is turned on and off in the desktop app.": "浏览器访问只能在桌面应用中开启或关闭。",
    "About": "关于",
  };
  const tables = { zh };
  let lang = "en";

  // setLang picks the language: "auto" follows the system.
  window.setLang = function (pick) {
    if (!pick || pick === "auto") pick = (navigator.language || "en").toLowerCase().startsWith("zh") ? "zh" : "en";
    lang = tables[pick] || pick === "en" ? pick : "en";
    document.documentElement.lang = lang;
  };
  window.t = function (s) {
    const tb = tables[lang];
    return (tb && tb[s]) || s;
  };
})();
