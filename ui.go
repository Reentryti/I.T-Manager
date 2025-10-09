/*
	User Interface implementation
*/

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/* 	Main UI Fonction to integrate all parts of block*/
func userInterface(app *tview.Application) tview.Primitive {

	/*	Header Block
		text view for the application title display
	*/
	header := buildHeader("IP Table Manager")

	/*	Iptable rules Block
		Show all iptables rules available
		Function build layout and display rules
	*/
	blockRule, listRule := buildBlockRules()

	/*	Filter rules chain Block
		Display rules depending on chain values
	*/
	blockFilter := buildBlockFilter()

	/*	Forms Block
		Realize CRUD operations
		Function show dynamic forms depend on operations
	*/
	blockForm := buildBlockForm() //	form options for fill

	/*	Informations Block
		Display rules meaning
	*/
	blockInfos := buildBlockInfos()

	/*
		Key commands Block
		Show all commands for crud operations
		Build with this options (add, delete, update, refresh, save & quit)
	*/
	blockKey := buildBlockOptions(app, listRule, func(actionMenu string) {

		blockForm.Clear(true) //	Clear the form block

		/*	Form block content depend on the user action*/
		switch actionMenu {

		case "add":
			//	Form Title
			blockForm.SetTitle("Ajouter une regle")

			/*	Form values*/
			blockForm.AddDropDown("Chaine", []string{"Input", "Output", "Forward"}, 0, nil).
				AddDropDown("Methode d'ajout", []string{"-A(append)", "-I(insert)"}, 0, nil).
				AddInputField("IP Source", "", 20, nil, nil).
				AddInputField("IP Destination", "", 20, nil, nil).
				AddInputField("Protocol", "", 30, nil, nil).
				AddDropDown("Action", []string{"Accept", "Drop"}, 0, nil).
				AddButton("Ok", func() { //	addRules function call with error management and action confirmation

					/*	Retrieve  form values*/
					chainId, chain := blockForm.GetFormItemByLabel("Chaine").(*tview.DropDown).GetCurrentOption()
					if chainId < 0 {
						errorDisplay(app, "Chaine non selectionnée")
						return
					}

					optId, option := blockForm.GetFormItemByLabel("Methode d'ajout").(*tview.DropDown).GetCurrentOption()

					if optId < 0 {
						errorDisplay(app, "Methode d'ajout non selectionée")
						return
					}
					optionAdd := strings.Split(option, "")[0]

					srcInput := blockForm.GetFormItemByLabel("IP Source").(*tview.InputField).GetText()
					dstInput := blockForm.GetFormItemByLabel("IP Destination").(*tview.InputField).GetText()
					protoInput := blockForm.GetFormItemByLabel("Protocole").(*tview.InputField).GetText()

					_, action := blockForm.GetFormItemByLabel("Action").(*tview.DropDown).GetCurrentOption()

					if srcInput == "" || dstInput == "" {
						errorDisplay(app, "Veuillez renseigner tous les champs")
						return
					}

					/*	Define the current rule with the get values*/
					ruleAdd := RuleIptable{
						chain:       chain,
						source:      srcInput,
						destination: dstInput,
						protocol:    protoInput,
						action:      action,
					}
					/*	Result on this display message*/
					mesg := fmt.Sprintf("Ajouter la regle suivante? \n %s\nSource %s\nDestination %s\nProtocole", srcInput, dstInput, protoInput)

					/*	Confirmation on add rule option*/
					confirmAction(app, mesg, func() {
						err := addRules(ruleAdd, optionAdd)
						if err != nil {
							errorDisplay(app, fmt.Sprintf("Erreur : %v", err))
							return
						}

						refreshRules(listRule)
					})

				})

		case "delete":
			blockForm.SetTitle("Supprimer une regle")
			blockForm.AddInputField("Numero de regle", "", 5, nil, nil).
				AddButton("Ok", func() { //	delete function call

					numField := blockForm.GetFormItemByLabel("Numero de regle").(*tview.InputField) //get the rule number to delete
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

		case "update":
			blockForm.SetTitle("Modifier une regle")
			blockForm.AddInputField("Numero de regles", "", 5, nil, nil).
				AddInputField("Nouvelle source", "", 20, nil, nil).
				AddDropDown("Action", []string{"", ""}, 0, nil).
				AddButton("Ok", func() {
					//	modify function call
					/*	get text input*/
					lineTxt := blockForm.GetFormItemByLabel("Numero de regles").(*tview.InputField).GetText()
					lineNum, err := strconv.Atoi(lineTxt) //	convert to integer
					if err != nil || lineNum <= 0 {
						errorDisplay(app, "Chaine non selectionnée")
						return
					}

					/* Create the new rule*/
					newSource := blockForm.GetFormItemByLabel("Nouvelle source").(*tview.InputField).GetText()
					_, newAction := blockForm.GetFormItemByLabel("Action").(*tview.DropDown).GetCurrentOption()

					ruleChange := RuleIptable{
						source: newSource,
						action: newAction,
					}

					/*	Handle confirmation before changing choosen rules*/
					confirmAction(app, "Voulez vous vraiment modifier cette regle ?", func() {
						err := updateRules(ruleChange, lineNum)
						if err != nil {
							errorDisplay(app, fmt.Sprintf("Erreur %v", err))
							return
						}

						refreshRules(listRule) //	refresh rules after any rule change
					})
				})

		case "refresh":
			blockForm.SetTitle("Rafraichir la liste des regles")
			blockForm.AddButton("Refresh", func() {
				//	refresh function call
				refreshRules(listRule)
			})

		case "save":
			blockForm.SetTitle("Sauvegarder les regles")
			blockForm.AddButton("Save", func() {
				//	save function call
				confirmAction(app, "Vous etes entrain de sauvegarder toutes vos actions", func() {
					err := save()
					if err != nil {
						errorDisplay(app, fmt.Sprintf("Erreur %v", err))
						return
					}
				})
			})
		case "quit":
			app.Stop()
		}
	})

	/* Layout design*/
	/*	Left block (Filter + Keys)*/
	leftFrame := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(blockFilter, 0, 1, false).
		AddItem(blockKey, 0, 2, false)

	/*	Bottom block (Form + Informations)*/
	bottomFrame := tview.NewFlex().
		AddItem(blockForm, 0, 2, false).
		AddItem(blockInfos, 0, 1, false)

	/*	Main Block (iptables rules + leftFrame)*/
	mainFrame := tview.NewFlex().
		AddItem(leftFrame, 25, 1, true).
		AddItem(blockRule, 0, 3, false)

	/*	Full layout*/
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(mainFrame, 0, 6, true).
		AddItem(bottomFrame, 10, 2, false)

	/* Handling cursor focus to navigate between different frames*/
	curs := []tview.Primitive{blockFilter, blockKey, blockRule, blockForm, listRule}
	navigation(app, curs)

	return layout
}

/*	Simply building block to set border and title*/
func buildBlock(title string) *tview.Box {
	return tview.NewBox().
		SetBorder(true).
		SetTitle(title)
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

/*	Show values recap of current rule*/
func buildBlockInfos() *tview.TextView {

	infos := tview.NewTextView().
		SetText("[yellow]Informations: ")
	infos.SetBorder(true).SetTitle("Informations")

	return infos
}

/*	*/
func buildBlockFilter() *tview.Form {
	filter := tview.NewForm().
		AddInputField("Source", "", 15, nil, nil).
		AddInputField("Destination", "", 15, nil, nil).
		AddDropDown("Action", []string{"Accept", "Drop"}, 0, nil).
		AddButton("Filtre", func() {
			//filter callback
			fmt.Println("Filtrage")
		})
	//	Design
	filter.SetBorder(true).SetTitle("Filter")

	return filter
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
