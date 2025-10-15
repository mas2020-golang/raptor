/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

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
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newEditCmd())
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newNavCmd())
	rootCmd.AddCommand(newPrintCmd())
	rootCmd.AddCommand(newOpenCmd())
	rootCmd.AddCommand(newEncryptCmd())
	rootCmd.AddCommand(newDecryptCmd())
	rootCmd.AddCommand(newInfoCmd())

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Give more information about the command execution")
}
