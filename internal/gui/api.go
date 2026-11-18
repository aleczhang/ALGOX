// Package gui serves algox's page: in the desktop window (app.go), in a
// browser through the app's Browser access (apppage.go), and from
// `algox web` (web.go). All three use Handler.
package gui

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"algox/internal/settings"
)

//go:embed assets
var embedded embed.FS

// Assets are the page's files, assets/ at their root.
func Assets() fs.FS {
	sub, _ := fs.Sub(embedded, "assets")
	return sub
}

// Handler serves the page's files and its JSON API.
type Handler struct {
	Version string
	// Web is true when the page is in a browser, not the app's own window.
	Web bool
	// Page is the app's Browser access; nil where there is none (`algox web`).
	// Its routes are on the app's window only, so a browser can neither read
	// the key nor turn the page off under itself.
	Page *AppPage

	mux *http.ServeMux
}

// NewHandler builds the routes.
func NewHandler(version string, web bool, page *AppPage) *Handler {
	h := &Handler{Version: version, Web: web, Page: page, mux: http.NewServeMux()}
	h.mux.HandleFunc("GET /api/state", h.state)
	h.mux.HandleFunc("POST /api/settings", h.saveSettings)
	if !web && page != nil {
		h.mux.HandleFunc("GET /api/web-page", h.webPage)
		h.mux.HandleFunc("POST /api/web-page", h.setWebPage)
	}
	h.mux.Handle("GET /", noStore(http.FileServerFS(Assets())))
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.mux.ServeHTTP(w, r) }

type stateJSON struct {
	Version  string            `json:"version"`
	Web      bool              `json:"web"`
	Settings settings.Settings `json:"settings"`
	// Problem says why settings.json couldn't be read, when it couldn't.
	Problem string `json:"problem,omitempty"`
}

func (h *Handler) stateNow() stateJSON {
	s, err := settings.Load()
	st := stateJSON{Version: h.Version, Web: h.Web, Settings: s}
	if err != nil {
		st.Problem = err.Error()
	}
	return st
}

func (h *Handler) state(w http.ResponseWriter, _ *http.Request) { writeJSON(w, h.stateNow()) }

func (h *Handler) saveSettings(w http.ResponseWriter, r *http.Request) {
	var s settings.Settings
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := settings.Save(s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, h.stateNow())
}

func (h *Handler) webPage(w http.ResponseWriter, _ *http.Request) { writeJSON(w, h.Page.Status()) }

func (h *Handler) setWebPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Mode   Mode `json:"mode"`
		Port   int  `json:"port"`
		NewKey bool `json:"newKey"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Page.Apply(req.Mode, req.Port, req.NewKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, h.Page.Status())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(v)
}

// noStore keeps a browser from serving an old script after an update.
func noStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
