/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package box

import (
	"github.com/spf13/cobra"
)

// boxCmd represents the box command
var BoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Manage the box",
	Long: `Manage the box: you can create, delete and list your boxes using this
command`,
}

func init() {
	// Here you will define your flags and configuration settings.
	BoxCmd.AddCommand(createCmd)
	BoxCmd.AddCommand(listCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	BoxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}