# ALGOX

A desktop app (Wails v3) whose page can also be opened in a browser.

- **Desktop:** `algox` opens the window.
- **Browser:** turn on Settings › Browser access (This computer or Local
  network), then open the link it shows. Or, with no desktop, run
  `algox web --port 18420`.

## Develop

```
go test ./...
go build -o bin/algox.exe .     # then ./bin/algox.exe
wails3 dev                      # rebuild on change
wails3 package                  # installer
```

Layout, conventions and design references: [AGENTS.md](AGENTS.md).
