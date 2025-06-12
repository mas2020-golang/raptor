/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/edit"
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
func newEditCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "edit",
		Aliases: []string{"secret", "sr"},
		Short:   "Edit a raptor object",
		Long:    `Edit a raptor object: secret`,
	}
	// Here you will define your flags and configuration settings.
	// Here you will define your flags and configuration settings.
	c.AddCommand(edit.EditSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
