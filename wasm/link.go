package main

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
)

//https://github.com/marwan-at-work/vecty-router

func redirect(route string) {
	js.Global().Get("history").Call(
		"pushState",
		map[string]interface{}{"redirectRoute": route},
		route,
		route,
	)
}

func link(c vecty.Component, route, text string) *vecty.HTML {
	return elem.Anchor(
		vecty.Markup(
			prop.Href(route),
			event.Click(onClick(c, route)).PreventDefault(),
		),
		vecty.Text(text),
	)
}

func onClick(c vecty.Component, route string) func(e *vecty.Event) {
	return func(e *vecty.Event) {
		redirect(route)
		vecty.Rerender(c)
	}
}
