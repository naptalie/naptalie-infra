package main
import (
        "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


func NewSimpleMeowTUI() *tview.Application {
    app := tview.NewApplication()
    text := tview.NewTextView().
        SetText("? Meow! Press Ctrl+C to quit").
        SetTextColor(tcell.ColorPink)
    
    app.SetRoot(text, true)
    return app
}
