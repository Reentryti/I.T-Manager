/*
	User Interface implementation
*/

package main

import (
	"github.com/rivo/tview"
)

/* 	Main UI Fonction to integrate all parts of block*/
func userInterface(app *tview.Application) tview.Primitive {

	header := buildHeader("IP Table Manager") //text view for title display

	rulesBox := buildBlockRules() // iptable rules display

	formBox := buildBlockForm()

	optionsBox := buildBlockOptions(func(actionMenu string) {

		formBox.Clear(true) //	Clear the form block

		//
		switch actionMenu {

		case "add":
			formBox.SetTitle("Ajouter une regle")
			formBox.AddInputField("IP Source", "", 20, nil, nil).
				AddInputField("IP Destination", "", 20, nil, nil).
				AddDropDown("Action", []string{"Accept", "Drop"}, 0, nil).
				AddButton("Ok", func() {
					//	addRules function call
				})

		case "delete":
			formBox.SetTitle("Supprimer une regle")
			formBox.AddInputField("Numero de regle", "", 5, nil, nil).
				AddButton("Ok", func() {
					//	delete function call
				})

		case "modify":
			formBox.SetTitle("Modifier une regle")
			formBox.AddInputField("Numero de regle", "", 5, nil, nil).
				AddInputField("Nouvelle source", "", 20, nil, nil).
				AddDropDown("Action", []string{"", ""}, 0, nil).
				AddButton("Ok", func() {
					//	modify function call
				})

		case "save":
			formBox.SetTitle("Sauvegarder les regles")
			formBox.AddButton("Save", func() {
				//	save function call
			})
		case "quit":
			app.Stop()
		}
	})

	/*othersBox := buildBlockFacult()
	Coming soon*/
	body := tview.NewFlex().
		AddItem(optionsBox, 20, 1, true).
		AddItem(rulesBox, 0, 3, false)

	render := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(body, 0, 5, true).
		AddItem(formBox, 12, 1, false)

	return render
}

/*	Header Tui block*/
func buildHeader(title string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter). //	text alignemnt
		SetText(title)
}

/*	Building the bloc that contains all current iptables rules*/
func buildBlockRules() tview.Primitive {
	rules := currentRules() //	slice of rules

	listRule := tview.NewList()

	frameRule := tview.NewFrame(listRule).
		SetBorder(true). //	border set
		SetTitle("Active Rules")

	//adding rules as item
	for _, rule := range rules {
		listRule.AddItem(rule, "", 0, nil)
	}

	return frameRule
}

/*	Building the block that contains all options commands */
func buildBlockOptions(onAction func(action string)) *tview.List {

	options := tview.NewList()

	options.AddItem("Add rule", "", 'a', func() {
		onAction("add")
	}) //	options for rules add
	options.AddItem("Delete rule", "", 'd', func() {
		onAction("delete")
	}) //	options for rules deletion
	options.AddItem("Change rule", "", 'c', func() {
		onAction("modify")
	}) //	options for rules change
	options.AddItem("Save rule", "", 's', func() {
		onAction("save")
	}) //	options to save current rules
	//	Quit application command
	options.AddItem("Quit", "", 'q', func() {
		onAction("quit")
	})

	return options
}

/*	Dont know what to make herer*/
func buildBlockFacult() {

}

/*	*/
func buildBlockForm() *tview.Form {
	form := tview.NewForm()
	//	SetBorder(true)

	form.AddTextView("info", "Selectionner une action à effectuer", 40, 2, false, false)
	form.SetBorder(true).
		SetTitle("Entrees").
		SetTitleAlign(tview.AlignLeft)

	return form
}
