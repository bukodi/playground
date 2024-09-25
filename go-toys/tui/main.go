package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	middlePane := tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top) (F3)")

	form := tview.NewForm().
		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddTextArea("Address", "", 40, 0, 0, nil).
		AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", nil).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Left (1/2 x width of Top)")
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF5 {
			app.SetFocus(middlePane)
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(form, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(middlePane, 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
