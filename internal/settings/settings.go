// Package settings reads and writes settings.json in algox's data folder.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"algox/internal/appdir"
)

// Settings are the user's choices the page shows in Settings › General.
type Settings struct {
	// Theme is "auto" (follow the system), "light" or "dark".
	Theme string `json:"theme"`
	// Lang is "auto" (follow the system), "zh" or "en".
	Lang string `json:"lang"`
}

// Default is what a fresh install starts with.
func Default() Settings { return Settings{Theme: "auto", Lang: "auto"} }

var themes = map[string]bool{"auto": true, "light": true, "dark": true}
var langs = map[string]bool{"auto": true, "zh": true, "en": true}

// Check refuses a value the page never offers.
func (s Settings) Check() error {
	if !themes[s.Theme] {
		return fmt.Errorf("unknown theme %q", s.Theme)
	}
	if !langs[s.Lang] {
		return fmt.Errorf("unknown language %q", s.Lang)
	}
	return nil
}

var mu sync.Mutex

const fileName = "settings.json"

// Load reads settings.json; a missing file is Default. Fields the file leaves
// out keep their default.
func Load() (Settings, error) {
	mu.Lock()
	defer mu.Unlock()
	s := Default()
	path, err := appdir.File(fileName)
	if err != nil {
		return s, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return Default(), fmt.Errorf("%s: %w", path, err)
	}
	if s.Check() != nil {
		return Default(), nil
	}
	return s, nil
}

// Save checks s and writes it.
func Save(s Settings) error {
	if err := s.Check(); err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	path, err := appdir.File(fileName)
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return appdir.WriteFile(path, b, 0o600)
}
