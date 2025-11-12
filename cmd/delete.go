/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/delete"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
func newDeleteCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete",
		Short: "Delete a raptor object",
		Long:  `Delete a raptor object: secret`,
	}
	// Here you will define your flags and configuration settings.
	c.AddCommand(delete.DeleteSecretCmd)
	c.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "to get more information use the verbose mode")

	return c
}
