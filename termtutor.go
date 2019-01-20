package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/robertkrimen/otto"
)

type QnA struct {
	Question string
	Answer   string
}

var (
	vm    = otto.New()
	rarr  = []rune("")
	rarrp = 0
	mode  = 1
	replm = 1
	qnas  []QnA
	qn = 0
)

func main() {
	// Read in json data
	file, err := ioutil.ReadFile("./data.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &qnas)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
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
	if v, err := g.SetView("javascript_repl", 0, 0, maxX-1, maxY-10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "javascript_repl"
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = true
		if _, err = setCurrentViewOnTop(g, "javascript_repl"); err != nil {
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
	if v, err := g.SetView("questions", 0, 0, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "questions"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		qn = 0
		fmt.Fprintln(v, qnas[qn].Question + "\n\n\n")
		qn += 1
		if _, err = setCurrentViewOnTop(g, "questions"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("answers", maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "answers"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		if _, err = setCurrentViewOnTop(g, "answers"); err != nil {
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
	javascript_repl, _ := g.View("javascript_repl")
	cmd := string(rarr)
	rarr = []rune{}
	v.SetCursor(0, 0)
	v.Clear()
	if cmd != "" {
		fmt.Fprintln(javascript_repl, "> "+cmd)
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		v.Clear()
		value, err := vm.Eval(cmd)
		fmt.Fprintln(javascript_repl, value)
		if err != nil {
			fmt.Fprintln(javascript_repl, err)
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
	// if err := g.SetKeybinding("repl_textbox", gocui.KeyTab, gocui.ModNone, togglereplm); err != nil {
	// 	log.Panicln(err)
	// }
	// if err := g.SetKeybinding("javascript_repl", gocui.KeyTab, gocui.ModNone, togglereplm); err != nil {
	// 	log.Panicln(err)
	// }

	// answers
	if err := g.SetKeybinding("answers", gocui.KeyEnter, gocui.ModNone, nextquestion); err != nil {
		log.Panicln(err)
	}

	// Swap modes
	if err := g.SetKeybinding("repl_textbox", gocui.KeyCtrlQ, gocui.ModNone, switchlesson); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("answers", gocui.KeyCtrlW, gocui.ModNone, switchrepl); err != nil {
		log.Panicln(err)
	}

	// Handle all legal javascript characters
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789{};.=<>()[]'\"/\\-+!@#$%^&*~:," {
		g.SetKeybinding("repl_textbox", c, gocui.ModNone, mkEvtHandler(c))
	}
	if err := g.SetKeybinding("repl_textbox", gocui.KeySpace, gocui.ModNone, mkEvtHandler(' ')); err != nil {
		log.Panicln(err)
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
func togglereplm(g *gocui.Gui, v *gocui.View) error {
	if replm == 1 {
		g.SetCurrentView("javascript_repl")
		replm = 0
	} else {
		g.SetCurrentView("repl_textbox")
		replm = 1
	}
	return nil
}
func switchlesson(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("questions")
	g.SetViewOnTop("answers")
	g.SetCurrentView("answers")
	return nil
}
func switchrepl(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("javascript_repl")
	g.SetViewOnTop("repl_textbox")
	g.SetCurrentView("repl_textbox")
	return nil
}
func nextquestion(g *gocui.Gui, v *gocui.View) error {
	w, _ := g.View("questions")
	wlen := len(w.BufferLines())
	if qn > 0 {
		fmt.Fprintln(v, "\n" + qnas[qn-1].Answer)
	} else {
		fmt.Fprintln(v, "\n" + qnas[9].Answer)
	}
	fmt.Fprintln(w, qnas[qn].Question + "\n\n\n")
	qn += 1
	qn = qn % 10
	newlines(v, w, wlen)
	return nil
}
func newlines(v *gocui.View, w *gocui.View, wlen int) error {
	// wlen := len(w.BufferLines())
	vlen := len(v.BufferLines())
	n := wlen - vlen
	for i := 0; i < n - 1; i++ {
		fmt.Fprintf(v, "\n")
	}
	vlen = len(v.BufferLines())
	v.SetCursor(0, vlen)
	return nil
}
