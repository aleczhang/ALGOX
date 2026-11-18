# Code standards

## What a change contains

- One change has one goal. Unrelated refactors and fixes go in their own commit.
- A change to a documented subsystem updates its reference in
  [`docs/subsystems`](subsystems/README.md) in the same commit.
- A behavior change comes with a test that fails without it.

## Tests

- `go test ./...` passes before a commit.
- Tests set `ALGOX_HOME` to `t.TempDir()`; none reads or writes the user's
  real data folder.
- A test of a network path listens on a free port (`127.0.0.1:0`), never a
  fixed one.

## Go

- Write `any`, not `interface{}`.
- `main.go` only calls `internal/cli`. Code goes in `internal/<package>`,
  one package per job.
- Read an environment variable that names a folder with `appdir.Getenv`,
  which ignores a relative path.
- Write a file of algox's with `appdir.WriteFile` (temp file and rename).
  A file holding a key is 0600 and stays out of `settings.json`.
- A comment says what the code is for in a sentence or two. Name
  functions, not line numbers.

## GUI

- The page is plain HTML/CSS/JS in `internal/gui/assets`, embedded with
  `go:embed`. No framework, no bundler, no Wails bindings: the page calls
  Go through `/api/*` with `App.api()`.
- A page is one script in `assets/pages/` that calls `App.page()`, listed
  in `index.html`.
- Every string the user sees goes through `t()`, with its Chinese in
  `i18n.js`. The English text is the key.
- Use no native `<select>`. Use `App.seg()` or the app's own menu.
- Colors come from the CSS variables in `app.css`, so both themes work.
  Don't use colored left-border stripes; mark state with a dot.
- A route that only the desktop window may use (`/api/web-page`) is
  registered only on the window's `Handler` (`Web` false).

## Commits

The subject names the area and what the user now sees:
`gui: the Settings page saves the theme at once`. The body says how it was
verified.
