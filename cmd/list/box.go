/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
)

// // String returns a formatted string representation of the box
// func (b utils.Box) String() string {
// 	return
// }

//

// NewListBoxCmd creates and returns a new list boxes command
func NewListBoxCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:     "boxes",
		Aliases: []string{"bo", "box"},
		Short:   "List the boxes",
		Long: `List the boxes present in the box folder. You can filter
using a regular expression.`,
		Example: `$ cryptex ls box --filter 'test.*'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListBoxes(filter)
		},
	}

	cmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter boxes by name using regexp (e.g. 'test.*')")

	return cmd
}

func runListBoxes(filter string) error {
	boxes, err := ListBoxes(filter)
	if err != nil {
		return fmt.Errorf("error retrieving boxes: %w", err)
	}

	printBoxes(boxes)
	return nil
}

func printBoxes(boxes []utils.Box) {
	fmt.Printf("%-25s%s\n", "NAME", "SIZE")
	for _, b := range boxes {
		fmt.Printf("%-25s%s\n", b.Name, strconv.FormatInt(b.Size, 10))
	}
}

func ListBoxes(filter string) ([]utils.Box, error) {
	folderBox := utils.GetFolderBox()
	files, err := os.ReadDir(folderBox)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", folderBox, err)
	}

	var filterRegex *regexp.Regexp
	if filter != "" {
		filterRegex, err = regexp.Compile(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern '%s': %w", filter, err)
		}
	}

	var boxes []utils.Box
	for _, file := range files {
		if shouldIncludeFile(file, filterRegex) {
			box, err := createBoxFromFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to process file %s: %w", file.Name(), err)
			}
			boxes = append(boxes, box)
		}
	}

	return boxes, nil
}

func shouldIncludeFile(file fs.DirEntry, filterRegex *regexp.Regexp) bool {
	if filterRegex == nil {
		return true
	}
	return filterRegex.MatchString(file.Name())
}

func createBoxFromFile(file fs.DirEntry) (utils.Box, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return utils.Box{}, err
	}

	return utils.Box{
		Name: file.Name(),
		Size: fileInfo.Size(),
	}, nil
}
