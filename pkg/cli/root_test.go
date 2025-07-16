package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "help flag",
			args:     []string{"--help"},
			wantHelp: true,
		},
		{
			name:     "version flag",
			args:     []string{"--version"},
			wantHelp: true,
		},
		{
			name:     "no args shows help",
			args:     []string{},
			wantHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the root command for testing
			cmd := &cobra.Command{
				Use:     rootCmd.Use,
				Short:   rootCmd.Short,
				Long:    rootCmd.Long,
				Version: rootCmd.Version,
			}

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.wantHelp {
				// For help and version, we expect output but no error
				if err != nil && !strings.Contains(err.Error(), "help requested") {
					t.Errorf("Expected help output, got error: %v", err)
				}

				output := buf.String()
				if output == "" {
					t.Error("Expected help output, got empty string")
				}
			}
		})
	}
}

func TestRootCmdProperties(t *testing.T) {
	if rootCmd.Use != "glypto" {
		t.Errorf("Expected Use to be 'glypto', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}

	if rootCmd.Version != "0.1.0" {
		t.Errorf("Expected Version to be '0.1.0', got '%s'", rootCmd.Version)
	}
}
