/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/mas2020-golang/cryptex/cmd"
	"github.com/mas2020-golang/cryptex/packages/utils"
)

var GitCommit, BuildDate string

func main() {
	// Retrieve the log level from the environment variable
	logLevel := getLogLevelFromEnv()

	// Create a new logger with the specified log level
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	// set the default logger for the whole application
	slog.SetDefault(logger)
	utils.GitCommit = GitCommit
	utils.BuildDate = BuildDate
	cmd.Execute()
}

// getLogLevelFromEnv retrieves the log level from the RAPTOR_LOGLEVEL environment variable.
// Defaults to Info level if the variable is not set or contains an invalid value.
func getLogLevelFromEnv() slog.Level {
	levelStr := strings.ToLower(os.Getenv("RAPTOR_LOGLEVEL"))
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
