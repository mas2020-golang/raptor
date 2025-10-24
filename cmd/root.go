/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose    bool
	createCmd  *cobra.Command
	listCmd    *cobra.Command
	getCmd     *cobra.Command
	editCmd    *cobra.Command
	printCmd   *cobra.Command
	openCmd    *cobra.Command
	encryptCmd *cobra.Command
	decryptCmd *cobra.Command
	infoCmd    *cobra.Command
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "raptor <command> [flags]",
	Short: "Use raptor to keep safe your personal information",
	Long: `Raptor is a Go-based application designed to securely store your personal information within an encrypted "box."
By leveraging robust encryption techniques, SecureBox ensures that your sensitive data remains confidential and accessible
only to you.

Features:
    - Strong Encryption: Utilizes AES (Advanced Encryption Standard) for encrypting your data, ensuring high levels of security.
    slingacademy.com

    - Password Protection: Access your encrypted box using a password, providing an additional layer of security.

    - Data Integrity: Ensures that your personal information remains intact and unaltered during storage and retrieval.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "boxes", Title: "Box Management"})
	rootCmd.AddGroup(&cobra.Group{ID: "encryption", Title: "Secret Operations"})

	createCmd = newCreateCmd()
	listCmd = newListCmd()
	getCmd = newGetCmd()
	editCmd = newEditCmd()
	printCmd = newPrintCmd()
	openCmd = newOpenCmd()
	encryptCmd = newEncryptCmd()
	decryptCmd = newDecryptCmd()
	infoCmd = newInfoCmd()
	
	listCmd.GroupID = "boxes"
	createCmd.GroupID = "boxes"
	getCmd.GroupID = "boxes"
	editCmd.GroupID = "boxes"
	openCmd.GroupID = "boxes"
	printCmd.GroupID = "boxes"
	encryptCmd.GroupID = "encryption"
	decryptCmd.GroupID = "encryption"

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(infoCmd)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Give more information about the command execution")
}
