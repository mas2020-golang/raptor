/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	pwd string
)

func newOpenCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "open [BOX-NAME]",
		Aliases: []string{"op", "open"},
		// Args:    cobra.MinimumNArgs(1),
		Short: "Open a box in interactive mode",
		Long: `Interactive mode is the simplest way to get the secrets giving you the ability to keep your box open.
You can open a box giving:
- the fs path (e.g. /Users/test/mybox)
- the name: raptor will search the box in the CRYPTEX_FOLDER env variable (if set) or in the $HOME/.cryptex folder

If you omit the name raptor will try to fetch the CRYPTEX_BOX env variable value.
If any of the previous checks doesn't give the right path raptor will throw an error.`,
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
	inputChan := make(chan string)
	doneChan := make(chan bool)
	secStringTimeout := os.Getenv("RAPTOR_TIMEOUT_SEC")
	timeout := 600 * time.Second
	if len(secStringTimeout) > 0 {
		// Convert the string to an integer
		secsTimeout, err := strconv.Atoi(secStringTimeout)
		if err != nil {
			output.Warning("", fmt.Sprintf("Error converting RAPTOR_LOGLEVEL to int: %v", err))
			return nil
		}
		timeout = time.Duration(secsTimeout) * time.Second
	}
	output.InfoBox(fmt.Sprintf("RAPTOR_TIMEOUT_SEC is set to %v", timeout))

	// open the box
	exit := false
	_, _, box, err := utils.OpenBox(boxName, pwd)
	utils.BufferBox = box
	utils.Check(err, "")
	output.Success("Box is ready for you!")

	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			slog.Debug("waiting input...")
			fmt.Print("raptor> ")
			if !scanner.Scan() {
				// Handle error or end of input
				break
			}
			input := strings.TrimSpace(scanner.Text())
			if err != nil {
				output.Error("", fmt.Sprintf("error reading input: %v", err))
				close(inputChan)
				return
			}

			if err := scanner.Err(); err != nil {
				return
			}
			slog.Debug("writing into the channel", "input", input)
			inputChan <- input
			slog.Debug("waiting for done channel")
			<-doneChan // Wait for processing to complete
		}
	}()

	for {
		select {
		case input, ok := <-inputChan:
			slog.Debug("input read from the channel", "input", input)
			if !ok {
				output.Warning("", "input channel closed, exiting")
				return nil
			}
			switch input {
			case "quit", "q", "exit", "bye":
				fmt.Println("see you for the next secret to whisper...")
				exit = true
			}
			if exit {
				return nil
			}

			switch input {
			case "clear", "cl", "wipe", "clean":
				clearScreen()
				doneChan <- true // Signal that processing is complete
				continue
			}

			// fetch the command
			clearScreen()
			os.Args = []string{"raptor"}
			os.Args = append(os.Args, strings.Split(input, " ")...)
			slog.Debug("executing the command", "args", os.Args)
			rootCmd.Execute()
			doneChan <- true // Signal that processing is complete
		case <-time.After(timeout):
			fmt.Printf("no input received for %v, exiting\n", timeout)
			return nil
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
