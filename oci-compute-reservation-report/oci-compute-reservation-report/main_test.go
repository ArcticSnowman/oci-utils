package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantDebug bool
		wantError bool
		wantHelp  bool
	}{
		{name: "default", args: []string{"--compartment", "engineering"}},
		{name: "debug", args: []string{"--compartment", "engineering", "--debug"}, wantDebug: true},
		{name: "debug shorthand", args: []string{"--compartment", "engineering", "-d"}, wantDebug: true},
		{name: "debug disabled", args: []string{"--compartment", "engineering", "--debug=false"}},
		{name: "help", args: []string{"--help"}, wantHelp: true},
		{name: "missing compartment", wantError: true},
		{name: "unexpected argument", args: []string{"unexpected"}, wantError: true},
		{name: "unknown flag", args: []string{"--unknown"}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantError {
				t.Fatalf("Execute() error = %v, wantError = %v", err, tt.wantError)
			}
			if got := strings.Contains(stderr.String(), "level=debug"); got != tt.wantDebug {
				t.Errorf("debug output = %v, want %v; stderr: %s", got, tt.wantDebug, stderr.String())
			}
			if tt.wantError {
				return
			}
			if tt.wantHelp {
				if !strings.Contains(stdout.String(), "-c, --compartment") {
					t.Errorf("help does not document the compartment flag: %s", stdout.String())
				}
				if !strings.Contains(stdout.String(), "--profile") {
					t.Errorf("help does not document the profile flag: %s", stdout.String())
				}
				return
			}
			if !strings.Contains(stdout.String(), "not implemented yet") {
				t.Errorf("missing starter message: %s", stdout.String())
			}
		})
	}
}
