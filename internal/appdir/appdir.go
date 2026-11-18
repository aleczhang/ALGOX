// Package appdir names the folder algox keeps its files in.
package appdir

import (
	"os"
	"path/filepath"
)

// Getenv returns an environment variable that names a folder, or "" when it
// is unset or relative: a relative path would depend on where algox was
// started from.
func Getenv(name string) string {
	v := os.Getenv(name)
	if v == "" || !filepath.IsAbs(v) {
		return ""
	}
	return v
}

// Dir is algox's data folder: $ALGOX_HOME, else the user's config folder's
// "algox". It is made (0700) when missing.
func Dir() (string, error) {
	dir := Getenv("ALGOX_HOME")
	if dir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(base, "algox")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

// File is the path of name in Dir.
func File(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

// WriteFile writes data through a temp file and a rename, so a crash never
// leaves half a file.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
