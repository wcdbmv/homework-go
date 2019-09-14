package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	flags  Flags
	output string
)

func init() {
	flag.BoolVar(&flags.caseInsensitive, "f", false, "fold lower case to upper case characters")
	flag.BoolVar(&flags.unique,          "u", false, "output only the first of an equal run")
	flag.BoolVar(&flags.reverse,         "r", false, "reverse the result of comparisons")
	flag.BoolVar(&flags.numeric,         "n", false, "compare according to string numerical value")
	flag.UintVar(&flags.column,           "k", 0,     "sort by k column")
	flag.StringVar(&output      ,        "o", "",    "write result to FILE instead of standard output")
}

func main() {
	flag.Parse()

	handle(checkArgs())

	lines, err := ReadLines(flag.Args()[0])
	handle(err)

	sortedLines, err := Sorted(lines, flags)
	handle(err)

	handle(WriteLines(output, sortedLines))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func checkArgs() error {
	if len(flag.Args()) != 1 {
		flag.Usage()
		return errors.New("flag.Args() != 1")
	}
	return nil
}

func ReadLines(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	text := strings.Split(string(content), "\n")
	if text[len(text) - 1] == "" {
		text = text[:len(text) - 1]
	}
	return text, nil
}

func WriteLines(output string, text []string) error {
	lines := strings.Join(text, "\n") + "\n"
	if output != "" {
		return ioutil.WriteFile(output, []byte(lines), 0644)
	}
	_, err := fmt.Print(lines)
	return err
}
