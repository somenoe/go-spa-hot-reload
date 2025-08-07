package main

import (
	"fmt"
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
}

func (h *hello) Render() app.UI {
	return app.Div().Class("m-4 p-4 bg-gray-100 rounded-lg border border-gray-200 shadow-xl").Body(
		app.H1().Class("text-xl font-bold").Text("VERSION: "+VERSION),
		app.Hr().Class("my-2 text-gray-400"),
		app.P().Class("italic").Text("Hello World! 😎"),
	)
}

// OnMount shares reload-on-update logic via reload.go
func (h *hello) OnMount(ctx app.Context) {
	reloadOnUpdate(ctx)
}

type lastUpdated struct {
	timestamp string
}

func (n *lastUpdated) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(n.timestamp))
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
