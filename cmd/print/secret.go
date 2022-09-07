/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package print

import (
	"fmt"
	"strings"
	"time"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	unsecure bool
	boxName  string
)

// boxCmd represents the box command
var PrintSecretCmd = &cobra.Command{
	Use:     "secret <NAME>",
	Aliases: []string{"sr"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Print the info of a secret",
	Long: `Print all the info related to the secret. If you specify --unsecure flag you will get also the sensitive
information in clear on the screen (use it carefully)`,
	Example: `$ cryptex secret print foo --box test // to print the info of the foo secret`,
	Run: func(cmd *cobra.Command, args []string) {
		print(args[0], cmd)
	},
}

func init() {
	PrintSecretCmd.Flags().BoolVarP(&unsecure, "unsecure", "u", false, "If passed sensitive info are shown")
	PrintSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
}

func print(name string, cmd *cobra.Command) {
	// open the box
	boxPath, _, box, err := utils.OpenBox(boxName)
	utils.Check(err, "")
	s, err := getSecret(name, box)
	utils.Check(err, "")
	err = showToStdOut(s, unsecure, cmd, boxPath)
	utils.Check(err, "")
}

// getSecret searches the secret into the box.
func getSecret(name string, box *protos.Box) (*protos.Secret, error) {
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

func showToStdOut(s *protos.Secret, unsecure bool, cmd *cobra.Command, boxPath string) error {
	// load the local timezone
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}
	lastUpdated := s.LastUpdated.AsTime().In(loc).Format("Jan 2 15:04 2006 MST")
	fmt.Println(output.GreenS(strings.Repeat("-", 35)))
	fmt.Printf("%s %s\n", output.BlueS("Version:"), s.Version)
	fmt.Printf("%s %s\n", output.BlueS("Login:"), s.Login)
	if unsecure {
		fmt.Printf("%s %s\n", output.BlueS("Pwd:"), s.Pwd)
	} else {
		fmt.Printf("%s %s\n", output.BlueS("Pwd:"), "xxx")
	}
	fmt.Printf("%s %s\n", output.BlueS("Url:"), s.Url)
	fmt.Printf("%s\n%s\n", output.BlueS("\nNotes:"), s.Notes)
	fmt.Println(output.BlueS(strings.Repeat("-", 35)))
	if s.Others != nil && len(s.Others) > 0 {
		output.Bold("Items:")
		for k, v := range s.Others {
			if unsecure {
				fmt.Printf("%-2s.%s -> %s\n", "", output.BlueS(k), v)
			} else {
				fmt.Printf("%-2s.%s\n", "", output.BlueS(k))
			}
		}
	}
	fmt.Printf("%s %s\n", output.BoldS("\nUpdated on:"), lastUpdated)
	v, _ := cmd.Parent().Flags().GetBool("verbose")
	if v {
		fmt.Printf("\n%s\n", output.YellowS("additional info:"))
		output.InfoBox(fmt.Sprintf("secret read from the %s box\n", output.BlueS(boxPath)))
	}

	return nil
}
