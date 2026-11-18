//go:build !nogui

package gui

import (
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// HasGUI is false in a `-tags nogui` build.
const HasGUI = true

// Run starts the desktop app: one window showing the page in-process, and
// the Browser access the user turned on in Settings. Quitting the app stops
// the browser page with it.
func Run(version string) error {
	// `wails3 dev` always points the app at a frontend dev server (Vite's).
	// ALGOX has none: the page is embedded and served by Handler, so a change
	// to it is picked up by dev mode's rebuild instead.
	os.Unsetenv("FRONTEND_DEVSERVER_URL")

	page := NewAppPage(NewHandler(version, true, nil))
	window := NewHandler(version, false, page)

	app := application.New(application.Options{
		Name:        "ALGOX",
		Description: "ALGOX",
		Assets:      application.AssetOptions{Handler: window},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		OnShutdown: page.Stop,
	})
	page.Start()

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "ALGOX",
		Width:     1100,
		Height:    720,
		MinWidth:  720,
		MinHeight: 480,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 46,
			Backdrop:                application.MacBackdropNormal,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(244, 244, 246),
		URL:              "/",
	})
	return app.Run()
}
