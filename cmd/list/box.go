/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package list

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var ListBoxCmd = &cobra.Command{
	Use:     "boxes",
	Aliases: []string{"bo", "box"},
	Short:   "List the boxes",
	Long: `List the boxes present in the box folder. You can filter
using a regular expression.`,
	Example: `$ cryptex ls box --name 'test.*''`,
	Run: func(cmd *cobra.Command, args []string) {
		listBoxes(cmd)
	},
}

func init() {
	ListBoxCmd.Flags().StringVarP(&filter, "filter", "f", "", "The name of the box as a regexp (e.g. 'test.*')")
}

func listBoxes(cmd *cobra.Command) {
	var files []fs.DirEntry
	// get the folder box
	folderBox, err := utils.GetFolderBox()
	utils.Check(err, "")

	files, err = os.ReadDir(folderBox)
	utils.Check(err, "")

	fmt.Printf("%-25s%s\n", "NAME", "SIZE")
	for _, file := range files {
		fi, err := file.Info()
		utils.Check(err, fmt.Sprintf("problem getting the size of the file %s", file.Name()))
		// if --name is present check for the regular expression
		if len(filter) > 0 {
			r, _ := regexp.Compile(filter)
			if r.MatchString(file.Name()) {
				fmt.Printf("%-25s%d\n", file.Name(), fi.Size())
			}
		} else {
			fmt.Printf("%-25s%d\n", file.Name(), fi.Size())
		}
	}
	v, _ := cmd.Parent().Flags().GetBool("verbose")
	if v {
		fmt.Println()
		output.InfoBox(fmt.Sprintf("box folder set to %s\n", output.BlueS(folderBox)))
	}
}
