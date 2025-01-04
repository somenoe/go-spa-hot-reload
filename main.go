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
	PORT               = ":8003"
	DELAY_CHECK_UPDATE = 2

	// TODO: change this to see hot reload on web browser
	VERSION = "0.1"
)

type hello struct {
	app.Compo
	lastUpdatedTime string
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("VERSION: "+VERSION),
		app.P().Text("Hello World! 😎"),
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

func (h *hello) OnMount(ctx app.Context) {
	h.lastUpdatedTime = ""
	for {
		newUpdatedTime := fetchUpdatedTime()

		if h.lastUpdatedTime == "" {
			h.lastUpdatedTime = newUpdatedTime
		} else if h.lastUpdatedTime != newUpdatedTime {
			ctx.Reload()
		}

		time.Sleep(DELAY_CHECK_UPDATE * time.Second)
	}

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
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	http.Handle("/updated", &lastUpdated{
		timestamp: time.Now().Local().UTC().String(),
	})

	fmt.Println("Running on: http://localhost" + PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
