/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/console"
	"github.com/spf13/cobra"
)

var ()

func newDecryptCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "decrypt <FILE|FOLDER>",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"en", "encrypt"},
		// Args:    cobra.MinimumNArgs(1),
		Short: "Decrypt a file or a folder",
		Long: `The decryption is accepting a file or a folder. The command will automatically delete
each file in in the path.`,
		Example: `$ raptor decrypt /test/file`,
		Run: func(cmd *cobra.Command, args []string) {
			slog.Debug("encrypt run", "path", args[0])
			if err := decrypt(args[0]); err != nil {
				if !errors.Is(err, security.ErrInvalidFile) {
					console.Error(err.Error(), true)
					//output.Error("", err.Error())
				}
			} else {
				console.OK("Decryption succeded")
			}
		},
	}
	// Here you will define your flags and configuration settings.
	//	c.Flags().StringVarP(&pwd, "pwd", "p", "", "pwd to open the box (use ONLY FOR DEBUG MODE)")

	return c
}

func decrypt(path string) error {
	// does the path exists
	info, err := os.Stat(path)
	utils.Verbosity(fmt.Sprintf("decryption starting on path %s", path), verbose)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the given path does not exist")
		} else {
			return fmt.Errorf("error accessing the path %s: %v", path, err)
		}
	}

	passphrase, err := utils.AskForPassword("Password: ", false)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	// encrypt the file or the folder
	if info.IsDir() {
		return security.DecryptDirectory(path, passphrase)
	} else {
		return security.DecryptFile(path, passphrase)
	}
}
