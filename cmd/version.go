package cmd

import (
	"fmt"

	"github.com/mas2020-golang/cryptex/packages/utils"
	// "github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}

func NewVersionCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the application version",
		Long:  "Displays version information including Git commit and build date.",
		Run: func(cmd *cobra.Command, args []string) {
			// output.ActivityBox("Application info:")
			fmt.Printf("raptor info:\n")
			fmt.Printf("  %-11s %s\n", "Version:", utils.Version)
			fmt.Printf("  %-10s %s\n", "Git commit:", utils.GitCommit)
			fmt.Printf("  %-10s %s\n", "Build date:", utils.BuildDate)
		},
	}
	return newCmd
}