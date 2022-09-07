/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mas2020-golang/cryptex/packages/protos"
	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	owner string
)

// boxCmd represents the box command
var AddBoxCmd = &cobra.Command{
	Use:     "box <NAME>",
	Aliases: []string{"bo", "box"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Create a new box",
	Long: `Create a new box and the .cryptex folder structure in case
it doesn't exist yet`,
	Example: `$ cryptex box create 'test' --owner me`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := create(args); err != nil {
			output.Error("", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	AddBoxCmd.Flags().StringVarP(&owner, "owner", "o", "", "The owner of the box (e.g. --owner bar)")
}

func create(args []string) error {
	var err error
	if err = createHomeFolder(); err != nil {
		return err
	}

	return createBox(args[0], owner)
}

func createHomeFolder() error {
	// get the folder box
	boxPath, err := utils.GetFolderBox()
	utils.Check(err, "problem to determine th folder box")
	_, err = os.Stat(boxPath)
	if err != nil {
		if os.IsNotExist(err) {
			// create the directory structure
			err = os.MkdirAll(boxPath, 0777)
			utils.Check(err, "")
			utils.Success(fmt.Sprintf("folder created in %s", boxPath))
		}
	}
	return err
}

func createBox(name, owner string) error {
	// get the folder box
	boxPath, err := utils.GetFolderBox()
	if err != nil {
		return fmt.Errorf("problem to determine the folder box: %v", err)
	}

	b := protos.Box{
		Name:        name,
		Owner:       owner,
		Version:     "1",
		LastUpdated: timestamppb.Now(),
	}

	// check if the box already exists
	_, err = os.Stat(path.Join(boxPath, b.Name))
	if err == nil {
		return fmt.Errorf("problem to determine the folder box: %v", err)
	}

	out, err := proto.Marshal(&b)
	if err != nil {
		return fmt.Errorf("failed to encode the box: %v", err)
	}
	// ask for the password
	key, err := utils.AskForPassword("Password: ", true, 6)
	if err != nil {
		return err
	}
	// encrypt the box
	encOut, err := security.EncryptBox(out, key)
	if err != nil {
		return err
	}
	// write the box into the disk
	err = ioutil.WriteFile(path.Join(boxPath, b.Name), encOut, 0644)
	if err != nil {
		return fmt.Errorf("failed to write the box: %v", err)
	}
	fmt.Println()
	utils.Success(fmt.Sprintf("Box %q created successfully!", name))
	return nil
}
