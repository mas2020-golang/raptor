/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"
	"regexp"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	items           bool
	boxName, filter string
)

// boxCmd represents the box command
var ListSecretCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret", "sr"},
	Short:   "List secret",
	Long: `List all the secret in the --box given flag. Use the flag --name
to filter using a regular expression.`,
	Example: `$ raptor secret ls --box test
$ raptor secret ls --box test --name '^secret.*test$`,
	Run: func(cmd *cobra.Command, args []string) {
		listSecrets(cmd)
	},
}

func init() {
	ListSecretCmd.Flags().StringVarP(&filter, "filter", "f", "", "The secret name as a regexp (e.g. 'test.*')")
	ListSecretCmd.Flags().BoolVarP(&items, "items", "i", false, "Show the items' keys for the items saved into the secret")
	ListSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
}

func listSecrets(cmd *cobra.Command) {
	// load the local timezone
	// loc, err := time.LoadLocation("Local")
	// utils.Check(err, "")
	// output variables
	name, version, url, login := "", "", "", ""
	boxPath, _, box, err := utils.OpenBox(boxName, "")
	utils.Check(err, "")
	// get the max length for the NAME, LOGIN attribute
	maxName := getMaxNameLenght(box)
	maxLogin := getMaxLoginLenght(box)
	// compose the format for the header row
	formatS := fmt.Sprintf("%s%ds%s%d%s", "%-", maxName+2, "%-9s%-", maxLogin+2, "s%-47s%-11s%s\n")
	fmt.Printf(formatS, "NAME", "VERSION", "LOGIN", "URL", "ITEMS", "LAST-UPD")
	for _, s := range box.Secrets {
		loginFormatS := fmt.Sprintf("%s%ds", "%-", maxLogin+2)
		login = fmt.Sprintf(loginFormatS, s.Login)
		if len(s.Version) > 9 {
			version = s.Version[0:6] + "..."
		}
		url = s.Url
		if len(s.Url) > 44 {
			url = s.Url[0:42] + "..."
		}
		url = output.BlueS(fmt.Sprintf("%-47s", url))
		if len(s.Notes) > 30 {
			s.Notes = s.Notes[0:26] + "..."
		}
		nameFormatS := fmt.Sprintf("%s%ds", "%-", maxName+2)
		name = output.RedS(output.BoldS(fmt.Sprintf(nameFormatS, s.Name)))
		lastUpdated := s.LastUpdated
		// check the name flag
		if len(filter) > 0 {
			r, _ := regexp.Compile(filter)
			if r.MatchString(name) {
				fmt.Printf("%s%-9s%s%s%-11d%s\n", name, version, login, url, len(s.Others),
					lastUpdated)
				showItems(s)
			}
		} else {
			fmt.Printf("%s%-9s%s%s%-11d%s\n", name, version, login, url, len(s.Others),
				lastUpdated)
			showItems(s)
		}
	}
	v, _ := (*cmd).Parent().Flags().GetBool("verbose")
	if v {
		fmt.Println()
		output.InfoBox(fmt.Sprintf("secret read from the %s box\n", output.BlueS(boxPath)))
	}
}

func showItems(s *utils.Secret) {
	if !items {
		return
	}
	if s.Others != nil && len(s.Others) > 0 {
		//fmt.Println(" - items:")
		for k, _ := range s.Others {
			fmt.Printf("%-1s.%s\n", "", output.BoldS(k))
		}
	}
}

// getMaxNameLenght return the max lenght for the NAME attribute
func getMaxNameLenght(box *utils.Box) int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}
	return max
}

// getMaxLoginLenght return the max lenght for the LOGIN attribute
func getMaxLoginLenght(box *utils.Box) int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Login) > max {
			max = len(s.Login)
		}
	}
	return max
}
