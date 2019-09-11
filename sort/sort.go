package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var (
	caseInsensitive bool
	unique          bool
	reverse         bool
	numeric         bool
	column          int
	output          string
)

func init() {
	flag.BoolVar(&caseInsensitive, "f", false, "fold lower case to upper case characters")
	flag.BoolVar(&unique,          "u", false, "output only the first of an equal run")
	flag.BoolVar(&reverse,         "r", false, "reverse the result of comparisons")
	flag.BoolVar(&numeric,         "n", false, "compare according to string numerical value")
	flag.IntVar(&column,           "k", 0,     "sort by k column")
	flag.StringVar(&output,        "o", "",    "write result to FILE instead of standard output")
}

type StringMap struct {
	key   string
	value *string
}

func main() {
	flag.Parse()

	handle(checkArgs())

	text, err := ReadLines(flag.Args()[0])
	handle(err)

	lines := link(text)

	if column > 0 {
		handle(extractColumn(lines, column))
	}

	if caseInsensitive {
		lowerize(lines)
	}

	if unique {
		if numeric {
			lines, err = removeDuplicatesNumeric(lines)
			handle(err)
		} else {
			lines = removeDuplicates(lines)
		}
	}

	less := func(i, j int) bool { return lines[i].key < lines[j].key }
	if numeric {
		less = func(i, j int) bool { return atof(lines[i].key) < atof(lines[j].key) }
	}

	sort.Slice(lines, less)

	if reverse {
		for i, j := 0, len(lines) - 1; i < j; i, j = i + 1, j - 1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	err = Write(output, join(lines))
	handle(err)
}

func checkArgs() error {
	if len(flag.Args()) != 1 {
		flag.Usage()
		return errors.New("flag.Args() != 1")
	}
	if column < 0 {
		return errors.New("column cannot be non-positive")
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

func Write(output string, data string) error {
	if output != "" {
		return ioutil.WriteFile(output, []byte(data), 0644)
	}
	_, err := fmt.Print(data)
	return err
}

func link(text []string) []StringMap {
	lines := make([]StringMap, len(text))
	for i, line := range text {
		lines[i] = StringMap{line, &text[i]}
	}
	return lines
}

func extractColumn(lines []StringMap, k int) error {
	for _, line := range lines {
		strSlice := strings.Fields(line.key)
		if len(strSlice) < k {
			return errors.New(fmt.Sprint("Not enough column in row %i", k))
		}
		line.key = strSlice[k - 1]
	}
	return nil
}

func lowerize(lines []StringMap) {
	for i, line := range lines {
		lines[i].key = strings.ToLower(line.key)
	}
}

func removeDuplicates(lines []StringMap) []StringMap {
	used := make(map[string]bool)
	var unique []StringMap
	for _, line := range lines {
		if _, value := used[line.key]; !value {
			used[line.key] = true
			unique = append(unique, line)
		}
	}
	return unique
}

func removeDuplicatesNumeric(lines []StringMap) ([]StringMap, error) {
	used := make(map[float64]bool)
	var unique []StringMap
	for _, line := range lines {
		f, err := strconv.ParseFloat(line.key, 64)
		if err != nil {
			return nil, err
		}
		if _, value := used[f]; !value {
			used[f] = true
			unique = append(unique, line)
		}
	}
	return unique, nil
}

func join(lines []StringMap) string {
	if len(lines) == 0 {
		return "\n"
	}
	var text string
	for _, line := range lines {
		text = text + *line.value + "\n"
	}
	return text
}

func atof(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
