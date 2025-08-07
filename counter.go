package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// counter page component
// Implementing increment/decrement functionality
type counter struct {
	app.Compo
	count int
}

func (c *counter) OnPreRender(ctx app.Context) {
	ctx.Page().SetTitle("Counter Page")
}

// OnMount initializes the counter value
func (c *counter) OnMount(ctx app.Context) {
	c.count = 0
	go func() {
		reloadOnUpdate(ctx)
	}()
}

// Render builds the UI for the counter page
func (c *counter) Render() app.UI {
	// determine color class based on count value
	colorClass := "text-black underline"
	if c.count < 0 {
		colorClass = "text-red-500"
	} else if c.count > 0 {
		colorClass = "text-green-500"
	}

	return app.Div().Class("m-4 p-4").Body(
		app.H1().
			Class("text-2xl font-bold mb-2").
			Text("Counter Page"),
		app.Div().
			Class("mb-2", colorClass).
			Text(fmt.Sprintf("Count: %d", c.count)),
		app.Button().
			Class("px-2 py-1 bg-red-400 min-w-10 text-white rounded mr-2").
			Text("-").
			OnClick(c.onDecrement),
		app.Button().
			Class("px-2 py-1 bg-green-400 min-w-10 text-white rounded").
			Text("+").
			OnClick(c.onIncrement),
	)
}

// onIncrement increases the counter and triggers an update
func (c *counter) onIncrement(ctx app.Context, e app.Event) {
	c.count++
	app.Log("Counter incremented to:", c.count)
	ctx.Update()
}

// onDecrement decreases the counter and triggers an update
func (c *counter) onDecrement(ctx app.Context, e app.Event) {
	c.count--
	app.Log("Counter decremented to:", c.count)
	ctx.Update()
}
