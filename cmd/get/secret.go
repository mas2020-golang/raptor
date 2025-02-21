/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

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
	// open the box
	_, _, box, err := utils.OpenBox(boxName, "")
	utils.Check(err, "")
	s, err := searchSecretPwd(name, box)
	utils.Check(err, "")
	if len(s) == 0 {
		output.Warning("", fmt.Sprintf("no secret %q found in %s", name, boxName))
		return
	}
	// copy the secret into the clipboard
	// Initialize the clipboard
	// err = clipboard.Init()
	// if err != nil {
	// 	panic(err)
	// }
	// Write text to the clipboard
	err = clipboard.WriteAll(s)
	if err != nil {
		output.Error("", err.Error())
	}
	// err = execCmd(s)
	// utils.Check(err, "")
	fmt.Println()
	utils.Success(output.BoldS("the secret is in your clipboard"))
}

// searchSecretPwd searches for the secret and the value:
// e.g. foo, foo.test
func searchSecretPwd(name string, box *utils.Box) (value string, err error) {
	var (
		secretName, secretItem string
	)
	elems := strings.Split(name, ".")
	secretName = elems[0]
	if len(elems) > 1 {
		secretItem = elems[1]
	}
	if box.Secrets != nil {
		for _, s := range box.Secrets {
			// fmt.Printf("secret in the box: %s\n", s.Name)
			// fmt.Printf("secret name given: %s\n", secretName)
			if secretName == s.Name {
				// fmt.Printf("found the secret name %s\n", secretName)
				if len(secretItem) > 0 {
					if len(s.Others) > 0 {
						for k, v := range s.Others {
							if k == secretItem {
								value = v
							}
						}
					}
				} else {
					value = s.Pwd
				}
			}
		}
	}
	return
}

// execCmd the command passing arg as the standard input. The command that will be executed is:
//   - pbcopy: for Mac and Linux
//   - clip: for Windows
func execCmd(arg string) error {
	cmdName := "pbcopy"
	os := runtime.GOOS
	switch os {
	case "windows":
		cmdName = "clip"
	case "darwin":
		cmdName = "pbcopy"
	case "linux":
		cmdName = "pbcopy"
	default:
		cmdName = "pbcopy"
	}
	cmd := exec.Command(cmdName)
	cmd.Stdin = strings.NewReader(arg)
	return cmd.Run()
}
