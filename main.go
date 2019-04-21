package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"gocap/panel"
)

func main() {
	g := panel.New(gocui.Output256)
	defer g.Close()

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
