/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package delete

import (
	"fmt"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var boxName string

// DeleteSecretCmd represents the delete secret command
var DeleteSecretCmd = &cobra.Command{
	Use:     "secret <NAME>",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"sr"},
	Short:   "Delete an existing secret",
	Long: `Delete a secret by name from the specified box.
The secret will be permanently removed from the encrypted box.`,
	Example: `$ raptor delete secret 'my-secret' --box test`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteSecret(args[0])
	},
}

func init() {
	DeleteSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to delete the secret")
}

func deleteSecret(name string) {
	// open the box
	boxPath, key, box, err := utils.OpenBox(boxName, "")
	utils.Check(err, "")

	// find and delete the secret
	deleted, err := removeSecret(name, box)
	if err != nil {
		output.Error("", err.Error())
		return
	}

	if !deleted {
		output.Warning("", fmt.Sprintf("no secret %q found in box %s", name, boxPath))
		return
	}

	fmt.Println()
	// save the box
	err = utils.SaveBox(boxPath, key, box)
	utils.Check(err, "")
	utils.Success(output.BoldS("secret deleted and box saved!"))
}

// removeSecret searches for the secret in the box and removes it if found
// Returns true if the secret was found and removed, false otherwise
func removeSecret(name string, box *utils.Box) (bool, error) {
	if box.Secrets == nil {
		return false, nil
	}

	for i, s := range box.Secrets {
		if s.Name == name {
			// Remove the secret by slicing
			box.Secrets = append(box.Secrets[:i], box.Secrets[i+1:]...)
			return true, nil
		}
	}

	return false, nil
}
