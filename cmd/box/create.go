/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package box

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	owner string
)

// boxCmd represents the box command
var createCmd = &cobra.Command{
	Use:   "create <NAME>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a new box",
	Long: `Create a new box and the .cryptex folder structure in case
it doesn't exist yet`,
	Example: `$ cryptex box create 'test' --owner me`,
	Run: func(cmd *cobra.Command, args []string) {
		create(args)
	},
}

func init() {
	createCmd.Flags().StringVarP(&owner, "owner", "o", "", "The owner of the box (e.g. --owner bar)")
}

func create(args []string) {
	createHomeFolder()
	createBox(args[0], owner)
}

func createHomeFolder() {
	// check the folder .cryptex
	home, err := os.UserHomeDir()
	utils.Check(err, "")
	_, err = os.Stat(path.Join(home, ".cryptex"))
	if err != nil {
		if os.IsNotExist(err) {
			// create the directory structure
			err = os.MkdirAll(path.Join(home, ".cryptex", "boxes"), 0777)
			utils.Check(err, "")
			utils.Success(fmt.Sprintf("home folder created in %s", home))
		}
	}
}

func createBox(name, owner string) {
	home, err := os.UserHomeDir()
	utils.Check(err, "")
	b := protos.Box{
		Name:        name,
		Owner:       owner,
		Version:     "1",
		LastUpdated: timestamppb.Now(),
	}

	// Write the new address book back to disk.
	_, err = os.Stat(path.Join(home, ".cryptex", "boxes", b.Name))
	if err == nil {
		utils.Error("the box already exists, try with a different name")
		os.Exit(1)
	}
	out, err := proto.Marshal(&b)
	utils.Check(err, "failed to encode the box")
	// ask for the password
	key, err := utils.AskForPassword("Password: ", true)
	utils.Check(err, "")
	// encrypt the box
	encOut, err := security.EncryptBox(out, key)
	utils.Check(err, "")
	// write the box into the disk
	err = ioutil.WriteFile(path.Join(home, ".cryptex", "boxes", b.Name), encOut, 0644)
	utils.Check(err, "failed to write the box")
	fmt.Println()
	utils.Success(fmt.Sprintf("Box %q created successfully!", name))
}
