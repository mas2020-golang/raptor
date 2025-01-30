/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var key, boxName string

// boxCmd represents the box command
var AddSecretCmd = &cobra.Command{
	Use:     "secret <NAME>",
	Aliases: []string{"sr"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Create a new secret",
	Long:    `Create a new secret adding the one to the existing secret for the box`,
	Example: `$ cryptex secret add 'new-secret' --box test`,
	Run: func(cmd *cobra.Command, args []string) {
		add(args[0])
	},
}

func init() {
	AddSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
}

func add(name string) {
	// open the box
	boxPath, key, box, err := utils.OpenBox(boxName)
	utils.Check(err, "")
	// add the secret
	err = addSecret(name, box)
	utils.Check(err, "")
	fmt.Println()
	// save the box
	err = utils.SaveBox(boxPath, key, box)
	utils.Check(err, "")
	utils.Success(output.BoldS("box saved!"))
}

func addSecret(name string, box *utils.Box) error {
	if err := search(name, box); err != nil {
		return err
	}
	output.Bold("\n==> adding a new secret\n")
	output.RedOut("(to exit without saving press CTRL+C)\n")
	fmt.Println(output.GreenS(strings.Repeat("-", 35)))
	if box.Secrets == nil {
		box.Secrets = make([]*utils.Secret, 0)
	}
	// new secret
	s := utils.Secret{}
	s.Name = name
	// read from standard input
	r := bufio.NewReader(os.Stdin)
	fmt.Println(output.BlueS("Name:"), name)
	fmt.Printf("%s [%s]: ", output.BlueS("Version"), output.BoldS("1.0.0"))
	s.Version = utils.GetText(r)
	if len(s.Version) == 0 {
		s.Version = "1.0.0"
	}
	fmt.Print(output.BlueS("Login: "))
	s.Login = utils.GetText(r)
	fmt.Print(output.BlueS("Password: "))
	input, err := utils.ReadPassword("")
	utils.Check(err, "")
	if len(input) != 0 {
		fmt.Printf("\n%s [%s]: ", output.BlueS("Confirm pwd"), output.BoldS("xxx"))
		input2, err := utils.ReadPassword("")
		utils.Check(err, "")
		if input != input2 {
			fmt.Println()
			return fmt.Errorf("the pwd mismatched")
		}
		s.Pwd = input
	}
	fmt.Print(output.BlueS("\nUrl: "))
	s.Url = utils.GetText(r)
	fmt.Println(output.BlueS("Notes (to save type '>>' and press ENTER):"))
	s.Notes = utils.GetTextWithEsc(r)
	fmt.Println(output.BlueS(strings.Repeat("-", 35)))
	utils.GetText(r)
	fmt.Printf("Do you have other secres to add? [Y/n] ")
	answer := utils.GetText(r)
	if answer == "Y" {
		s.Others = make(map[string]string)
		for {
			fmt.Print(output.BlueS("Name: "))
			n := utils.GetText(r)
			fmt.Print(output.BlueS("Value: "))
			v := utils.GetText(r)
			s.Others[n] = v
			fmt.Printf("\nDo you have other secres to add? [Y/n] ")
			answer = utils.GetText(r)
			if answer != "Y" {
				break
			}
		}
	}
	s.LastUpdated = time.Now().Format(time.RFC3339)
	box.Secrets = append(box.Secrets, &s)
	return nil
}

// search goes into the secret and throws an error if a secret with the same
// name already exists
func search(name string, box *utils.Box) error {
	for _, s := range box.Secrets {
		if (*s).Name == name {
			return fmt.Errorf("a secret with the name %s already exists", name)
		}
	}
	return nil
}
