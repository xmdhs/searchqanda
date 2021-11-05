package main

import (
	"github.com/hexops/vecty"
)

func main() {
	vecty.SetTitle("问答版搜索")
	s := &search{}
	vecty.RenderBody(s)
	select {}
}
