package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
)

// Launch launches the Fyne GUI application.
func Launch() {
	a := app.New()
	w := a.NewWindow("JSON to CSV/Excel")
	// UIコンポーネントの追加は後で記述
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}