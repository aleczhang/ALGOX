package gui

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strconv"
	"sync/atomic"
)

// NewKey makes a random key for a browser's link.
func NewKey() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// webGuard lets a browser in only with the key: once from the link (?k=),
// which sets a cookie and redirects to the same address without the key,
// then by the cookie. Writes from another origin are refused
// (net/http CrossOriginProtection), as is any request whose Sec-Fetch-Site
// names another site, since the Lax cookie rides along with those.
type webGuard struct {
	key    atomic.Pointer[string]
	cookie string
	next   http.Handler
	cop    *http.CrossOriginProtection
}

func newWebGuard(key string, port int, next http.Handler) *webGuard {
	g := &webGuard{cookie: "algox_web_" + strconv.Itoa(port), next: next, cop: http.NewCrossOriginProtection()}
	g.key.Store(&key)
	return g
}

// SetKey swaps the key; every browser signed in with the old one is let go.
func (g *webGuard) SetKey(key string) { g.key.Store(&key) }

func (g *webGuard) token() string {
	sum := sha256.Sum256([]byte("algox-web:" + *g.key.Load()))
	return hex.EncodeToString(sum[:])
}

func same(a, b string) bool { return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1 }

func (g *webGuard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if k := r.URL.Query().Get("k"); k != "" {
		if !same(k, *g.key.Load()) {
			denied(w)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: g.cookie, Value: g.token(), Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
		u := *r.URL
		q := u.Query()
		q.Del("k")
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.RequestURI(), http.StatusSeeOther)
		return
	}
	c, err := r.Cookie(g.cookie)
	if err != nil || !same(c.Value, g.token()) {
		denied(w)
		return
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site == "cross-site" || site == "same-site" {
		http.Error(w, "cross-origin request refused", http.StatusForbidden)
		return
	}
	g.cop.Handler(g.next).ServeHTTP(w, r)
}

func denied(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>ALGOX</title>` +
		`<body style="font-family:system-ui;padding:40px;color:#444">` +
		`<p>需要访问链接。请在 ALGOX 桌面应用的「设置 › 浏览器访问」中复制链接。</p>` +
		`<p>Open the link from ALGOX's Settings › Browser access.</p></body>`))
}
