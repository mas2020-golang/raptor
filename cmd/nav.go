package cmd

import (
	"github.com/mas2020-golang/cryptex/cmd/nav"
	"github.com/spf13/cobra"
)

func newNavCmd() *cobra.Command {
	return nav.NewCmd()
}
