package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

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

var linesOut1_noflags = []string{
	"Apple",
	"BOOK",
	"Book",
	"Go",
	"Hauptbahnhof",
	"January",
	"January",
	"Napkin",
}

var linesOut1_r = []string{
	"Napkin",
	"January",
	"January",
	"Hauptbahnhof",
	"Go",
	"Book",
	"BOOK",
	"Apple",
}

var linesOut1_u = []string{
	"Apple",
	"BOOK",
	"Book",
	"Go",
	"Hauptbahnhof",
	"January",
	"Napkin",
}

var linesOut1_uf = []string{
	"Apple",
	"BOOK",
	"Go",
	"Hauptbahnhof",
	"January",
	"Napkin",
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

var linesOut2_k1ufr = []string{
	"Napkin         1000",
	"January        1.2e3",
	"Hauptbahnhof   -0.123",
	"Go             0.01",
	"BOOK           -123e-3",
	"Apple          123.4",
}

var linesOut2_k2ufnr = []string{
	"January        1.2e3",
	"Napkin         1000",
	"Apple          123.4",
	"Go             0.01",
	"January        0.0e0",
	"Book           -0.01",
	"BOOK           -123e-3",
}

func TestUsual(t *testing.T) {
	in1 := linesIn1

	out1_noflags, err := Sorted(in1, Flags{})
	require.Equal(t, linesOut1_noflags, out1_noflags, errMsg("nop", err))

	out1_r, err := Sorted(in1, Flags{reverse:true})
	require.Equal(t, linesOut1_r, out1_r, errMsg("-r", err))

	out1_u, err := Sorted(in1, Flags{unique:true})
	require.Equal(t, linesOut1_u, out1_u, errMsg("-u", err))

	out1_uf, err := Sorted(in1, Flags{caseInsensitive:true, unique:true})
	require.Equal(t, linesOut1_uf, out1_uf, errMsg("-u -f", err))
}

func TestNonTrivial(t *testing.T) {
	in2 := linesIn2

	out2_k1ufr, err := Sorted(in2, Flags{column:1, unique:true, caseInsensitive:true, reverse:true})
	require.Equal(t, linesOut2_k1ufr, out2_k1ufr, errMsg("-k=1 -u -f -r", err))

	out2_k2ufnr, err := Sorted(in2, Flags{
		caseInsensitive: true,
		unique:          true,
		reverse:         true,
		numeric:         true,
		column:          2,
	})
	require.Equal(t, linesOut2_k2ufnr, out2_k2ufnr, errMsg("-k=2 -u -f -n -r", err))
}

func TestFail(t *testing.T) {
	in2 := linesIn2

	_, err := Sorted(in2, Flags{column:3})
	require.Error(t, err, errMsg("-k=3", err))

	_, err = Sorted(in2, Flags{column:1, numeric:true, unique:true})
	require.Error(t, err, errMsg("-k=1 -u -n", err))
}

func errMsg(flags string, err error) string {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}
	return fmt.Sprintf("flags: %s [err is `%s`]", flags, errStr)
}
