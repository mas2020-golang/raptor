/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package box

import (
	"fmt"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"regexp"
)

var (
	boxName string
)

// boxCmd represents the box command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the boxes",
	Long: `List the boxes present in the box folder. You can filter
using a regular expression.`,
	Example: `$ cryptex box ls --name 'test.*''`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	listCmd.Flags().StringVarP(&boxName, "name", "n", "", "The name of the box as a regexp (e.g. 'test.*')")
}

func list() {
	// check the folder .cryptex
	home, err := os.UserHomeDir()
	utils.Check(err, "")
	// read the file in the home dir
	files, err := os.ReadDir(path.Join(home, ".cryptex", "boxes"))
	utils.Check(err, "")
	fmt.Printf("%-16s%-10s\n", "NAME", "SIZE")
	for _, file := range files {
		fi, err := file.Info()
		utils.Check(err, fmt.Sprintf("problem getting the size of the file %s", file.Name()))
		// if --name is present check for the regular expression
		if len(boxName) > 0{
			r, _ := regexp.Compile(boxName)
			if r.MatchString(file.Name()){
				fmt.Printf("%-16s%-10d\n", file.Name(), fi.Size())
			}
		}else{
			fmt.Printf("%-16s%-10d\n", file.Name(), fi.Size())
		}



	}
}
