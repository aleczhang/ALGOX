package gui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"algox/internal/appdir"
)

// Mode is who the app's Browser access lets reach the page.
type Mode string

const (
	ModeOff   Mode = "off"
	ModeLocal Mode = "local" // 127.0.0.1 only: this computer, or a tunnel pointed at it
	ModeLAN   Mode = "lan"   // every address: other computers on the network
)

// DefaultPort is the Browser access port until the user picks another.
const DefaultPort = 18420

const pageFile = "web-page.json"

// pageConfig is web-page.json, beside settings.json and never in it, so the
// key stays on this computer.
type pageConfig struct {
	Mode Mode   `json:"mode"`
	Port int    `json:"port"`
	Key  string `json:"key,omitempty"`
}

func (c pageConfig) check() error {
	switch c.Mode {
	case ModeOff, ModeLocal, ModeLAN:
	default:
		return fmt.Errorf("unknown mode %q", c.Mode)
	}
	if c.Port < 1024 || c.Port > 65535 {
		return fmt.Errorf("port %d: pick one from 1024 to 65535", c.Port)
	}
	return nil
}

// AppPage serves the app's page to browsers while the desktop app runs,
// turned on and off from the app's Settings.
type AppPage struct {
	mu      sync.Mutex
	cfg     pageConfig
	page    http.Handler // the page as a browser sees it (Handler with Web)
	guard   *webGuard
	srv     *http.Server
	problem string // why the page isn't served though it should be
}

// NewAppPage reads web-page.json. A file that can't be read leaves the page
// off and says why.
func NewAppPage(page http.Handler) *AppPage {
	p := &AppPage{page: page, cfg: pageConfig{Mode: ModeOff, Port: DefaultPort}}
	cfg, err := loadPageConfig()
	if err != nil {
		p.problem = err.Error()
		return p
	}
	p.cfg = cfg
	return p
}

func loadPageConfig() (pageConfig, error) {
	cfg := pageConfig{Mode: ModeOff, Port: DefaultPort}
	path, err := appdir.File(pageFile)
	if err != nil {
		return cfg, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err == nil {
		err = json.Unmarshal(b, &cfg)
	}
	if err == nil {
		err = cfg.check()
	}
	if err != nil {
		return pageConfig{Mode: ModeOff, Port: DefaultPort}, fmt.Errorf("%s: %w; browser access is off", path, err)
	}
	return cfg, nil
}

func savePageConfig(cfg pageConfig) error {
	path, err := appdir.File(pageFile)
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return appdir.WriteFile(path, b, 0o600)
}

// Start serves the page as web-page.json says; the app calls it once.
func (p *AppPage) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.serveLocked()
}

// Stop stops serving, as the app quits.
func (p *AppPage) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopLocked()
}

// Apply saves a whole choice and serves by it. newKey replaces the key,
// letting go every browser signed in with the old one.
func (p *AppPage) Apply(mode Mode, port int, newKey bool) error {
	if port == 0 {
		port = DefaultPort
	}
	cfg := pageConfig{Mode: mode, Port: port}
	if err := cfg.check(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	cfg.Key = p.cfg.Key
	if cfg.Key == "" && mode != ModeOff || newKey {
		cfg.Key = NewKey()
	}
	if err := savePageConfig(cfg); err != nil {
		return err
	}
	moved := cfg.Mode != p.cfg.Mode || cfg.Port != p.cfg.Port
	p.cfg = cfg
	if !moved && p.srv != nil {
		p.guard.SetKey(cfg.Key)
		return nil
	}
	p.stopLocked()
	p.serveLocked()
	return nil
}

func (p *AppPage) serveLocked() {
	if p.cfg.Mode == ModeOff {
		return
	}
	host := "127.0.0.1"
	if p.cfg.Mode == ModeLAN {
		host = ""
	}
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(p.cfg.Port)))
	if err != nil {
		p.problem = fmt.Sprintf("port %d is taken: %v", p.cfg.Port, err)
		return
	}
	p.problem = ""
	p.guard = newWebGuard(p.cfg.Key, p.cfg.Port, p.page)
	p.srv = &http.Server{Handler: p.guard, ReadHeaderTimeout: 10 * time.Second}
	go func(srv *http.Server) { _ = srv.Serve(ln) }(p.srv)
}

func (p *AppPage) stopLocked() {
	if p.srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = p.srv.Shutdown(ctx)
	p.srv, p.guard = nil, nil
}

// PageStatus is what Settings › Browser access shows.
type PageStatus struct {
	Mode    Mode     `json:"mode"`
	Port    int      `json:"port"`
	Serving bool     `json:"serving"`
	Links   []string `json:"links"`
	Problem string   `json:"problem,omitempty"`
}

// Status says how the page is served, with the links that carry the key.
func (p *AppPage) Status() PageStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	st := PageStatus{Mode: p.cfg.Mode, Port: p.cfg.Port, Serving: p.srv != nil, Problem: p.problem, Links: []string{}}
	if p.srv == nil {
		return st
	}
	hosts := []string{"127.0.0.1"}
	if p.cfg.Mode == ModeLAN {
		hosts = append(hosts, LANAddrs()...)
	}
	for _, h := range hosts {
		st.Links = append(st.Links, Link(h, p.cfg.Port, p.cfg.Key))
	}
	return st
}

// Link is the address a browser opens, with the key.
func Link(host string, port int, key string) string {
	return "http://" + net.JoinHostPort(host, strconv.Itoa(port)) + "/?k=" + key
}

// LANAddrs are this computer's IPv4 addresses on its networks.
func LANAddrs() []string {
	var out []string
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok && !n.IP.IsLoopback() && !n.IP.IsLinkLocalUnicast() {
			if ip4 := n.IP.To4(); ip4 != nil {
				out = append(out, ip4.String())
			}
		}
	}
	return out
}
