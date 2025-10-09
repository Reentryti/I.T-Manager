package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*	Navigation function to handle cursor placement on ui*/
func navigation(app *tview.Application, curs []tview.Primitive) {

	position := 0
	cssFocus(curs, position)
	app.SetFocus(curs[position])

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {

		/*	Cursor handle with Tab touch*/
		case tcell.KeyTab:
			position = (position + 1) % len(curs)
			app.SetFocus(curs[position])
			cssFocus(curs, position)
			return nil

		/*	Cursor handler with tab touch*/
		case tcell.KeyBacktab:
			position = (position - 1 + len(curs)) % len(curs)
			app.SetFocus(curs[position])
			cssFocus(curs, position)
			return nil

		/*	Keybinds to close the app with ctrl + C or X*/
		case tcell.KeyCtrlC, tcell.KeyCtrlX:
			app.Stop()
			return nil
		default:
			return event
		}
	})
}

/*	Styling navigation between blocks*/
func cssFocus(curs []tview.Primitive, position int) {

	for i, p := range curs {

		switch box := p.(type) {

		case *tview.Frame:
			if i == position {
				box.SetBorderColor(tcell.ColorOrangeRed)
			} else {
				box.SetBorderColor(tcell.ColorGray)
			}

		case *tview.Box:

		case *tview.List:

		case *tview.Form:
		}
	}
}
