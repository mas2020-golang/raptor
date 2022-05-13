/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package secret

import (
	"bufio"
	"fmt"
	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var key string

// boxCmd represents the box command
var addCmd = &cobra.Command{
	Use:     "add <NAME>",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Create a new secret",
	Long:    `Create a new secret adding the one to the existing secret for the box`,
	Example: `$ cryptex secret add 'new-secret' --box test`,
	Run: func(cmd *cobra.Command, args []string) {
		add(args[0])
	},
}

func init() {
}

func add(name string) {
	// open the box
	boxPath, err := openBox(boxName)
	utils.Check(err, "")
	// add the secret
	err = addSecret(name)
	utils.Check(err, "")
	fmt.Println()
	// save the box
	err = saveBox(boxPath)
	utils.Check(err, "")
	utils.Success(utils.BoldS("box saved!"))
}

func openBox(name string) (string, error) {
	// search the CRYPTEX_BOX env if name is empty
	if len(name) == 0{
		name = os.Getenv("CRYPTEX_BOX")
		if len(name) == 0{
			return "", fmt.Errorf("--box args is not given and the env var CRYPTEX_BOX is empty")
		}
	}
	// check the folder .cryptex
	home, err := os.UserHomeDir()
	utils.Check(err, "")
	// read the file in the home dir
	path := path.Join(home, ".cryptex", "boxes", name)
	// Read the existing address book.
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading the file box in %s: %v", path, err)
	}

	// ask for the password
	key, err = utils.AskForPassword(false)
	utils.Check(err, "")
	// encrypt the box
	decIn, err := security.DecryptBox(in, key)
	utils.Check(err, "")

	box = &protos.Box{}
	err = proto.Unmarshal(decIn, box)
	if err != nil {
		return "", fmt.Errorf("failed to read the box: %v. Maybe an incorrect pwd?", err)
	}
	return path, nil
}

func addSecret(name string) error {
	if err := search(name); err != nil{
		return err
	}
	utils.BoldOut("==> add a new secret (only for the NOTES: to save the press CTRL+D)\n")
	utils.RedOut("(to exit without saving type CTRL+C)\n")
	fmt.Println(strings.Repeat("_", 45))
	if box.Secrets == nil {
		box.Secrets = make([]*protos.Secret, 0)
	}
	// new secret
	s := protos.Secret{}
	s.Name = name
	// read from standard input
	r := bufio.NewReader(os.Stdin)
	fmt.Println(utils.BoldS("Name: "), name)
	fmt.Printf(utils.BoldS("Version: "))
	s.Version = utils.GetText(r)
	if len(s.Version) == 0 {
		fmt.Println()
	}
	fmt.Printf(utils.BoldS("Login: "))
	s.Login = utils.GetText(r)
	fmt.Printf("Password: ")
	s.Pwd = utils.GetText(r)
	fmt.Printf("Url: ")
	s.Url = utils.GetText(r)
	fmt.Printf("Notes: ")
	s.Notes = utils.GetText(r)
	fmt.Printf("Do you have other secres to add? [Y/n] ")
	answer := utils.GetText(r)
	if answer == "Y" {
		s.Others = make(map[string]string)
		for {
			fmt.Printf("Name: ")
			n := utils.GetText(r)
			fmt.Printf("Value: ")
			v := utils.GetText(r)
			s.Others[n] = v
			fmt.Printf("Do you have other secres to add? [Y/n] ")
			answer = utils.GetText(r)
			if answer != "Y" {
				break
			}
		}
	}
	s.LastUpdated = timestamppb.Now()
	box.Secrets = append(box.Secrets, &s)
	return nil
}

func saveBox(path string) error {
	out, err := proto.Marshal(box)
	if err != nil {
		return fmt.Errorf("failed to encode the box: %v", err)
	}
	// encrypt the box
	encOut, err := security.EncryptBox(out, key)
	utils.Check(err, "")
	if err := ioutil.WriteFile(path, encOut, 0644); err != nil {
		return fmt.Errorf("failed to write the box: %v", err)
	}
	return nil
}

// search goes into the secret and throws an error if a secret with the same
// name already exists
func search(name string) error {
	for _, s := range box.Secrets{
		if (*s).Name == name{
			return fmt.Errorf("a secret with the name %s already exists", name)
		}
	}
	return nil
}