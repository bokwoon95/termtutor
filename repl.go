package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/robertkrimen/otto"
)

var (
	vm = otto.New()
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	maxX, maxY := g.Size()
	if v, err := g.SetView("repl_window", 0, 0, maxX-1, maxY-10); err != nil {
		v.Title = "repl_window"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		g.SetCurrentView("repl_window")
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyEnter, gocui.ModNone, send); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func send(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
