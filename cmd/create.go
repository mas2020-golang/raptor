/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/create"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Add an object to raptor",
		Long: `Add an object to raptor: you can create a box, a secret or add an item
to an existing secret as well`,
	}
	// Here you will define your flags and configuration settings.
	c.AddCommand(create.AddBoxCmd)
	c.AddCommand(create.AddSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
