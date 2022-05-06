/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"fmt"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"regexp"
	"time"
)

var (
	secretName string
)

// boxCmd represents the box command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List secret",
	Long:    `List all the secret in the --box given flag. Use the flag --name
to filter using a regular expression.`,
	Example: `$ cryptex secret ls --box test
$ cryptex secret ls --box test --name '^secret.*test$`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	listCmd.Flags().StringVarP(&secretName, "name", "n", "", "The secret name as a regexp (e.g. 'test.*')")
}

func list() {
	// load the local timezone
	loc, err := time.LoadLocation("Local")
	utils.Check(err, "")
	// check the folder .cryptex
	home, err := os.UserHomeDir()
	utils.Check(err, "")
	// read the file in the home dir
	boxF := path.Join(home, ".cryptex", "boxes", boxName)

	// open the box
	err = openBox(boxF)
	utils.Check(err, "")
	// ls the secret
	fmt.Printf("%-15s%-10s%-31s%-31s%-11s%s\n", "NAME", "VERSION", "URL", "NOTES", "ITEMS", "LAST-UPD")
	for _, s := range box.Secrets{
		if len(s.Version) > 9{
			s.Version =  s.Version[0:6] + "..."
		}
		if len(s.Name) > 14{
			s.Name =  s.Name[0:10] + "..."
		}
		if len(s.Url) > 30{
			s.Url =  s.Url[0:26] + "..."
		}
		if len(s.Notes) > 30{
			s.Notes =  s.Notes[0:26] + "..."
		}
		s.Name = utils.LightRedS(utils.BoldS(fmt.Sprintf("%-15s", s.Name)))
		lastUpdated := s.LastUpdated.AsTime().In(loc).Format("Jan 2 15:04 2006 MST")
		// check the name flag
		if len(secretName) > 0{
			r, _ := regexp.Compile(secretName)
			if r.MatchString(s.Name){
				fmt.Printf("%s%-10s%-31s%-31s%-11d%s\n", s.Name, s.Version, s.Url, s.Notes, len(s.Others),
					lastUpdated)
			}
		}else{
			fmt.Printf("%s%-10s%-31s%-31s%-11d%s\n", s.Name, s.Version, s.Url, s.Notes, len(s.Others),
				lastUpdated)
		}
	}
}
