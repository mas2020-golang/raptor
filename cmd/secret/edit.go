/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// boxCmd represents the box command
var editCmd = &cobra.Command{
	Use:   "edit <NAME>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Edit an existing secret",
	Long: `Edit a secret by name:
The <NAME> argument is in the following format:
  - 'secret name.item': to update an item
  - 'secret name.pwd': to update only the pwd
  - 'secret name': to update every elements
`,
	Example: `$ cryptex secret edit 'new-secret' --box test`,
	Run: func(cmd *cobra.Command, args []string) {
		edit(args[0])
	},
}

func init() {
}

func edit(name string) {
	// open the box
	boxPath, err := openBox()
	utils.Check(err, "")
	// add the secret
	err = editSecret(name)
	utils.Check(err, "")
	fmt.Println()
	// save the box
	err = saveBox(boxPath)
	utils.Check(err, "")
	utils.Success(utils.BoldS("box saved!"))
}

func editSecret(name string) error {
	// get the secret to edit
	s := findSecret(name)
	if s == nil {
		return fmt.Errorf("the secret doesn't exist")
	}

	utils.Note(utils.BoldS("\npress ENTER without typing to skip the field"))
	utils.RedOut("(to exit without saving type CTRL+C)\n")
	fmt.Println(strings.Repeat("-", 35))
	// read from standard input
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%s]: ", utils.BlueS("Name"), utils.BoldS(s.Name))
	input := utils.GetText(r)
	if len(input) != 0 {
		s.Name = input
	}
	fmt.Printf("%s [%s]: ", utils.BlueS("Version"), utils.BoldS(s.Version))
	input = utils.GetText(r)
	if len(input) != 0 {
		s.Version = input
	}
	fmt.Printf("%s [%s]: ", utils.BlueS("Login"), utils.BoldS(s.Login))
	input = utils.GetText(r)
	if len(input) != 0 {
		s.Login = input
	}
	fmt.Printf("%s [%s]: ", utils.BlueS("Pwd"), utils.BoldS("xxx"))
	input, err := utils.ReadPassword("")
	utils.Check(err, "")
	if len(input) != 0 {
		fmt.Printf("\n%s [%s]: ", utils.BlueS("Confirm pwd"), utils.BoldS("xxx"))
		input2, err := utils.ReadPassword("")
		utils.Check(err, "")
		if input != input2 {
			fmt.Println()
			return fmt.Errorf("the pwd mismatched")
		}
		s.Pwd = input
	}
	fmt.Printf("\n%s [%s]: ", utils.BlueS("Url"), utils.BoldS(s.Url))
	input = utils.GetText(r)
	if len(input) != 0 {
		s.Url = input
	}
	fmt.Printf("%s:\n%s\n", utils.BlueS("Notes"), s.Notes)
	fmt.Println(utils.BlueS(strings.Repeat("-", 35)))
	fmt.Printf("Do you want to change the notes? [Y/n] ")
	answer := utils.GetText(r)
	if answer == "Y" {
		fmt.Println(utils.BlueS("\nNew Notes: (to save type '>>' and press ENTER)"))
		input = utils.GetTextWithEsc(r)
		if len(input) != 0 && input != "ERROR!" {
			s.Notes = input
		}
		fmt.Println(utils.BlueS(strings.Repeat("-", 35)))
		utils.GetText(r)
	}
	if s.Others != nil {
		utils.BlueOut(fmt.Sprintf("\nItems:\n"))
		for k, v := range s.Others {
			fmt.Printf("%s [%s]: ", k, v)
			input = utils.GetText(r)
			if len(input) != 0 {
				s.Others[k] = input
			}
		}
	}

	s.LastUpdated = timestamppb.Now()
	return nil
}

// findSecret searches for the secret into the box and returns the one corresponding or nil
// value
func findSecret(name string) *protos.Secret {
	for _, s := range box.Secrets {
		if (*s).Name == name {
			return s
		}
	}
	return nil
}
