package gui

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// ServeWeb serves the page to browsers with no desktop app (`algox web`).
// The key is made for this run unless key is given. host "" listens on
// every address.
func ServeWeb(version, host string, port int, key string) error {
	if key == "" {
		key = NewKey()
	}
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	shown := host
	if shown == "" || shown == "0.0.0.0" {
		shown = "127.0.0.1"
	}
	fmt.Println("ALGOX web:", Link(shown, port, key))
	if host == "" || host == "0.0.0.0" {
		for _, a := range LANAddrs() {
			fmt.Println("          ", Link(a, port, key))
		}
	}
	srv := &http.Server{
		Handler:           newWebGuard(key, port, NewHandler(version, true, nil)),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv.Serve(ln)
}
