/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package create

import (
	"fmt"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var secretName string

func NewAddItemCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "item <NAME>",
		Args:    cobra.MinimumNArgs(1),
		Short:   "Add an item to a secret",
		Long:    `You can add an item to the items collection of the given secret`,
		Example: "cryptex add-item 'my-new-item' --secret first-secret --box test-box",
		Run: func(cmd *cobra.Command, args []string) {
			err := addItem(args[0], cmd)
			if err != nil {
				output.Error("", fmt.Sprintf("an error occurred during the item add: %s", err))
			}
		},
	}

	// Here you will define your flags and configuration settings.
	c.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box")
	c.Flags().StringVarP(&secretName, "secret", "s", "", "The name of the secret where to add the item")
	c.MarkFlagRequired("secret")
	return c
}

func addItem(name string, cmd *cobra.Command) error {
	found := false
	// open the box
	boxPath, key, box, err := utils.OpenBox(boxName)
	if err != nil {
		return err
	}

	// search for the secret
	for _, s := range box.Secrets {
		if s.Name == secretName {
			found = true
			if _, ok := s.Others[name]; ok {
				return fmt.Errorf("the item already exists for the --secret %s", secretName)
			}

			// add the new item
			fmt.Println(output.GreenS(strings.Repeat("-", 35)))
			if s.Others == nil {
				s.Others = make(map[string]string)
			}
			pwd, err := utils.AskForPassword("Item password: ", true, 1)
			if err != nil {
				return err
			}
			s.Others[name] = pwd
			s.LastUpdated = timestamppb.Now()
			// save the box
			err = utils.SaveBox(boxPath, key, box)
			if err != nil {
				return err
			}
			fmt.Println()
			utils.Success(output.BoldS("new item saved!"))
			v, _ := (*cmd).Parent().Flags().GetBool("verbose")
			if v {
				fmt.Println()
				output.InfoBox(fmt.Sprintf("item added to %q into the box %s\n", secretName, output.BlueS(boxPath)))
			}
		}
	}
	if !found {
		return fmt.Errorf("no secret %q found for the box %q", secretName, boxPath)
	}

	return nil
}
