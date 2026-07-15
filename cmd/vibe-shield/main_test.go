package main

import (
	"reflect"
	"testing"
)

func TestStripNoColorFlag(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		wantArgs  []string
		wantColor bool
	}{
		{
			name:     "no flag",
			input:    []string{"npm", "run"},
			wantArgs: []string{"npm", "run"},
		},
		{
			name:      "leading flag",
			input:     []string{"--no-color", "npm", "run"},
			wantArgs:  []string{"npm", "run"},
			wantColor: true,
		},
		{
			name:     "flag only in child args",
			input:    []string{"npm", "run", "--no-color"},
			wantArgs: []string{"npm", "run", "--no-color"},
		},
		{
			name:      "duplicate leading flags",
			input:     []string{"--no-color", "--no-color", "npm"},
			wantArgs:  []string{"npm"},
			wantColor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotColor := stripNoColorFlag(tt.input)
			if gotColor != tt.wantColor {
				t.Errorf("noColor = %v, want %v", gotColor, tt.wantColor)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
