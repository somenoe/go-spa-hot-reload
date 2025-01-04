//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Runs go mod download and then installs the binary.
func Build() error {
	if err := sh.RunWith(map[string]string{
		"GOARCH": "wasm",
		"GOOS":   "js",
	}, "go", "build", "-o", "web/app.wasm"); err != nil {
		return err
	}

	return sh.Run("go", "build", "-o", "tmp/main.exe")
}
