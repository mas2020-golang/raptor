/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mas2020-golang/cryptex/cmd/list"
	"github.com/mas2020-golang/cryptex/packages/utils"
	"github.com/spf13/cobra"
)

// Define styles
var (
	// Main container style
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			PaddingTop(0).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2)
		//Margin(1)

	// Header style
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	// Section header style
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Underline(true).
			MarginTop(1).
			MarginBottom(1)

	// Key style (left column) - ensuring consistent alignment
	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(false).
			Width(19). // Increased width to accommodate longest key
			Align(lipgloss.Left)

	// Value style (right column)
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			MarginLeft(2)

	// Missing value style
	missingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Italic(true).
			MarginLeft(2)

	// Version specific styles
	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("118")).
			Bold(true).
			MarginLeft(2)

	commitStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			MarginLeft(2)
)

// boxCmd represents the box command
// boxCmd represents the box command
func newInfoCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "info",
		Short: "Print useful information about raptor",
		Long: `The info command prints the following information:
    - env variable values
    - box list
    - version and commit info
`,
		Run: func(cmd *cobra.Command, args []string) {
			DisplayEnvironmentInfo()
		},
	}

	return c
}

type EnvVar struct {
	Key         string
	Description string
	Required    bool
}

func DisplayEnvironmentInfo() {
	// Define your environment variables
	envVars := []EnvVar{
		{"CRYPTEX_FOLDER", "Cryptex folder path", false},
		{"CRYPTEX_BOX", "Cryptex box configuration", false},
		{"RAPTOR_LOGLEVEL", "Logging level for Raptor", false},
		{"RAPTOR_TIMEOUT_SEC", "Timeout in seconds for Raptor", false},
	}

	var content strings.Builder

	// Header
	// header := headerStyle.Render("🚀 APPLICATION CONFIGURATION")
	// content.WriteString(header + "\n")

	// Application Info Section
	appSection := sectionStyle.Render("📋 Application Info")
	content.WriteString(appSection + "\n")

	// Version
	versionLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyStyle.Render("Version:"),
		versionStyle.Render(utils.Version),
	)
	content.WriteString(versionLine + "\n")

	// Commit
	commitLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyStyle.Render("Commit:"),
		commitStyle.Render(utils.GitCommit),
	)
	content.WriteString(commitLine + "\n")

	// Environment Variables Section
	envSection := sectionStyle.Render("🌍 Environment Variables")
	content.WriteString(envSection + "\n")

	// Display each environment variable
	for _, envVar := range envVars {
		value := os.Getenv(envVar.Key)

		var valuePart string
		if value != "" {
			// Special formatting for specific variables
			switch envVar.Key {
			case "RAPTOR_TIMEOUT_SEC":
				if timeout, err := strconv.Atoi(value); err == nil {
					valuePart = valueStyle.Render(fmt.Sprintf("%d seconds", timeout))
				} else {
					valuePart = valueStyle.Render(value)
				}
			case "RAPTOR_LOGLEVEL":
				// Add emoji based on log level
				emoji := getLogLevelEmoji(value)
				valuePart = valueStyle.Render(fmt.Sprintf("%s %s", emoji, value))
			default:
				valuePart = valueStyle.Render(value)
			}
		} else {
			if envVar.Required {
				valuePart = missingStyle.Render("❌ REQUIRED - NOT SET")
			} else {
				valuePart = missingStyle.Render("Not set")
			}
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			keyStyle.Render(envVar.Key+":"),
			valuePart,
		)
		content.WriteString(line + "\n")
	}

	// Add configuration
	boxSection := sectionStyle.Render("Boxes")
	content.WriteString(boxSection + "\n")

	folderBox, boxes, err := getBoxes()
	renderedBoxes := ""
	if err != nil {
		renderedBoxes = err.Error()
	} else {
		for _, b := range boxes {
			renderedBoxes += fmt.Sprintf("- %s\n", b.Name)
		}
	}
	boxesLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyStyle.Render(renderedBoxes),
	)
	// content.WriteString(statusLine + "\n")
	content.WriteString(boxesLine)
	content.WriteString("\nsearched in " + folderBox)
	// Wrap everything in the container
	final := containerStyle.Render(content.String())
	fmt.Println(final)
}

// Helper function to get emoji for log levels
func getLogLevelEmoji(level string) string {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return "🐛"
	case "INFO":
		return "ℹ️"
	case "WARN", "WARNING":
		return "⚠️"
	case "ERROR":
		return "❌"
	case "FATAL":
		return "💀"
	default:
		return "📝"
	}
}

// Helper function to determine configuration status
func getBoxes() (string, []utils.Box, error) {
	return list.ListBoxes("")
}

// Alternative compact version
func DisplayCompactInfo() {
	// Create a table-like layout
	rows := [][]string{
		{"Version", utils.Version},
		{"Commit", utils.GitCommit},
		{"CRYPTEX_FOLDER", getEnvOrDefault("CRYPTEX_FOLDER", "Not set")},
		{"CRYPTEX_BOX", getEnvOrDefault("CRYPTEX_BOX", "Not set")},
		{"RAPTOR_LOGLEVEL", getEnvOrDefault("RAPTOR_LOGLEVEL", "Not set")},
		{"RAPTOR_TIMEOUT_SEC", getEnvOrDefault("RAPTOR_TIMEOUT_SEC", "Not set")},
	}

	var tableContent strings.Builder
	header := headerStyle.Render("🚀 Configuration Summary")
	tableContent.WriteString(header + "\n\n")

	for _, row := range rows {
		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			keyStyle.Render(row[0]+":"),
			valueStyle.Render(row[1]),
		)
		tableContent.WriteString(line + "\n")
	}

	final := containerStyle.Render(tableContent.String())
	fmt.Println(final)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// You can choose which version to use
	fmt.Println("Detailed Version:")
	DisplayEnvironmentInfo()

	fmt.Println("\nCompact Version:")
	DisplayCompactInfo()
}
