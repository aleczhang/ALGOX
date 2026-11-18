// Settings: General (theme, language) and Browser access. Browser access is
// drawn in the desktop window only: its route (/api/web-page) isn't on a
// browser's API.
App.page({
  id: "settings",
  title: "Settings",
  foot: true,
  icon: '<circle cx="8" cy="8" r="2.2"/><path d="M8 1.8v1.6M8 12.6v1.6M1.8 8h1.6M12.6 8h1.6M3.6 3.6l1.1 1.1M11.3 11.3l1.1 1.1M3.6 12.4l1.1-1.1M11.3 4.7l1.1-1.1"/>',
  render(el) {
    const { h } = App;
    const row = (n, d, right) => h("div", { class: "row" },
      h("div", { class: "l" }, h("div", { class: "n" }, t(n)), d && h("div", { class: "d" }, t(d))),
      h("div", { class: "r" }, right));

    async function save(change) {
      try {
        App.state = await App.api("/api/settings", { ...App.state.settings, ...change });
        App.applyPrefs();
        App.render();
        App.toast(t("Saved"));
      } catch (e) { App.toast(e.message); }
    }

    const s = App.state.settings;
    el.append(
      h("div", { class: "sec" }, t("General")),
      h("div", { class: "card rows" },
        row("Theme", "Light, dark, or the system's", App.seg([["auto", "System"], ["light", "Light"], ["dark", "Dark"]], s.theme, v => save({ theme: v }))),
        row("Language", "The language of the app's text", App.seg([["auto", "System"], ["zh", "中文"], ["en", "English"]], s.lang, v => save({ lang: v })))),
      App.state.problem ? h("div", { class: "note bad" }, App.state.problem) : "");

    el.append(h("div", { class: "sec" }, t("Browser access")));
    const box = h("div", { class: "card rows" });
    el.append(box);
    if (App.state.web) {
      box.append(h("div", { class: "note" }, t("Browser access is turned on and off in the desktop app.")));
    } else {
      drawWebPage(box);
    }

    el.append(h("div", { class: "sec" }, t("About")),
      h("div", { class: "card rows" }, row("Version", null, h("span", null, App.state.version))));

    async function drawWebPage(box, given) {
      let st = given;
      if (!st) {
        try { st = await App.api("/api/web-page"); } catch (e) { box.replaceChildren(h("div", { class: "note bad" }, e.message)); return; }
      }
      const port = h("input", { class: "field", type: "number", min: "1024", max: "65535", value: String(st.port), "aria-label": t("Port") });
      const apply = async (mode, newKey) => {
        try {
          const next = await App.api("/api/web-page", { mode, port: Number(port.value) || 0, newKey: !!newKey });
          drawWebPage(box, next);
          App.toast(t("Saved"));
        } catch (e) { App.toast(e.message); }
      };
      const kids = [
        row("Access", "Who can open this page in a browser while the app runs",
          App.seg([["off", "Off"], ["local", "This computer"], ["lan", "Local network"]], st.mode, v => apply(v))),
        row("Port", "From 1024 to 65535", [port, h("button", { type: "button", class: "btn", onclick: () => apply(st.mode) }, t("Apply"))]),
      ];
      if (st.mode !== "off") {
        kids.push(row(st.serving ? "Serving" : "Not served", null, h("span", { class: "dot " + (st.serving ? "ok" : "warn") })));
      }
      if (st.problem) kids.push(h("div", { class: "note bad" }, st.problem));
      if (st.links.length) {
        kids.push(h("div", { class: "note" }, t("Each link carries the key. Anyone with it can open ALGOX.")));
        for (const link of st.links) {
          kids.push(h("div", { class: "link" }, h("code", null, link),
            h("button", { type: "button", class: "btn", onclick: () => copy(link) }, t("Copy"))));
        }
        kids.push(row("New key", "Browsers signed in with the old key are signed out.",
          h("button", { type: "button", class: "btn", onclick: () => apply(st.mode, true) }, t("New key"))));
      }
      if (st.mode === "lan") kids.push(h("div", { class: "note" }, t("Not encrypted: other computers reach it over plain http.")));
      box.replaceChildren(...kids);
    }

    async function copy(text) {
      try { await navigator.clipboard.writeText(text); App.toast(t("Copied")); } catch (e) { App.toast(e.message); }
    }
  },
});
