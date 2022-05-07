/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"fmt"
	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"regexp"
	"time"
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
	// ls the secret
	fmt.Printf("%-15s%-10s%-31s%-31s%-11s%s\n", "NAME", "VERSION", "URL", "NOTES", "ITEMS", "LAST-UPD")
	for _, s := range box.Secrets {
		if len(s.Version) > 9 {
			s.Version = s.Version[0:6] + "..."
		}
		if len(s.Name) > 14 {
			s.Name = s.Name[0:10] + "..."
		}
		if len(s.Url) > 30 {
			s.Url = s.Url[0:26] + "..."
		}
		s.Url = utils.BlueS(fmt.Sprintf("%-31s", s.Url))
		if len(s.Notes) > 30 {
			s.Notes = s.Notes[0:26] + "..."
		}
		s.Name = utils.RedS(utils.BoldS(fmt.Sprintf("%-15s", s.Name)))
		lastUpdated := s.LastUpdated.AsTime().In(loc).Format("Jan 2 15:04 2006 MST")
		// check the name flag
		if len(filter) > 0 {
			r, _ := regexp.Compile(filter)
			if r.MatchString(s.Name) {
				fmt.Printf("%s%-10s%s%-31s%-11d%s\n", s.Name, s.Version, s.Url, s.Notes, len(s.Others),
					lastUpdated)
				showItems(s)
			}
		} else {
			fmt.Printf("%s%-10s%-31s%-31s%-11d%s\n", s.Name, s.Version, s.Url, s.Notes, len(s.Others),
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
			fmt.Printf("%-1s.%s\n","", utils.BoldS(k))
		}
	}
}
