/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"log/slog"

	"github.com/atotto/clipboard"
	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	includeNumbers, includeLetters, includeSpecial bool
	length                                         int
	pwd                                            string
)

// boxCmd represents the box command
func NewPasswordCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "password",
		Aliases: []string{"pwd", "pass"},
		Short:   "Generate a random password",
		Long:    `Generate a random password and copy it in memory`,
		Example: `$ raptor create pwd`,
		Run: func(cmd *cobra.Command, args []string) {
			slog.Debug("create.newDecryptCmd()")
			var err error
			if pwd, err = security.GeneratePassword(length, includeNumbers, includeLetters, includeSpecial); err != nil {
				output.Error("", err.Error())
				return
			}
			err = clipboard.WriteAll(pwd)
			if err != nil {
				output.Error("", err.Error())
				return
			} else {
				utils.Success("password is stored in your clipboard")
			}
			slog.Debug("create.newDecryptCmd()", "password", pwd)
		},
	}
	// Here you will define your flags and configuration settings.
	c.Flags().BoolVarP(&includeNumbers, "includeNumbers", "n", true, "A boolean indicating whether to include numbers (0-9)")
	c.Flags().BoolVarP(&includeLetters, "includeLetters", "l", true, "A boolean indicating whether to include letters (a-z, A-Z)")
	c.Flags().BoolVarP(&includeSpecial, "includeSpecial", "s", true, "A boolean indicating whether to include special characters (&*{}'\"<>!$@)")
	c.Flags().IntVarP(&length, "length", "d", 10, "The char lenght (default 10)")

	return c
}
