/*
	Temporary action register to store all action
	and apply them after save
*/

package main

import (
	"fmt"
	"os"
	"strings"
)

// temporary rules
var tempRules []RuleIptable

/* Load all rules on a temp file*/
func loadTempRules() {
	output, _ := sudoRun("iptables-save")
	tempRules = parseRules(string(output))
}

/*	Conversion on go slice*/
func parseRules(output string) []RuleIptable {
	var rules []RuleIptable

	lines := strings.Split(output, "\n")
	for _, line := range lines {

		if strings.HasPrefix(line, "-A") || strings.HasPrefix(line, "-I") {
			fields := strings.Fields(line)

			rule := RuleIptable{
				chain: fields[0],
			}

			for i := 1; i < len(fields); i++ {
				switch fields[i] {
				case "-s":
					rule.source = fields[i+1]
				case "-d":
					rule.destination = fields[i+1]
				case "-p":
					rule.protocol = fields[i+1]
				case "-j":
					rule.action = fields[i+1]
				default:

					if !strings.HasPrefix(fields[i], "-") {
						rule.chain = fields[i]
					}
				}
			}
			rules = append(rules, rule)
		}
	}

	return rules
}

/*	Save all temporary actions done*/
func save() error {
	tempTable := "" //temporay file path

	//	COnversion
	content := buildIptableSaveContent(tempRules)
	err := os.WriteFile(tempTable, []byte(content), 0644)

	if err != nil {
		return fmt.Errorf("Erreur lors de la creation du fichier temporaire : %v", err)
	}

	//
	_, err = sudoRun("iptables-restore", tempTable)

	if err != nil {
		return fmt.Errorf("Erreur lors de la restauration : %v", err)
	}
	return nil
}

/*	Build content compatible to iptable-save*/
func buildIptableSaveContent(rules []RuleIptable) string {

	var builder strings.Builder

	return builder.String()
}
