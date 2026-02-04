package utils

import (
	"reflect"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    string
		cmd     string
		args    []string
		wantErr bool
	}{
		{
			name: "single binary",
			spec: "nano",
			cmd:  "nano",
			args: []string{},
		},
		{
			name: "binary with args",
			spec: "code --wait --reuse-window",
			cmd:  "code",
			args: []string{"--wait", "--reuse-window"},
		},
		{
			name: "quoted binary path",
			spec: "\"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code\" --wait",
			cmd:  "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			args: []string{"--wait"},
		},
		{
			name:    "unterminated quote",
			spec:    "\"code --wait",
			wantErr: true,
		},
		{
			name:    "empty spec",
			spec:    "  ",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd, args, err := parseCommandLine(tc.spec)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.spec)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cmd != tc.cmd {
				t.Fatalf("command mismatch: got %q want %q", cmd, tc.cmd)
			}
			if !reflect.DeepEqual(args, tc.args) {
				t.Fatalf("args mismatch: got %#v want %#v", args, tc.args)
			}
		})
	}
}

func TestEnsureWaitFlagIfSupported(t *testing.T) {
	t.Parallel()

	name, args := ensureWaitFlagIfSupported("code", []string{"--reuse-window"})
	if name != "code" {
		t.Fatalf("unexpected command name: %q", name)
	}
	if !reflect.DeepEqual(args, []string{"--reuse-window", "--wait"}) {
		t.Fatalf("unexpected args: %#v", args)
	}

	_, args = ensureWaitFlagIfSupported("nano", []string{})
	if len(args) != 0 {
		t.Fatalf("nano should not receive wait flags, got %#v", args)
	}
}

func TestNormalizeForCompare(t *testing.T) {
	t.Parallel()

	input := "line1\r\nline2\r\n"
	got := normalizeForCompare(input)
	want := "line1\nline2"
	if got != want {
		t.Fatalf("normalize mismatch: got %q want %q", got, want)
	}
}

func TestParseBoolEnv(t *testing.T) {
	t.Setenv("RAPTOR_EDITOR_REQUIRE_WAIT", "true")
	if !parseBoolEnv("RAPTOR_EDITOR_REQUIRE_WAIT", false) {
		t.Fatalf("expected true")
	}

	t.Setenv("RAPTOR_EDITOR_REQUIRE_WAIT", "invalid")
	if !parseBoolEnv("RAPTOR_EDITOR_REQUIRE_WAIT", true) {
		t.Fatalf("expected fallback default true")
	}

	t.Setenv("RAPTOR_EDITOR_REQUIRE_WAIT", "0")
	if parseBoolEnv("RAPTOR_EDITOR_REQUIRE_WAIT", true) {
		t.Fatalf("expected false")
	}
}

func TestResolveEditorSpecs(t *testing.T) {
	t.Setenv("RAPTOR_EDITOR", "code --wait")
	t.Setenv("VISUAL", "vim")
	t.Setenv("EDITOR", "nano")

	specs := resolveEditorSpecs()
	if len(specs) == 0 {
		t.Fatalf("expected at least one editor spec")
	}
	if specs[0] != "code --wait" {
		t.Fatalf("expected RAPTOR_EDITOR first, got %q", specs[0])
	}
}
