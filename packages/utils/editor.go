package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

// GetNoteFromEditor opens a temporary file in the configured editor and returns
// the edited note, if it changed, and any execution error.
func GetNoteFromEditor(initial string) (string, bool, error) {
	editorSpecs := resolveEditorSpecs()
	requireWait := parseBoolEnv("RAPTOR_EDITOR_REQUIRE_WAIT", false)

	tmpFilePath, err := writeTempFile(initial)
	if err != nil {
		return "", false, fmt.Errorf("cannot create temp file: %w", err)
	}
	defer cleanupTempFile(tmpFilePath)

	var lastErr error
	for _, editorSpec := range editorSpecs {
		cmdName, cmdArgs, err := parseCommandLine(editorSpec)
		if err != nil {
			lastErr = err
			continue
		}

		if requireWait {
			cmdName, cmdArgs = ensureWaitFlagIfSupported(cmdName, cmdArgs)
		}

		fullArgs := append(cmdArgs, tmpFilePath)
		cmd := exec.Command(cmdName, fullArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if isEditorNotFound(err) {
				lastErr = err
				continue
			}
			return "", false, fmt.Errorf("editor failed: %w", err)
		}

		editedBytes, err := os.ReadFile(tmpFilePath)
		if err != nil {
			return "", false, fmt.Errorf("cannot read edited note: %w", err)
		}

		initialNorm := normalizeForCompare(initial)
		editedNorm := normalizeForCompare(string(editedBytes))
		changed := editedNorm != initialNorm
		return editedNorm, changed, nil
	}

	if lastErr != nil {
		return "", false, fmt.Errorf("editor failed: %w", lastErr)
	}

	return "", false, fmt.Errorf("editor failed: no editor command configured")
}

func resolveEditorSpecs() []string {
	defaultSpecs := defaultEditorSpecs()

	for _, key := range []string{"RAPTOR_EDITOR", "VISUAL", "EDITOR"} {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			specs := []string{value}
			for _, fallback := range defaultSpecs {
				if strings.EqualFold(value, fallback) {
					continue
				}
				specs = append(specs, fallback)
			}
			return specs
		}
	}

	return defaultSpecs
}

func defaultEditorSpecs() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"notepad"}
	case "darwin":
		return []string{"nano", "vi"}
	default:
		return []string{"nano", "vi"}
	}
}

func parseCommandLine(editorSpec string) (string, []string, error) {
	tokens, err := tokenizeCommandLine(editorSpec)
	if err != nil {
		return "", nil, fmt.Errorf("invalid editor command %q: %w", editorSpec, err)
	}

	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("empty editor command")
	}

	return tokens[0], tokens[1:], nil
}

func tokenizeCommandLine(input string) ([]string, error) {
	var (
		tokens       []string
		current      strings.Builder
		activeQuote  rune
		tokenStarted bool
	)

	flushToken := func() {
		if !tokenStarted {
			return
		}
		tokens = append(tokens, current.String())
		current.Reset()
		tokenStarted = false
	}

	runes := []rune(strings.TrimSpace(input))
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if activeQuote == 0 {
			switch {
			case unicode.IsSpace(r):
				flushToken()
			case r == '\'' || r == '"':
				activeQuote = r
				tokenStarted = true
			case r == '\\':
				tokenStarted = true
				if i+1 < len(runes) {
					next := runes[i+1]
					if unicode.IsSpace(next) || next == '\'' || next == '"' || next == '\\' {
						current.WriteRune(next)
						i++
						continue
					}
				}
				current.WriteRune(r)
			default:
				tokenStarted = true
				current.WriteRune(r)
			}
			continue
		}

		if activeQuote == '\'' {
			if r == '\'' {
				activeQuote = 0
				continue
			}
			current.WriteRune(r)
			continue
		}

		// Double-quoted token.
		if r == '"' {
			activeQuote = 0
			continue
		}
		if r == '\\' {
			if i+1 < len(runes) {
				next := runes[i+1]
				if next == '"' || next == '\\' {
					current.WriteRune(next)
					i++
					continue
				}
			}
			current.WriteRune(r)
			continue
		}
		current.WriteRune(r)
	}

	if activeQuote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}

	flushToken()
	return tokens, nil
}

func ensureWaitFlagIfSupported(name string, args []string) (string, []string) {
	baseName := strings.ToLower(filepath.Base(name))
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	hasArg := func(flags ...string) bool {
		for _, arg := range args {
			for _, flag := range flags {
				if strings.EqualFold(arg, flag) {
					return true
				}
			}
		}
		return false
	}

	switch baseName {
	case "code", "code-insiders", "subl", "sublime_text":
		if !hasArg("--wait", "-w") {
			args = append(args, "--wait")
		}
	case "gedit":
		if !hasArg("--wait") {
			args = append(args, "--wait")
		}
	}

	return name, args
}

func parseBoolEnv(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "1", "t", "true", "y", "yes", "on":
		return true
	case "0", "f", "false", "n", "no", "off":
		return false
	default:
		return defaultValue
	}
}

func writeTempFile(initial string) (string, error) {
	tmpFile, err := os.CreateTemp("", "raptor-note-*.md")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if err := tmpFile.Chmod(0o600); err != nil {
		// Windows may not fully support POSIX-style permissions.
		if runtime.GOOS != "windows" {
			return "", err
		}
	}

	if _, err := tmpFile.WriteString(initial); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func cleanupTempFile(filePath string) {
	_ = os.Remove(filePath)
}

func normalizeForCompare(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.TrimSuffix(s, "\n")
	return s
}

func isEditorNotFound(err error) bool {
	if errors.Is(err, exec.ErrNotFound) {
		return true
	}

	var execErr *exec.Error
	if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
		return true
	}

	var pathErr *os.PathError
	return errors.As(err, &pathErr) && errors.Is(pathErr.Err, fs.ErrNotExist)
}
