/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/get"
	"github.com/spf13/cobra"
)

var (
	boxName string
)

func newGetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Get an object from raptor (*available in interactive mode only)",
		Long: `Get an object to raptor: you can get a box, a secret or an item
of an existing secret as well`,
	}
	// Here you will define your flags and configuration settings.
	c.AddCommand(get.GetSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
