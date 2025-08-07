package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	PORT = ":8003"

	// TODO: change this to see hot reload on web browser
	VERSION = "0.1"
)

type hello struct {
	app.Compo
	lastUpdatedTime string
}

func (h *hello) Render() app.UI {
	return app.Div().Class("m-4 p-4 bg-gray-100 rounded-lg border border-gray-200 shadow-xl").Body(
		app.H1().Class("text-xl font-bold").Text("VERSION: "+VERSION),
		app.Hr().Class("my-2 text-gray-400"),
		app.P().Class("italic").Text("Hello World! 😎"),
	)
}

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

// clearServiceWorkers clears PWA service worker cache before reloading the page.
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

func (h *hello) OnMount(ctx app.Context) {
	h.lastUpdatedTime = ""
	for {
		newUpdatedTime := fetchUpdatedTime()

		if h.lastUpdatedTime == "" {
			h.lastUpdatedTime = newUpdatedTime
		} else if h.lastUpdatedTime != newUpdatedTime {
			clearServiceWorkers()

			time.Sleep(100 * time.Millisecond)
			ctx.Reload()
		}

		time.Sleep(500 * time.Millisecond)
	}

}

type lastUpdated struct {
	timestamp string
}

func (n *lastUpdated) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(n.timestamp))
}

// counter page component
// New counter component implementing increment/decrement functionality
type counter struct {
	app.Compo
	count int
}

func (c *counter) OnMount(ctx app.Context) {
	c.count = 0
}

func (c *counter) Render() app.UI {
	return app.Div().Class("m-4 p-4").Body(
		app.H1().Text("Counter Page").Class("text-2xl font-bold mb-2"),
		app.Div().Text(fmt.Sprintf("Count: %d", c.count)).Class("mb-2"),
		app.Button().Text("-").OnClick(c.onDecrement).Class("px-2 py-1 bg-red-500 text-white rounded mr-2"),
		app.Button().Text("+").OnClick(c.onIncrement).Class("px-2 py-1 bg-green-500 text-white rounded"),
	)
}

func (c *counter) onIncrement(ctx app.Context, e app.Event) {
	c.count++
	ctx.Update()
}

func (c *counter) onDecrement(ctx app.Context, e app.Event) {
	c.count--
	ctx.Update()
}

func main() {

	app.Route("/", func() app.Composer {
		return &hello{}
	})
	// Register counter page route
	app.Route("/counter", func() app.Composer {
		return &counter{}
	})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Scripts: []string{
			"https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4",
		},
	})

	http.Handle("/updated", &lastUpdated{
		timestamp: time.Now().Local().UTC().String(),
	})

	fmt.Println("Running on: http://localhost" + PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
