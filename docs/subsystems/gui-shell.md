# GUI app shell

ALGOX shows one web page in two places: the desktop app's window (Wails v3)
and a browser. Both are served by the same `Handler` from embedded assets;
the page talks to Go over a JSON API, with no bindings and no bundler.

## Responsibilities and sources of truth

| Part | Responsibility | Source |
| --- | --- | --- |
| Commands | `algox` / `algox app` runs `gui.Run`; `algox web [--host] [--port]` runs `gui.ServeWeb`. A `-tags nogui` build has no desktop app and says so. | [`main.go`](../../internal/cli/main.go), [`app_off.go`](../../internal/gui/app_off.go) |
| Desktop app | `Run` makes one window that loads `/` from a `Handler` with `Web` false, in-process (no port, no key). It starts Browser access (`AppPage.Start`) and stops it on shutdown. | [`app.go`](../../internal/gui/app.go) |
| API and assets | `NewHandler` serves `assets/` and `GET /api/state`, `POST /api/settings`. `GET/POST /api/web-page` exist only on the window's handler, so a browser can't read the key or turn its own access off. | [`api.go`](../../internal/gui/api.go) |
| Browser access | Settings › Browser access picks Off, This computer (`127.0.0.1`) or Local network (every address) and a port (1024–65535). Mode, port and key are kept in `web-page.json` (0600, beside `settings.json`, never in it). The key is made the first time it's turned on and kept across restarts; New key replaces it and signs every browser out (`webGuard.SetKey`). A file that can't be read leaves it off and says why; a taken port is shown as the problem. Quitting the app stops it. | [`apppage.go`](../../internal/gui/apppage.go) |
| Sign-in | `webGuard`: the link's `?k=` sets an HttpOnly, SameSite=Lax cookie (a hash of the key) and redirects without the key; later requests need the cookie. Cross-origin writes are refused (`http.CrossOriginProtection`), and any request whose `Sec-Fetch-Site` is `cross-site` or `same-site`. | [`guard.go`](../../internal/gui/guard.go) |
| `algox web` | The same page and guard with no desktop app; a key per run, or `ALGOX_WEB_KEY`. Defaults to `127.0.0.1:18420`. | [`web.go`](../../internal/gui/web.go) |
| The page | `app.js` is the shell: the sidebar (pages in groups, Settings at the foot), the breadcrumb header, and the open page. Pages register with `App.page({id, title, group, icon, foot, render})`. The open page is in the hash (`#/settings`) and in localStorage `algox.lastView`. `body.web` / `body.window` and `body.mac` mark where it runs. Strings go through `t()` (`i18n.js`, zh and en). Under 640px the sidebar becomes a bar on top. | [`app.js`](../../internal/gui/assets/app.js), [`app.css`](../../internal/gui/assets/app.css), [`pages/`](../../internal/gui/assets/pages) |
| Settings | `settings.json`: theme (auto/light/dark), language (auto/zh/en). `Check` refuses other values; a refused save writes nothing. | [`settings.go`](../../internal/settings/settings.go) |

## Runtime path

1. `algox` → `cli.Main` → `gui.Run`: builds the browser `Handler`, wraps it in `AppPage`, builds the window's `Handler`, starts Wails.
2. The page loads `/api/state` and draws the sidebar and the page the hash or `algox.lastView` names (else the first).
3. A setting is saved with `POST /api/settings`, which answers the new state.
4. Browser access is applied with `POST /api/web-page`, which saves the file, then restarts the listener (mode or port changed) or swaps the key.

## Verification

`go test ./...`: `TestGuardLetsInOnlyWithTheKey`,
`TestWebPageRoutesAreOnTheWindowOnly`, `TestAppPageServesWithItsKeyAndTurnsOff`,
`TestAppPageTurnsDownBadChoices`, `TestAppPageUnreadFileServesNothing`,
`TestSaveLoadAndRefuse`.
