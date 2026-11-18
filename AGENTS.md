# Notes for coding agents

## Subsystem design and review

Before changing a documented subsystem, read its
[design reference](docs/subsystems/README.md) and verify the linked source.
When preparing or reviewing a change, follow the same page's rules for
updating the reference and describing the change.

## Code standards

Write and review changes against the [code standards](docs/code-standards.md).

## Lessons

Before changing code, read [LESSONS.md](LESSONS.md): what earlier work got
wrong and the rule that would have caught it. Add an entry when a review or
a bug finds such a mistake.

## Layout

| Path | What it holds |
| --- | --- |
| `main.go` | Where the binary starts; calls `internal/cli.Main`. Nothing else. |
| `internal/cli` | The commands: `algox` / `algox app` (desktop), `algox web`, `algox version`. |
| `internal/gui` | The Wails app (`app.go`), the page's handler and JSON API (`api.go`), Browser access (`apppage.go`, `guard.go`), `algox web` (`web.go`). |
| `internal/gui/assets` | The page: plain HTML/CSS/JS, embedded, no bundler. `pages/` holds one script per page. |
| `internal/settings` | `settings.json`. |
| `internal/appdir` | algox's data folder and atomic file writes. |
| `build/` | Wails packaging per platform (icons, manifests, installers). |
| `docs/` | Code standards and subsystem references. |

## Commands

```
go test ./...                    # every test
go build -o bin/algox.exe .      # the desktop app (Windows)
go build -tags nogui -o bin/algox .   # no desktop: `algox web` only
./bin/algox web --port 18420     # the page in a browser
wails3 dev                       # desktop app, rebuilt on change
wails3 package                   # installer for this OS
```

Data lives in `$ALGOX_HOME`, else the user's config folder's `algox`.
Tests set `ALGOX_HOME` to a temp folder; never touch the real one.
