package nav

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/mas2020-golang/cryptex/internal/secretutil"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

// NewCmd creates the "nav" command that opens the secret URL (if any) and copies
// the secret password to the clipboard.
func NewCmd() *cobra.Command {
	var boxName string

	cmd := &cobra.Command{
		Use:   "nav <NAME>",
		Args:  cobra.MinimumNArgs(1),
		Short: "Open the secret URL and copy its password",
		Long: `Open the secret URL (when present) using the default browser and
copy the secret password to the clipboard.`,
		Example: `$ raptor nav foo --box test // open foo secret URL and copy the password
$ raptor nav foo.bar --box test // open the foo secret URL and copy the password`,
		Run: func(cmd *cobra.Command, args []string) {
			runNav(boxName, args[0])
		},
	}

	cmd.Flags().StringVarP(&boxName, "box", "b", "", "The name of the box where the secret is stored")

	return cmd
}

func runNav(boxName, name string) {
	result, _, err := secretutil.Lookup(boxName, name)
	utils.Check(err, "")

	if result == nil || result.Secret == nil {
		output.Warning("", fmt.Sprintf("no secret %q found in %s", name, boxName))
		return
	}

	// When the user refers to an item, ensure it exists.
	if result.Item != "" && len(result.Value) == 0 {
		output.Warning("", fmt.Sprintf("no secret %q found in %s", name, boxName))
		return
	}

	secretPwd := result.Secret.Pwd

	if result.Secret.Url != "" {
		if err := openBrowser(result.Secret.Url); err != nil {
			output.Error("", fmt.Sprintf("failed to open the browser: %v", err))
			return
		}
	} else {
		output.Warning("", fmt.Sprintf("secret %q does not define a URL to open", result.Secret.Name))
	}

	if len(secretPwd) == 0 {
		output.Warning("", fmt.Sprintf("secret %q does not have a password to copy", result.Secret.Name))
		return
	}

	if err := clipboard.WriteAll(secretPwd); err != nil {
		output.Error("", err.Error())
		return
	}

	fmt.Println()
	utils.Success(output.BoldS("the secret password is in your clipboard"))
}

func openBrowser(url string) error {
	if len(url) == 0 {
		return errors.New("empty URL")
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
