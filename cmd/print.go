/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/print"
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
// boxCmd represents the box command
func newPrintCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "print",
		Aliases: []string{"pr"},
		Short:   "Print a raptor object",
		Long:    `You can easily print the details of box or a secret`,
	}

	// Here you will define your flags and configuration settings.
	c.AddCommand(print.PrintSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
