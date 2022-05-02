/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
var SecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage the secrets",
	Long: `Manage the secrets: you can create, delete and list your secrets in a specified box using this
command`,
}

func init() {
	// Here you will define your flags and configuration settings.
	SecretCmd.AddCommand(addCmd)
}