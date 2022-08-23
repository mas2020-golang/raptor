/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get the sensitive info from a secret",
	Long: `Get the sensitive info from a secret. You can refer to the data as:
- <SECRET_NAME>: retrieves the root pwd for the secret
- <SECRET_NAME>.<ITEM_NAME>: retrieves the ITEM_NAME in the items collection`,
	Example: `$ cryptex secret get foo --box test // to retrieve the pwd of the foo secret
$ cryptex secret get foo.test --box test // to retrieve the test secret item of the foo secret`,
	Run: func(cmd *cobra.Command, args []string) {
		get(args[0])
	},
}

func init() {
}

func get(name string) {
	// open the box
	_, err := openBox()
	utils.Check(err, "")
	s, err := searchSecretPwd(name)
	utils.Check(err, "")
	if len(s) == 0 {
		utils.Warning(fmt.Sprintf("no secret %q found in %s", name, boxName))
		return
	}
	// copy the secret into the clipboard
	err = execCmd(s)
	utils.Check(err, "")
	fmt.Println()
	utils.Success(utils.BoldS("the secret is in your clipboard"))
}

// searchSecretPwd searches for the secret and the value:
// e.g. foo, foo.test
func searchSecretPwd(name string) (value string, err error) {
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
			if secretName == s.Name {
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
//  - pbcopy: for Mac and Linux
//  - clip: for Windows
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
