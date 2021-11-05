package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	router "marwan.io/vecty-router"
)

func main() {
	vecty.SetTitle("问答版搜索")
	b := &Body{}
	vecty.RenderBody(b)
	select {}
}

// Body renders the <body> tag
type Body struct {
	vecty.Core
}

// Render renders the <body> tag with the App as its children
func (b *Body) Render() vecty.ComponentOrHTML {
	s := search{}
	return elem.Body(
		router.NewRoute("/search/s", &s, router.NewRouteOpts{}),
	)

}
