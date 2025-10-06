/*
	User Interface implementation
*/

package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

/* 	Main UI Fonction to integrate all parts of block*/
func userInterface(app *tview.Application) tview.Primitive {

	header := buildHeader("IP Table Manager") //text view for title display

	//	rulesBox := buildBlockRules() // iptable rules display
	blockRule, listRule := buildBlockRules()

	formBox := buildBlockForm() //	form options for fill

	optionsBox := buildBlockOptions(app, listRule, func(actionMenu string) {

		formBox.Clear(true) //	Clear the form block

		//
		switch actionMenu {

		case "add":
			//	Form Title
			formBox.SetTitle("Ajouter une regle")
			
			//Form input
			formBox.AddInputField("Chaine", []string{"Input", "Output", "Forward"}, 0, nil).
				AddDropDown("Methode d'ajout", []string{"-A(append)", "-I(insert)"}, 0, nil).
				AddInputField("IP Source", "", 20, nil, nil).
				AddInputField("IP Destination", "", 20, nil, nil).
				AddInputField("Protocol", "", 30, nil, nil). 
				AddDropDown("Action", []string{"Accept", "Drop"}, 0, nil).
				AddButton("Ok", func() { //	addRules function call with error management and action confirmation

					/*	Retrieve  form values*/
					chainId, chain := formBox.GetFormItemByLabel("Chaine").(*tview.DropDown).GetCurrentOption()-
					if chainId < 0 {
						errorDisplay(app, "Chaine non selectionnée")
						return
					}

					optId, option := formBox.GetFormItemByLabel("Methode d'ajout").(*tview.DropDown).GetCurrentOption()

					if optId < 0{
						errorDisplay(app, "Methode d'ajout non selectionée")
						return
					}
					optionAdd := strings.Split(option, "")[0]

					srcInput := formBox.GetFormItemByLabel("IP Source").(*tview.InputField).GetText()
					dstInput := formBox.GetFormItemByLabel("IP Destination").(*tview.InputField).GetText()

					_, action := formBox.GetFormItemByLabel("Action").(*tview.DropDown).GetCurrentOption()

					if srcInput == "" || dstInput == "" {
						errorDisplay(app, "Veuillez renseigner tous les champs")
						return
					}

					/*	Confirmation on add rule option*/
					confirmAction(app, "Voulez vous vraiment ajouter cette regle?", func() {
						err := addRules(rule, optionAdd)
						if err != nil {
							errorDisplay(app, fmt.Sprintf("Erreur : %v", err))
							return
						}
						refreshRules(listRule)
					})

				})

		case "delete":
			formBox.SetTitle("Supprimer une regle")
			formBox.AddInputField("Numero de regle", "", 5, nil, nil).
				AddButton("Ok", func() { //	delete function call

					numField := formBox.GetFormItemByLabel("Numero de regle").(*tview.InputField) //get the rule number to delete
					lineStr := numField.GetText()

					line, err := strconv.Atoi(lineStr)
					/*	Error handling with the error function if no rule for this num*/
					if err != nil {
						errorDisplay(app, "Aucune regle ne correspond a ce numero")
						return
					}

					/*	Confirmation handling before deleting any kind of rules*/
					confirmAction(app, "Voulez vous vraiment supprimer cette regle ?", func() {
						err := delRules("INPUT", line)
						if err != nil {
							errorDisplay(app, fmt.Sprintf("Erreur %v", err))
							return
						}

						refreshRules(listRule) //	refresh rules after deletion
					})
				})

		case "modify":
			formBox.SetTitle("Modifier une regle")
			formBox.AddInputField("Numero de regle", "", 5, nil, nil).
				AddInputField("Nouvelle source", "", 20, nil, nil).
				AddDropDown("Action", []string{"", ""}, 0, nil).
				AddButton("Ok", func() {
					//	modify function call

					/*	Handle confirmation before changing choosen rules*/
					confirmAction(app, "Voulez vous vraiment modifier cette regle ?", func() {
						err := modRules()
						if err != nil {
							errorDisplay(app, fmt.Sprintf("Erreur %v", err))
							return
						}

						refreshRules(listRule) //	refresh rules after any rule change
					})
				})

		case "refresh":
			formBox.SetTitle("Rafraichir la liste des regles")
			formBox.AddButton("Refresh", func() {
				//	refresh function call
				refreshRules(listRule)
			})

		case "save":
			formBox.SetTitle("Sauvegarder les regles")
			formBox.AddButton("Save", func() {
				//	save function call
				confirmAction(app, "Vous etes entrain de sauvegarder toutes vos actions", func() {
					err := saveRule(listRule)
				})
			})
		case "quit":
			app.Stop()
		}
	})

	/*othersBox := buildBlockFacult()
	Coming soon*/
	body := tview.NewFlex().
		AddItem(optionsBox, 20, 1, true).
		AddItem(blockRule, 0, 3, false)

	render := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(body, 0, 5, true).
		AddItem(formBox, 12, 1, false)

	/* Handling cursor focus to navigate between different boxes/frames*/
	focusables := []tview.Primitive{listRule, optionsBox, formBox}
	current := 0

	/* 	Define Tab as navigation button*/
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			current = (current + 1) % len(focusables)
			app.SetFocus(focusables[current])
			return nil
		}
		return event
	})

	return render
}

/*	Header Tui block*/
func buildHeader(title string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter). //	text alignemnt
		SetText(title)
}

/*	Building the bloc that contains all current iptables rules*/
func buildBlockRules() (*tview.Flex, *tview.List) {
	// already call on refresh function rules := currentRules() //	slice of rules

	listRule := tview.NewList()
	refreshRules(listRule)

	frameRule := tview.NewFrame(listRule).
		AddText("Active Rules", true, tview.AlignCenter, tcell.ColorWhite).
		SetBorders(1, 1, 1, 0, 0, 0)

	blockRule := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frameRule, 0, 1, true)
	/*	already call on refresh function*/
	//adding rules as item
	//for _, rule := range rules {
	//	listRule.AddItem(rule, "", 0, nil)
	//}

	refreshRules(listRule)

	return blockRule, listRule
}

/*	Building the block that contains all options commands */
func buildBlockOptions(app *tview.Application, listRule *tview.List, onAction func(action string)) *tview.List {

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
	options.AddItem("Refresh list", "", 'r', func() {
		onAction("refresh")
	}) //	options to refresh rules list
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

/*	Refresh function to reload current iptables rules*/
func refreshRules(listRule *tview.List) {
	listRule.Clear()
	rules := currentRules()
	for _, rule := range rules {
		listRule.AddItem(rule, "", 0, nil)
	}
}

/*
Confirmation page before executing options

	parameter app for the main application
	parameter message for the text to display
	parameter onConfirm function to execut after user confirmation
*/
func confirmAction(app *tview.Application, message string, onConfirm func()) {
	confirmPage := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				onConfirm()
			}
			app.SetRoot(userInterface(app), true) //	if not close the modal
		})
	/*	make the confirmPage fullscreen */
	app.SetRoot(confirmPage, true)
}

/*
Error function display

	parameter app for the main application
	parameter message for the error to errorDisplay
*/
func errorDisplay(app *tview.Application, message string) {
	errorPage := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(userInterface(app), true) //	close the error page
		})
	app.SetRoot(errorPage, true)
}
