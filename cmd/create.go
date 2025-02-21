/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "os"

	"github.com/mas2020-golang/cryptex/cmd/create"
	// "github.com/mas2020-golang/cryptex/packages/utils"
	// "github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new", "cr"},
		Short:   "Add an object to raptor",
		Long: `Add an object to raptor: you can create a box, a secret or add an item
to an existing secret as well. (not available in interactive mode)`,
		Run: func(cmd *cobra.Command, args []string) {
			// if utils.BufferBox != nil {
			// 	output.Error("", "create is not available in interactive mode")
			// 	os.Exit(1)
			// }
		},
	}
	// Here you will define your flags and configuration settings.
	c.AddCommand(create.AddBoxCmd)
	c.AddCommand(create.AddSecretCmd)
	c.AddCommand(create.NewPasswordCmd())

	return c
}
