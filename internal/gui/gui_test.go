package gui

import (
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func home(t *testing.T) string {
	dir := t.TempDir()
	t.Setenv("ALGOX_HOME", dir)
	return dir
}

func TestGuardLetsInOnlyWithTheKey(t *testing.T) {
	g := newWebGuard("secret", 1234, NewHandler("t", true, nil))

	// no key, a wrong key
	for _, path := range []string{"/", "/api/state", "/?k=nope"} {
		w := httptest.NewRecorder()
		g.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("%s: %d, want 401", path, w.Code)
		}
	}

	// the link sets the cookie and drops the key from the address
	w := httptest.NewRecorder()
	g.ServeHTTP(w, httptest.NewRequest("GET", "/?k=secret", nil))
	if w.Code != http.StatusSeeOther || strings.Contains(w.Header().Get("Location"), "secret") {
		t.Fatalf("link: %d %q", w.Code, w.Header().Get("Location"))
	}
	cookie := w.Result().Cookies()[0]
	if strings.Contains(cookie.Value, "secret") || !cookie.HttpOnly {
		t.Fatalf("cookie %+v", cookie)
	}

	get := func(c *http.Cookie, site string) int {
		r := httptest.NewRequest("GET", "/api/state", nil)
		r.AddCookie(c)
		if site != "" {
			r.Header.Set("Sec-Fetch-Site", site)
		}
		w := httptest.NewRecorder()
		g.ServeHTTP(w, r)
		return w.Code
	}
	if c := get(cookie, ""); c != http.StatusOK {
		t.Fatalf("with cookie: %d", c)
	}
	if c := get(cookie, "cross-site"); c != http.StatusForbidden {
		t.Fatalf("cross-site: %d, want 403", c)
	}

	// a new key signs every browser out
	g.SetKey("other")
	if c := get(cookie, ""); c != http.StatusUnauthorized {
		t.Fatalf("old cookie after new key: %d", c)
	}
}

func TestWebPageRoutesAreOnTheWindowOnly(t *testing.T) {
	home(t)
	page := NewAppPage(NewHandler("t", true, nil))
	for _, tc := range []struct {
		h    http.Handler
		want int
	}{{NewHandler("t", false, page), 200}, {NewHandler("t", true, page), 404}} {
		w := httptest.NewRecorder()
		tc.h.ServeHTTP(w, httptest.NewRequest("GET", "/api/web-page", nil))
		// a browser's handler has no such route: the file server finds no file
		if w.Code != tc.want {
			t.Fatalf("got %d, want %d", w.Code, tc.want)
		}
	}
}

func freePort(t *testing.T) int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestAppPageServesWithItsKeyAndTurnsOff(t *testing.T) {
	dir := home(t)
	p := NewAppPage(NewHandler("t", true, nil))
	p.Start()
	defer p.Stop()
	port := freePort(t)

	if err := p.Apply(ModeLocal, port, false); err != nil {
		t.Fatal(err)
	}
	st := p.Status()
	if !st.Serving || len(st.Links) != 1 {
		t.Fatalf("status %+v", st)
	}
	if _, err := os.Stat(filepath.Join(dir, pageFile)); err != nil {
		t.Fatal(err)
	}

	// the link opens the page; the bare address doesn't
	jar, _ := cookiejar.New(nil)
	browser := &http.Client{Jar: jar}
	res, err := browser.Get(st.Links[0])
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 || !strings.Contains(string(body), "ALGOX") {
		t.Fatalf("link: %d", res.StatusCode)
	}
	res, err = http.Get(strings.Split(st.Links[0], "?")[0])
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("no key: %d", res.StatusCode)
	}

	// the key is kept across a restart of the app
	again := NewAppPage(nil)
	if again.cfg.Key != p.cfg.Key || again.cfg.Mode != ModeLocal {
		t.Fatalf("not kept: %+v", again.cfg)
	}

	if err := p.Apply(ModeOff, port, false); err != nil {
		t.Fatal(err)
	}
	if p.Status().Serving {
		t.Fatal("still serving when off")
	}
	if _, err := http.Get(st.Links[0]); err == nil {
		t.Fatal("port still answers when off")
	}
}

func TestAppPageTurnsDownBadChoices(t *testing.T) {
	home(t)
	p := NewAppPage(nil)
	for _, c := range []struct {
		m Mode
		p int
	}{{"everyone", DefaultPort}, {ModeLocal, 80}, {ModeLAN, 70000}} {
		if err := p.Apply(c.m, c.p, false); err == nil {
			t.Fatalf("%v accepted", c)
		}
	}
}

func TestAppPageUnreadFileServesNothing(t *testing.T) {
	dir := home(t)
	if err := os.WriteFile(filepath.Join(dir, pageFile), []byte("{nope"), 0o600); err != nil {
		t.Fatal(err)
	}
	p := NewAppPage(nil)
	p.Start()
	defer p.Stop()
	st := p.Status()
	if st.Mode != ModeOff || st.Serving || st.Problem == "" {
		t.Fatalf("status %+v", st)
	}
}
