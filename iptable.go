/*
	Iptable functions
*/

package main

import (
	"fmt"
	//"os/exec"
	"strings"
)

//	Strutures
//
// IP Table rules
type RuleIptable struct {
	chain       string
	action      string
	source      string
	destination string
	protocol    string
}

/*	running iptable command*/
func execCommand() {}

/*	get active rules*/
func currentRules() []string {

	//cmdIptable := exec.Command("iptables", "-L") //command to request all iptables rules available

	//Error management
	outputIptable, err := sudoRun("-L")
	if err != nil {
		return []string{"Erreur lors de la récupération des regles"}
	}

	//convert output types to bytes array to string array
	convertOutput := strings.Split(string(outputIptable), "\n")

	return convertOutput
}

/*	adding new rules*/
func addRules(rule RuleIptable, optionAdd string) error {

	//	verify chosen option
	if optionAdd != "-A" && optionAdd != "-I" {
		return fmt.Errorf("Option invalide: %s (utilise -A ou -I)", optionAdd)
	}

	//	command options and parameters
	argsAdd := []string{optionAdd, rule.chain}

	if rule.source != "" {
		argsAdd = append(argsAdd, "-s", rule.source)
	}
	if rule.destination != "" {
		argsAdd = append(argsAdd, "-d", rule.destination)
	}
	if rule.protocol != "" {
		argsAdd = append(argsAdd, "-p", rule.protocol)
	}
	argsAdd = append(argsAdd, "-j", rule.action)

	//	Command to add rules
	//cmdAdd := exec.Command("iptables", argsAdd...)

	outputAdd, err := sudoRun(argsAdd...)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'ajout de regles : %v \n Sortie: %s", err, string(outputAdd))
	}

	return nil
}

/*	deleting rules*/
func delRules(chain string, line int) error {

	//	Line number is strictly positive
	if line <= 0 {
		return fmt.Errorf("Numero de ligne invalide %d", line)
	}

	argsDel := []string{"-D", chain, fmt.Sprintf("%d", line)}

	//cmdDel := exec.Command("iptables", argsDel...)

	outputDel, err := sudoRun(argsDel...)
	if err != nil {
		return fmt.Errorf("Erreur lors de la suppression de regles : %v \n Sortie: %s", err, string(outputDel))
	}

	return nil
}

/*	modifying rules*/
func modRules(rule RuleIptable, line int) error {
	//	line must be a positive number
	if line <= 0 {
		return fmt.Errorf("Numero de ligne invalide :%d", line)
	}

	argsMod := []string{"-R", rule.chain, fmt.Sprintf("%d", line)}

	if rule.source != "" {
		argsMod = append(argsMod, "-s", rule.source)
	}
	if rule.destination != "" {
		argsMod = append(argsMod, "-d", rule.destination)
	}
	if rule.protocol != "" {
		argsMod = append(argsMod, "-p", rule.protocol)
	}

	argsMod = append(argsMod, "-j", rule.action)

	//cmdMod := exec.Command("iptables", argsMod...)
	outputMod, err := sudoRun(argsMod...)
	if err != nil {
		return fmt.Errorf("Erreur lors de la modification de regles: %v \n Sortie: %s", err, string(outputMod))
	}

	return nil
}

/*
refresh rules list
Not necessary right now
*/
func refreshRules() {}
