package main

import (
	"io"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// fetchUpdatedTime retrieves the latest timestamp from the server
func fetchUpdatedTime() string {
	r, err := http.Get("/updated")
	if err != nil {
		app.Log(err)
		return ""
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		app.Log(err)
		return ""
	}

	return string(b)
}

// clearServiceWorkers clears PWA service worker cache before reloading the page
func clearServiceWorkers() {
	window := app.Window()
	navigator := window.Get("navigator")
	if navigator.Get("serviceWorker").Truthy() {
		navigator.Get("serviceWorker").
			Call("getRegistrations").
			Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				regs := args[0]
				for i := 0; i < regs.Length(); i++ {
					reg := regs.Index(i)
					reg.Call("unregister")
					app.Log("Service worker unregistered")
				}

				return nil
			})).
			Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				app.Log("Error unregistering service worker:", args[0])

				return nil
			}))
	}
}

// reloadOnUpdate runs a loop that reloads the page when the server timestamp changes
func reloadOnUpdate(ctx app.Context) {
	var lastTime string
	for {
		newTime := fetchUpdatedTime()
		if lastTime == "" {
			lastTime = newTime
		} else if lastTime != newTime {
			clearServiceWorkers()
			time.Sleep(100 * time.Millisecond)
			ctx.Reload()
		}
		time.Sleep(500 * time.Millisecond)
	}
}
