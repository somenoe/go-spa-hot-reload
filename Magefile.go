//go:build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	DIFF_FORMATTER = "05.000"
)

func Wasm() {
	startTime := time.Now()
	sh.RunWith(map[string]string{
		"GOARCH": "wasm",
		"GOOS":   "js",
	}, "go", "build", "-o", "web/app.wasm")
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	formattedDiff := time.Time{}.Add(diff).Format(DIFF_FORMATTER)

	fmt.Println("wasm\t=>", formattedDiff)
}

func Build() {
	mg.Deps(Wasm)
	startTime := time.Now()
	sh.Run("go", "build", "-o", "tmp/main.exe")
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	formattedDiff := time.Time{}.Add(diff).Format(DIFF_FORMATTER)

	fmt.Println("build\t=>", formattedDiff)
}

func Run() {
	mg.Deps(Build)

	sh.Run("tmp/main.exe")
}

func Dev() {
	mg.Deps(Wasm)

	sh.Run("go", "run", ".")
}
