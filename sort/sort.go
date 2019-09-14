package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Flags struct {
	caseInsensitive bool
	unique          bool
	reverse         bool
	numeric         bool
	column          uint
}

type StringMap struct {
	key   string
	value *string
}

func Sorted(lines []string, flags Flags) ([]string, error) {
	mappedLines := link(lines)

	if flags.column > 0 {
		if err := extractColumn(mappedLines, flags.column); err != nil {
			return nil, err
		}
	}

	if flags.caseInsensitive {
		lowerize(mappedLines)
	}

	if flags.unique {
		uniq, err := uniqueize(mappedLines, flags.numeric)
		if err != nil {
			return nil, err
		}
		mappedLines = uniq
	}

	less := func(i, j int) bool { return mappedLines[i].key < mappedLines[j].key }
	if flags.numeric {
		less = func(i, j int) bool { return atof(mappedLines[i].key) < atof(mappedLines[j].key) }
	}
	sort.Slice(mappedLines, less)

	if flags.reverse {
		for i, j := 0, len(mappedLines) - 1; i < j; i, j = i + 1, j - 1 {
			mappedLines[i], mappedLines[j] = mappedLines[j], mappedLines[i]
		}
	}

	return load(mappedLines), nil
}

func link(lines []string) []StringMap {
	mappedLines := make([]StringMap, len(lines))
	for i := range lines {
		mappedLines[i] = StringMap{lines[i], &lines[i]}
	}
	return mappedLines
}

func extractColumn(mappedLines []StringMap, k uint) error {
	for i := range mappedLines {
		words := strings.Fields(mappedLines[i].key)
		if len(words) < int(k) {
			return errors.New(fmt.Sprint("Not enough columns in row ", i))
		}
		mappedLines[i].key = words[k - 1]
	}
	return nil
}

func lowerize(mappedLines []StringMap) {
	for i := range mappedLines {
		mappedLines[i].key = strings.ToLower(mappedLines[i].key)
	}
}

func uniqueize(lines []StringMap, numeric bool) ([]StringMap, error) {
	removeDuplicatesFunc := removeDuplicates
	if numeric {
		removeDuplicatesFunc = removeDuplicatesNumeric
	}

	return removeDuplicatesFunc(lines)
}

func removeDuplicates(mappedLines []StringMap) ([]StringMap, error) {
	used := make(map[string]bool)
	var unique []StringMap
	for _, line := range mappedLines {
		if _, value := used[line.key]; !value {
			used[line.key] = true
			unique = append(unique, line)
		}
	}
	return unique, nil
}

func removeDuplicatesNumeric(mappedLines []StringMap) ([]StringMap, error) {
	used := make(map[float64]bool)
	var unique []StringMap
	for _, line := range mappedLines {
		f, err := strconv.ParseFloat(line.key, 64)
		if err != nil {
			return nil, err
		}
		if _, was := used[f]; !was {
			used[f] = true
			unique = append(unique, line)
		}
	}
	return unique, nil
}

func load(mappedLines []StringMap) []string {
	lines := make([]string, len(mappedLines))
	for i := range mappedLines {
		lines[i] = *mappedLines[i].value
	}
	return lines
}

func atof(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
