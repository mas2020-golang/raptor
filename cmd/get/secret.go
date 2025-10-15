/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"

	"github.com/mas2020-golang/cryptex/cmd/internal/secretutil"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"

	// "golang.design/x/clipboard"
	"github.com/atotto/clipboard"
)

var boxName string

// boxCmd represents the box command
var GetSecretCmd = &cobra.Command{
	Use:     "secret <NAME>",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"sr"},
	Short:   "Get the sensitive data from a secret",
	Long: `Get the sensitive data from a secret. You can refer to the data as:
- <SECRET_NAME>: retrieves the root sensitive data for the secret
- <SECRET_NAME>.<ITEM_NAME>: retrieves the ITEM_NAME sensitive data in the items collection`,
	Example: `$ raptor get secret foo --box test // to retrieve the pwd of the foo secret
$ raptor get secret foo.test --box test // to retrieve the test secret item of the foo secret`,
	Run: func(cmd *cobra.Command, args []string) {
		get(args[0])
	},
}

func init() {
	GetSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
}

func get(name string) {
	result, _, err := secretutil.Lookup(boxName, name)
	utils.Check(err, "")
	if result == nil || len(result.Value) == 0 {
		output.Warning("", fmt.Sprintf("no secret %q found in %s", name, boxName))
		return
	}
	// copy the secret into the clipboard
	err = clipboard.WriteAll(result.Value)
	if err != nil {
		output.Error("", err.Error())
	}
	fmt.Println()
	utils.Success(output.BoldS("the secret is in your clipboard"))
}
