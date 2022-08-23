/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"fmt"
	"strings"
	"time"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
)

var unsecure bool

// boxCmd represents the box command
var printCmd = &cobra.Command{
	Use:     "print <NAME>",
	Aliases: []string{"pr"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Print the info of a secret",
	Long: `Print all the info related to the secret. If you specify --unsecure flag you will get also the sensitive
information in clear on the screen (use it carefully)`,
	Example: `$ cryptex secret print foo --box test // to print the info of the foo secret`,
	Run: func(cmd *cobra.Command, args []string) {
		print(args[0])
	},
}

func init() {
	printCmd.Flags().BoolVarP(&unsecure, "unsecure", "u", false, "If passed sensitive info are shown")
}

func print(name string) {
	// open the box
	_, err := openBox()
	utils.Check(err, "")
	s, err := getSecret(name)
	utils.Check(err, "")
	err = showToStdOut(s, unsecure)
	utils.Check(err, "")
}

// getSecret searches the secret into the box.
func getSecret(name string) (*protos.Secret, error) {
	if len(box.Secrets) == 0 {
		return nil, fmt.Errorf("no secret found in the box")
	}
	for _, secret := range box.Secrets {
		if secret.Name == name {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("no secret found in the box")
}

func showToStdOut(s *protos.Secret, unsecure bool) error {
	// load the local timezone
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}
	lastUpdated := s.LastUpdated.AsTime().In(loc).Format("Jan 2 15:04 2006 MST")
	fmt.Println(utils.GreenS(strings.Repeat("-", 35)))
	fmt.Printf("%s %s\n", utils.BlueS("Version:"), s.Version)
	fmt.Printf("%s %s\n", utils.BlueS("Login:"), s.Login)
	if unsecure {
		fmt.Printf("%s %s\n", utils.BlueS("Pwd:"), s.Pwd)
	} else {
		fmt.Printf("%s %s\n", utils.BlueS("Pwd:"), "xxx")
	}
	fmt.Printf("%s %s\n", utils.BlueS("Url:"), s.Url)
	fmt.Printf("%s\n%s\n", utils.BlueS("\nNotes:"), s.Notes)
	fmt.Println(utils.BlueS(strings.Repeat("-", 35)))
	if s.Others != nil && len(s.Others) > 0 {
		utils.BoldOut("Items:\n")
		for k, v := range s.Others {
			if unsecure {
				fmt.Printf("%-2s.%s -> %s\n", "", utils.BlueS(k), v)
			} else {
				fmt.Printf("%-2s.%s\n", "", utils.BlueS(k))
			}
		}
	}
	fmt.Printf("%s %s\n", utils.BoldS("\nUpdated on:"), lastUpdated)
	return nil
}
