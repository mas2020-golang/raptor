/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/mas2020-golang/cryptex/packages/security"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	Example: `$ cryptex create box 'test' --owner me`,
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

// TODO: refactor the box creation using the YAML box (look at the `YAML info.md` file)
func createBox(name, owner string) error {
	// get the folder box
	boxPath, err := utils.GetFolderBox()
	if err != nil {
		return fmt.Errorf("problem to determine the folder box: %v", err)
	}

	b := utils.Box{
		Name:        name,
		Owner:       owner,
		Version:     "1",
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	// check if the box already exists
	_, err = os.Stat(path.Join(boxPath, b.Name))
	if err == nil {
		return fmt.Errorf("problem to determine the folder box: %v", err)
	}

	out, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to encode the box: %v", err)
	}
	fmt.Printf("box creation YAML details:\n%s", string(out)) //TODO: remove this line

	// ask for the password
	key, err := utils.AskForPassword("Password: ", true)
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
