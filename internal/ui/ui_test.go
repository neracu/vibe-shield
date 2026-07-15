package ui

import (
	"testing"

	"github.com/fatih/color"
)

func TestConfigureColor(t *testing.T) {
	orig := color.NoColor
	t.Cleanup(func() { color.NoColor = orig })

	tests := []struct {
		name        string
		noColorFlag bool
		env         map[string]string
		wantNoColor bool
	}{
		{
			name: "default",
			env: map[string]string{
				"NO_COLOR": "",
				"TERM":     "xterm-256color",
			},
			wantNoColor: false,
		},
		{
			name:        "flag",
			noColorFlag: true,
			wantNoColor: true,
		},
		{
			name:        "NO_COLOR set",
			env:         map[string]string{"NO_COLOR": "1"},
			wantNoColor: true,
		},
		{
			name:        "TERM dumb",
			env:         map[string]string{"TERM": "dumb"},
			wantNoColor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color.NoColor = false
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			ConfigureColor(tt.noColorFlag)
			if color.NoColor != tt.wantNoColor {
				t.Errorf("color.NoColor = %v, want %v", color.NoColor, tt.wantNoColor)
			}
		})
	}
}
