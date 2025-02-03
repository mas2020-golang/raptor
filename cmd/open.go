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

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	pwd string
)

func newOpenCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "open <BOX-NAME>",
		Aliases: []string{"op", "open"},
		// Args:    cobra.MinimumNArgs(1),
		Short:   "Open a box in interactive mode",
		Long:    ``,
		Example: `$ raptor open 'test'`,
		Run: func(cmd *cobra.Command, args []string) {
			boxName := ""
			if len(args) != 0 {
				boxName = args[0]
			}

			if err := interactiveOpen(boxName); err != nil {
				output.Error("", err.Error())
				os.Exit(1)
			}
		},
	}
	// Here you will define your flags and configuration settings.
	c.Flags().StringVarP(&pwd, "pwd", "p", "", "pwd to open the box (use ONLY FOR DEBUG MODE)")

	return c
}

func interactiveOpen(boxName string) error {
	// open the box
	// open the box if not interactive mode
	exit := false
	_, _, box, err := utils.OpenBox(boxName, pwd)
	utils.BufferBox = box
	utils.Check(err, "")
	output.Success("Box is ready for you!")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("raptor> ")
		if !scanner.Scan() {
			// Handle error or end of input
			break
		}
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "quit", "q", "exit", "bye":
			fmt.Println("see you for the next secret to whisper...")
			exit = true
		}
		if exit {
			break
		}

		switch input {
		case "clear", "cl", "wipe", "clean":
			clearScreen()
			continue
		}

		// fetch the command
		clearScreen()
		os.Args = []string{"raptor"}
		os.Args = append(os.Args, strings.Split(input, " ")...)
		if err := rootCmd.Execute(); err != nil {
			output.Error("", err.Error())
		}

		//showOutput()
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
