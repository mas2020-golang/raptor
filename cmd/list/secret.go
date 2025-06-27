/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

var (
	items           bool
	boxName, filter string
)

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	headerStyle  = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
	messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lightGray)
)

// boxCmd represents the box command
var ListSecretCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret", "sr"},
	Short:   "List secret",
	Long: `List all the secret in the --box given flag. Use the flag --name
to filter using a regular expression.`,
	Example: `$ raptor secret ls --box test
$ raptor secret ls --box test --name '^secret.*test$`,
	Run: func(cmd *cobra.Command, args []string) {
		listSecrets(cmd)
	},
}

func init() {
	ListSecretCmd.Flags().StringVarP(&filter, "filter", "f", "", "The secret name as a regexp (e.g. 'test.*')")
	ListSecretCmd.Flags().BoolVarP(&items, "items", "i", false, "Show the items' keys for the items saved into the secret")
	ListSecretCmd.PersistentFlags().StringVarP(&boxName, "box", "b", "", "The name of the box where to add the secret")
}

func listSecrets(cmd *cobra.Command) {
	// load the local timezone
	// loc, err := time.LoadLocation("Local")
	// utils.Check(err, "")
	// output variables
	name, version, url, login := "", "", "", ""
	boxPath, _, box, err := utils.OpenBox(boxName, "")
	utils.Check(err, "")
	// get the max length for the NAME, LOGIN attribute
	maxName := getMaxNameLenght(box)
	maxLogin := getMaxLoginLenght(box)

	if len(box.Secrets) == 0 {
		fmt.Println(messageStyle.Render("No secrets yet..."))
		return
	}
	// table format
	t := table.New().
		Headers("NAME", "VERSION", "LOGIN", "URL", "ITEMS", "LAST-UPD").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	for _, s := range box.Secrets {
		loginFormatS := fmt.Sprintf("%%-%ds", maxLogin+2)
		login = fmt.Sprintf(loginFormatS, s.Login)
		if len(s.Version) > 9 {
			version = s.Version[0:6] + "..."
		}
		url = s.Url
		if len(s.Url) > 44 {
			url = s.Url[0:42] + "..."
		}
		url = output.BlueS(fmt.Sprintf("%-47s", url))
		nameFormatS := fmt.Sprintf("%%-%ds", maxName+2)
		name = output.RedS(output.BoldS(fmt.Sprintf(nameFormatS, s.Name)))
		lastUpdated := s.LastUpdated
		// check the name flag
		if len(filter) > 0 {
			r, _ := regexp.Compile("(?i)" + filter) // case insensitive regexp
			if err != nil {
				output.Error("", fmt.Sprintf("Invalid filter: %v\n", err))
				return
			}
			if r.MatchString(name) {
				t.Row(name, version, login, url, strconv.Itoa(len(s.Others)), lastUpdated)
				showItems(s, t)
			}
		} else {
			t.Row(name, version, login, url, strconv.Itoa(len(s.Others)), lastUpdated)
			showItems(s, t)
		}
	}
	v, _ := (*cmd).Parent().Flags().GetBool("verbose")
	if v {
		fmt.Println()
		output.InfoBox(fmt.Sprintf("secret read from the %s box\n", output.BlueS(boxPath)))
	}

	fmt.Println(t.Render())
}

func showItems(s *utils.Secret, t *table.Table) {
	if !items {
		return
	}
	if len(s.Others) > 0 {
		for k := range s.Others {
			t.Row("", "", fmt.Sprintf(" .%s", output.BoldS(k)), "", "", "")
		}
	}
}

// getMaxNameLenght return the max lenght for the NAME attribute
func getMaxNameLenght(box *utils.Box) int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}
	return max
}

// getMaxLoginLenght return the max lenght for the LOGIN attribute
func getMaxLoginLenght(box *utils.Box) int {
	max := 10
	for _, s := range box.Secrets {
		if len(s.Login) > max {
			max = len(s.Login)
		}
	}
	return max
}
