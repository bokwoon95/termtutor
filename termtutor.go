package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("repl_window", 0, 0, maxX-1, maxY-10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "repl_window"
		v.Editable = true
		v.Wrap = true
		if _, err = setCurrentViewOnTop(g, "repl_window"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("repl_textbox", 0, maxY-9, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "repl_textbox"
		v.Editable = true
		v.Wrap = true
		if _, err = setCurrentViewOnTop(g, "repl_textbox"); err != nil {
			return err
		}
	}
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	// Quit Application
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	return nil
}

// Handler functions
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
// Handler functions

// Helper functions
func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
// Helper functions
