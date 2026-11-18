//go:build nogui

package gui

import "errors"

// HasGUI is false in a `-tags nogui` build.
const HasGUI = false

// Run says this build has no desktop app.
func Run(string) error {
	return errors.New("this build has no desktop app (-tags nogui); use `algox web`")
}
