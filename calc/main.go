package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	handle(checkArgs())

	result, err := Calculate(os.Args[1])
	handle(err)

	fmt.Println(result)
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkArgs() error {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run calc.go <expression>")
		return errors.New("len(os.Args) != 2")
	}
	return nil
}