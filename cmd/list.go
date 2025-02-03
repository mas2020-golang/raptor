/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/list"
	"github.com/spf13/cobra"
)

// TODO: move the command in the cmd package
var (
	ListCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "Show the specified raptor objects",
		Long:    `Show the specified raptor objects: boxes, secrets, items`,
	}
	filter string
)

func newListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "Show the specified raptor objects",
		Long:    `Show the specified raptor objects: boxes, secrets, items`,
	}
	// Here you will define your flags and configuration settings.
	c.AddCommand(list.ListBoxCmd)
	c.AddCommand(list.ListSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
