/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	owner string
)

var OpenBoxCmd = &cobra.Command{
	Use:     "open <BOX-NAME>",
	Aliases: []string{"op", "open"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Open a box",
	Long:    `xxx xxx xxx`,
	Example: `$ raptor open 'test'`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := interactiveOpen(); err != nil {
			output.Error("", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	OpenBoxCmd.Flags().StringVarP(&owner, "owner", "o", "", "The owner of the box (e.g. --owner bar)")
}

func interactiveOpen() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			// Handle error or end of input
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" || input == "q" {
			fmt.Println("Exiting...")
			break
		}
		switch input {
		case "clear", "cl", "wipe", "clean":
			clearScreen()
			continue
		}
		clearScreen()
		showOutput()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showOutput() {
	for i := 0; i < 10; i++ {
		fmt.Println(`
Box: NAME
Last update: DATE`)
	}
	fmt.Println()
}
