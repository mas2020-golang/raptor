/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package box

import (
	"fmt"
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your create",
	Long: `Create longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("box called")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
