// The sample home page.
App.page({
  id: "home",
  title: "Home",
  group: "Workspace",
  icon: '<path d="M2.5 7.2 8 2.5l5.5 4.7V13a.5.5 0 0 1-.5.5H10V10H6v3.5H3a.5.5 0 0 1-.5-.5Z"/>',
  render(el) {
    const { h } = App;
    const st = App.state;
    const stat = (k, v) => h("div", { class: "card stat" }, h("div", { class: "k" }, t(k)), h("div", { class: "v" }, v));
    el.append(
      h("h1", null, t("Welcome to ALGOX")),
      h("p", { class: "lead" }, t("This is the sample home page. Add your own pages beside it.")),
      h("div", { class: "cards" },
        stat("Version", st.version),
        stat("Opened in", t(st.web ? "Browser" : "Desktop app"))),
      h("div", { class: "sec" }, t("Getting started")),
      h("div", { class: "card rows" },
        row("Add a page", "Put a script in assets/pages/ that calls App.page(), then list it in index.html."),
        row("Add an API", "Add a route in internal/gui/api.go and call it with App.api().")));

    function row(n, d) {
      return h("div", { class: "row" }, h("div", { class: "l" }, h("div", { class: "n" }, t(n)), h("div", { class: "d" }, t(d))));
    }
  },
});
