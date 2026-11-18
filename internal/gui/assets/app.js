// The page's shell: the sidebar, the header's breadcrumb and the open page.
// No framework: each page registers itself with App.page() and draws into
// the element it is given. The page talks to Go over /api/* (api.go).
(function () {
  const pages = [];
  const App = {
    state: null,
    view: "",

    // page registers a page. opts: id, title (English, through t()), group
    // (a sidebar group's English name), icon (an SVG's inner markup), foot
    // (true: at the sidebar's foot, as Settings is), render(el).
    page(opts) { pages.push(opts); },

    // api calls a JSON route; body given, it is a POST. Answers the JSON, or
    // throws with the route's own message.
    async api(path, body) {
      const init = body === undefined ? {} : { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) };
      const r = await fetch(path, init);
      if (!r.ok) throw new Error((await r.text()).trim() || r.statusText);
      return r.json();
    },

    // h makes an element: h("div", {class: "x", onclick}, child, "text").
    h(tag, attrs, ...kids) {
      const el = document.createElement(tag);
      for (const [k, v] of Object.entries(attrs || {})) {
        if (v == null || v === false) continue;
        if (k.startsWith("on")) el.addEventListener(k.slice(2), v);
        else if (k === "html") el.innerHTML = v;
        else el.setAttribute(k, v === true ? "" : v);
      }
      for (const k of kids.flat()) if (k != null && k !== false) el.append(k instanceof Node ? k : document.createTextNode(String(k)));
      return el;
    },

    // seg is the app's segmented switch (no native <select>).
    seg(options, value, onpick) {
      return App.h("div", { class: "seg", role: "radiogroup" }, options.map(([v, label]) =>
        App.h("button", { type: "button", role: "radio", "aria-checked": String(v === value), class: v === value ? "on" : null, onclick: () => v !== value && onpick(v) }, t(label))));
    },

    toast(msg) {
      const el = document.getElementById("toast");
      el.textContent = msg;
      el.hidden = false;
      clearTimeout(App._toast);
      App._toast = setTimeout(() => { el.hidden = true; }, 1800);
    },

    // applyPrefs puts the theme and language of the saved settings in place.
    applyPrefs() {
      const s = App.state.settings;
      if (s.theme === "auto") delete document.documentElement.dataset.theme;
      else document.documentElement.dataset.theme = s.theme;
      setLang(s.lang);
    },

    show(id) {
      const p = pages.find(x => x.id === id) || pages[0];
      App.view = p.id;
      try { localStorage.setItem("algox.lastView", p.id); } catch (_) { }
      if (location.hash !== "#/" + p.id) history.replaceState(null, "", "#/" + p.id);
      App.render();
    },

    render() {
      renderNav();
      const p = pages.find(x => x.id === App.view);
      document.getElementById("crumb").replaceChildren(App.h("span", null, t(p.title)));
      document.title = "ALGOX · " + t(p.title);
      const el = document.getElementById("page");
      el.replaceChildren();
      p.render(el);
    },

    async start() {
      document.body.classList.toggle("mac", /Mac/.test(navigator.platform));
      App.state = await App.api("/api/state");
      document.body.classList.toggle("web", App.state.web);
      document.body.classList.toggle("window", !App.state.web);
      App.applyPrefs();
      window.addEventListener("hashchange", () => App.show(location.hash.slice(2)));
      let first = location.hash.slice(2);
      if (!first) try { first = localStorage.getItem("algox.lastView") || ""; } catch (_) { }
      App.show(first);
    },
  };

  function navButton(p) {
    return App.h("button", {
      type: "button", class: "nav-i" + (p.id === App.view ? " on" : ""), "aria-current": p.id === App.view ? "page" : null,
      onclick: () => App.show(p.id),
    }, App.h("i", { class: "nav-ic", "aria-hidden": "true", html: '<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">' + (p.icon || "") + "</svg>" }),
      App.h("span", null, t(p.title)));
  }

  function renderNav() {
    const nav = document.getElementById("nav");
    const kids = [];
    let group = null;
    for (const p of pages.filter(x => !x.foot)) {
      if (p.group && p.group !== group) kids.push(App.h("div", { class: "nav-gl" }, t(p.group)));
      group = p.group;
      kids.push(navButton(p));
    }
    nav.replaceChildren(...kids);
    document.getElementById("navFoot").replaceChildren(...pages.filter(x => x.foot).map(navButton));
  }

  window.App = App;
})();
