package main

import "testing"

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
	compareTestResult(out1_noflags, linesOut1_noflags, err, t, "Usual noflags")

	out1_r, err := Sorted(in1, Flags{reverse:true})
	compareTestResult(out1_r, linesOut1_r, err, t, "Usual -r")

	out1_u, err := Sorted(in1, Flags{unique:true})
	compareTestResult(out1_u, linesOut1_u, err, t, "Usual -u")

	out1_uf, err := Sorted(in1, Flags{caseInsensitive:true, unique:true})
	compareTestResult(out1_uf, linesOut1_uf, err, t, "Usual -u -f")
}

func TestNonTrivial(t *testing.T) {
	in2 := linesIn2

	out2_k1ufr, err := Sorted(in2, Flags{column:1, unique:true, caseInsensitive:true, reverse:true})
	compareTestResult(out2_k1ufr, linesOut2_k1ufr, err, t, "NonTrivial -k=1 -u -f -r")

	out2_k2ufnr, err := Sorted(in2, Flags{
		caseInsensitive: true,
		unique:          true,
		reverse:         true,
		numeric:         true,
		column:          2,
	})
	compareTestResult(out2_k2ufnr, linesOut2_k2ufnr, err, t, "NonTrivial -k=2 -u -f -n -r")
}

func TestFail(t *testing.T) {
	in2 := linesIn2

	_, err := Sorted(in2, Flags{column:3})
	checkFail(err, t, "Fail -k=3")

	_, err = Sorted(in2, Flags{column:1, numeric:true, unique:true})
	checkFail(err, t, "Fail -k=1 -u -n")
}

func equal(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func compareTestResult(got, expected []string, err error, t *testing.T, testName string) {
	if err != nil {
		t.Errorf("Test%s failed: %s", testName, err)
	}
	if !equal(got, expected) {
		t.Errorf("Test%s failed, result not match", testName)
	}
}

func checkFail(err error, t *testing.T, testName string) {
	if err == nil {
		t.Errorf("Test%s failed: err is nil", testName)
	}
}
