package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/robertkrimen/otto"
)

var (
	vm = otto.New()
	rarr = []rune("")
	rarrp = 0
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true

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
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = true
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

// All views
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// repl_textbox
func eval_repltextbox(g *gocui.Gui, v *gocui.View) error {
	repl_window, _ := g.View("repl_window")
	cmd := string(rarr)
	rarr = []rune{}
	v.SetCursor(0, 0)
	v.Clear()
	if cmd != "" {
		fmt.Fprintln(repl_window, "> "+cmd)
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		v.Clear()
		value, err := vm.Eval(cmd)
		fmt.Fprintln(repl_window, value)
		if err != nil {
			fmt.Fprintln(repl_window, err)
		}
	}
	return nil
}
func fillO(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprintf(v, "OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
	return nil
}
func clear_repltextbox(g *gocui.Gui, v *gocui.View) error {
	v.SetCursor(0, 0)
	v.Clear()
	rarr = []rune("")
	rarrp = -1
	return nil
}

// End Handler functions

// Helper functions
func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
// End Helper functions

func initKeybindings(g *gocui.Gui) error {
	// Quit Application
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// repl_textbox
	if err := g.SetKeybinding("repl_textbox", gocui.KeyEnter, gocui.ModNone, eval_repltextbox); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyCtrlO, gocui.ModNone, fillO); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyCtrlU, gocui.ModNone, clear_repltextbox); err != nil {
		log.Panicln(err)
	}

	// Handle all legal javascript characters
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789{};.=<>()[]'\"/\\-+!@#$%^&*~:" {
		g.SetKeybinding("repl_textbox", c, gocui.ModNone, mkEvtHandler(c))
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyBackspace, gocui.ModNone, backspace); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyBackspace2, gocui.ModNone, backspace); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyArrowLeft, gocui.ModNone, arrowleft); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeyArrowRight, gocui.ModNone, arrowright); err != nil {
		log.Panicln(err)
	}

	return nil
}

func mkEvtHandler(ch rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		fmt.Fprintf(v, fmt.Sprintf("%c", ch))
		v.MoveCursor(1, 0, true)
		rarr = append(rarr, ch)
		rarrp += 1
		return nil
	}
}

func backspace(g *gocui.Gui, v *gocui.View) error {
	v.EditDelete(true)
	if rarrp > 0 {
		rarr = append(rarr[:rarrp-1], rarr[rarrp:]...)
		rarrp -= 1
	}
	return nil
}
func arrowleft(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(-1, 0, true)
	if rarrp > 0 {
		rarrp -= 1
	}
	return nil
}
func arrowright(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(1, 0, true)
	if rarrp < len(rarr)-1 {
		rarrp += 1
	}
	return nil
}
