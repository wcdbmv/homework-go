package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	flag.UintVar(&flags.column,          "k", 0,     "sort by k column")
	flag.StringVar(&output      ,        "o", "",    "write result to FILE instead of standard output")
}

func main() {
	flag.Parse()

	handle(checkArgs())

	var (
		reader io.Reader = os.Stdin
		fin    *os.File
		err    error
	)
	if len(flag.Args()) == 1 {
		fin, err = os.Open(flag.Args()[0])
		handle(err)
		defer fin.Close()
		reader = fin
	}

	lines, err := ReadLines(reader)
	handle(err)

	sortedLines, err := Sorted(lines, flags)
	handle(err)

	var (
		writer io.Writer = os.Stdout
		fout   *os.File
	)
	if output != "" {
		fout, err = os.Create(output)
		handle(err)
		defer fout.Close()
		writer = fout
	}

	handle(WriteLines(writer, sortedLines))
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkArgs() error {
	if len(flag.Args()) > 1 {
		flag.Usage()
		return errors.New("wrong usage: there must be no more than one argument")
	}
	return nil
}

func getWriter() (io.Writer, error) {
	if output == "" {
		return os.Stdout, nil
	}
	return os.Create(output)
}

func ReadLines(reader io.Reader) ([]string, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	if lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}
	return lines, nil
}

func WriteLines(writer io.Writer, lines []string) error {
	text := strings.Join(lines, "\n") + "\n"
	_, err := io.WriteString(writer, text)
	return err
}
