/*
	TUI firewall manager for iptables ou nftables rules
*/

// Packages
package main

// Libraries
import (
	"fmt"
	//"os"
	"os/exec" //for exec command
	//"strings"
	//from github
	"github.com/rivo/tview" // go widget library
)

// Fonctions

/* running iptable command*/
func executeCommand() {}

/* handle privileges run*/
func sudoRun(args ...string) ([]byte, error) {
	cmdAdmin := exec.Command("sudo", append([]string{"iptables"}, args...)...)
	return cmdAdmin.CombinedOutput()
}

/* all commands to interact with the tui*/
func actionCommand() {}

func main() {

	//application title
	const title = "IP-Table Manager"
	//create tui application
	app := tview.NewApplication()

	//interface display
	ui := userInterface(app)

	//add header on main ui element
	if err := app.SetRoot(ui, true).Run(); err != nil {
		panic(err)
	}
}

func mainiptablefunction() {
	//	test current rules
	rules := currentRules()
	fmt.Println("Regles actuelles")
	for i, r := range rules {
		fmt.Printf("[%d] %s\n", i+1, r)
	}
}
