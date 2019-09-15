package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSorted(t *testing.T) {
	var linesIn1 = []string{
		"Napkin",
		"Apple",
		"January",
		"BOOK",
		"January",
		"Hauptbahnhof",
		"Book",
		"Go",
	}
	var linesIn2 = []string{
		"Napkin         1000",
		"Apple          123.4",
		"January        1.2e3",
		"BOOK           -123e-3",
		"January        0.0e0",
		"Hauptbahnhof   -0.123",
		"Book           -0.01",
		"Go             0.01",
	}
	type args struct {
		lines []string
		flags Flags
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Usual noflags",
			args{
				linesIn1,
				Flags{},
			},
			[]string{"Apple", "BOOK", "Book", "Go", "Hauptbahnhof", "January", "January", "Napkin"},
			false,
		}, {
			"Usual -r",
		args{
			linesIn1,
			Flags{reverse:true},
		},
		[]string{"Napkin", "January", "January", "Hauptbahnhof", "Go", "Book", "BOOK", "Apple"},
		false,
		}, {
			"Usual -u",
			args{
				linesIn1,
				Flags{unique:true},
			},
		[]string{"Apple", "BOOK", "Book", "Go", "Hauptbahnhof", "January", "Napkin"},
		false,
		}, {
			"Usual -u -f",
			args{
				linesIn1,
				Flags{unique:true, caseInsensitive:true},
			},
			[]string{"Apple", "BOOK", "Go", "Hauptbahnhof", "January", "Napkin"},
			false,
		}, {
			"NonTrivial -k=1 -u -f -r",
			args{
				linesIn2,
				Flags{column:1, unique:true, caseInsensitive:true, reverse:true},
			},
			[]string{
				"Napkin         1000",
				"January        1.2e3",
				"Hauptbahnhof   -0.123",
				"Go             0.01",
				"BOOK           -123e-3",
				"Apple          123.4",
			},
			false,
		}, {
			"NonTrivial -k=2 -u -f -n -r",
			args{
				linesIn2,
				Flags{column:2, unique:true, caseInsensitive:true, numeric:true, reverse:true},
			},
			[]string{
				"January        1.2e3",
				"Napkin         1000",
				"Apple          123.4",
				"Go             0.01",
				"January        0.0e0",
				"Book           -0.01",
				"BOOK           -123e-3",
			},
			false,
		}, {
			"Fail -k=3",
			args{
				linesIn2,
				Flags{column:3},
			},
			nil,
			true,
		}, {
			"Fail -k=1 -u -n",
			args{
				linesIn2,
				Flags{column:1, unique:true, numeric:true},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sorted(tt.args.lines, tt.args.flags)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
