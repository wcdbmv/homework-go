package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	handle(checkArgs())

	if len(os.Args) == 2 {
		result, err := Calculate(os.Args[1])
		handle(err)

		fmt.Println(result)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		handle(err)  // handle only io err this

		result, err := Calculate(line)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkArgs() error {
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run calc.go [expression]")
		return errors.New("len(os.Args) > 2")
	}
	return nil
}