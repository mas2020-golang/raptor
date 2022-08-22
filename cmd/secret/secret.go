/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/spf13/cobra"
)

var (
	boxName string
	box     *protos.Box
)

// boxCmd represents the box command
var SecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage the secrets",
	Long: `Manage the secrets: you can create, delete and list your secrets in the 
specified box using this command.`,
}

func init() {
	// Here you will define your flags and configuration settings.
	SecretCmd.AddCommand(addCmd)
	SecretCmd.AddCommand(listCmd)
	SecretCmd.AddCommand(getCmd)
	SecretCmd.AddCommand(printCmd)
	SecretCmd.AddCommand(editCmd)
	SecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
	//SecretCmd.MarkPersistentFlagRequired("box")
}
