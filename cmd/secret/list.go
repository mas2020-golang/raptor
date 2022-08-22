/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
)

var (
	filter string
	items  bool
)

// boxCmd represents the box command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List secret",
	Long: `List all the secret in the --box given flag. Use the flag --name
to filter using a regular expression.`,
	Example: `$ cryptex secret ls --box test
$ cryptex secret ls --box test --name '^secret.*test$`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "The secret name as a regexp (e.g. 'test.*')")
	listCmd.Flags().BoolVarP(&items, "items", "i", false, "Show the items' keys for the items saved into the secret")
}

func list() {
	// load the local timezone
	loc, err := time.LoadLocation("Local")
	utils.Check(err, "")
	// open the box
	_, err = openBox(boxName)
	utils.Check(err, "")
	// get the max length for the NAME, LOGIN attribute
	maxName := getMaxNameLenght()
	maxLogin := getMaxLoginLenght()
	// compose the format for the header row
	formatS := fmt.Sprintf("%s%ds%s%d%s", "%-", maxName+2, "%-9s%-", maxLogin+2, "s%-47s%-11s%s\n")
	fmt.Printf(formatS, "NAME", "VERSION", "LOGIN", "URL", "ITEMS", "LAST-UPD")
	for _, s := range box.Secrets {
		loginFormatS := fmt.Sprintf("%s%ds", "%-", maxLogin+2)
		s.Login = fmt.Sprintf(loginFormatS, s.Login)
		if len(s.Version) > 9 {
			s.Version = s.Version[0:6] + "..."
		}
		if len(s.Url) > 44 {
			s.Url = s.Url[0:42] + "..."
		}
		s.Url = utils.BlueS(fmt.Sprintf("%-47s", s.Url))
		if len(s.Notes) > 30 {
			s.Notes = s.Notes[0:26] + "..."
		}
		nameFormatS := fmt.Sprintf("%s%ds", "%-", maxName+2)
		s.Name = utils.RedS(utils.BoldS(fmt.Sprintf(nameFormatS, s.Name)))
		lastUpdated := s.LastUpdated.AsTime().In(loc).Format("Jan 2 15:04 2006 MST")
		// check the name flag
		if len(filter) > 0 {
			r, _ := regexp.Compile(filter)
			if r.MatchString(s.Name) {
				fmt.Printf("%s%-9s%s%s%-11d%s\n", s.Name, s.Version, s.Login, s.Url, len(s.Others),
					lastUpdated)
				showItems(s)
			}
		} else {
			fmt.Printf("%s%-9s%s%s%-11d%s\n", s.Name, s.Version, s.Login, s.Url, len(s.Others),
				lastUpdated)
			showItems(s)
		}
	}
}

func showItems(s *protos.Secret) {
	if !items {
		return
	}
	if s.Others != nil && len(s.Others) > 0 {
		//fmt.Println(" - items:")
		for k, _ := range s.Others {
			fmt.Printf("%-1s.%s\n", "", utils.BoldS(k))
		}
	}
}

// getMaxNameLenght return the max lenght for the NAME attribute
func getMaxNameLenght() int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}
	return max
}

// getMaxLoginLenght return the max lenght for the LOGIN attribute
func getMaxLoginLenght() int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Login) > max {
			max = len(s.Login)
		}
	}
	return max
}
