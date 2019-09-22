package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		wantResult float64
		wantErr    bool
	}{
		{
			"It works",
			"2 + 2",
			4,
			false,
		}, {
			"It knows priority",
			"2 + 2 * 2",
			6,
			false,
		}, {
			"It knows parentheses",
			"((((2))+((2)) *((2) + (2))))",
			10,
			false,
		}, {
			"It parses floats good",
			"+1.5+-15e-1",
			0,
			false,
		}, {
			"It w0rKs",
			"(1)",
			1,
			false,
		}, {
			"It can catch division by zero",
			"1 / (1 - 1)",
			0,
			true,
		}, {
			"It wouldn't work with empty expression",
			"",
			0,
			true,
		}, {
			"It wouldn't work with whatever empty expression",
			"              ",
			0,
			true,
		}, {
			"It can detect wrong parenthesis sequance",
			"(((1) + 3)",
			0,
			true,
		}, {
			"It can detect wrong parenthesis sequance 2",
			"((1) + 3))",
			0,
			true,
		}, {
			"think before using it",
			"im dummy dumb blockhead",
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Calculate(tt.args)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.wantResult, gotResult)
		})
	}
}