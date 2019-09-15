package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"math"
	"strings"
	"unicode"
)

const (
	openBracket  byte = 40
	closeBracket byte = 41
	plus         byte = 43
	minus        byte = 45
	multiply     byte = 42
	divide       byte = 47
)

func priority(operator byte) int {
	switch operator {
	case openBracket:
		fallthrough
	case closeBracket:
		return 0
	case plus:
		fallthrough
	case minus:
		return 1
	case multiply:
		fallthrough
	case divide:
		return 2
	}
	log.Println("non-operator in priority func")
	return 0
}

func isOperator(c byte) bool {
	return c == plus || c == minus || c == multiply || c == divide || c == openBracket || c == closeBracket
}

// reads operand from string
// if not success — restores reader
func sreadOperand(reader *strings.Reader) (f float64, err error) {
	offset := reader.Size() - int64(reader.Len())
	_, err = fmt.Fscan(reader, &f)
	if err != nil {
		reader.Seek(offset, io.SeekStart)
	}
	return
}

func sreadOperator(reader *strings.Reader) (c byte, err error) {
	err = skipSpaces(reader)
	if err != nil {
		return
	}
	c, err = reader.ReadByte()
	if err != nil {
		return
	}
	if !isOperator(c) {
		_ = reader.UnreadByte()
		err = errors.New("not an operator")
	}
	return
}

func skipSpaces(reader *strings.Reader) error {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err  // io.EOF
		}
		if !unicode.IsSpace(rune(b)) {
			_ = reader.UnreadByte()
			return nil
		}
	}
}

func shuntingYard(infixExpr string) (postfixExpr string, err error) {
	reader := strings.NewReader(infixExpr)
	var operators Stack

	for checkOperand := true; reader.Len() != 0; {

		if checkOperand {
			// maybe it's an operand ?
			operand, err := sreadOperand(reader)
			if err == nil {
				postfixExpr = fmt.Sprintf("%s%f ", postfixExpr, operand)
				checkOperand = false
				continue
			}
		}

		// okay, it must be an operator
		operator, err := sreadOperator(reader)
		if err != nil {
			return "", errors.New("invalid format")
		}

		// push operator
		checkOperand = true
		if operator == openBracket {
			operators.Push(operator)
			continue
		}
		if operator == closeBracket {
			checkOperand = false
			for {
				if operators.Empty() {
					return "", errors.New("invalid format")
				}
				operator = operators.Pop().(byte)
				if operator == openBracket {
					break
				}
				postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operator))
			}
			continue
		}
		if !operators.Empty() && priority(operator) <= priority(operators.Top().(byte)) {
			postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operators.Pop().(byte)))
		}
		operators.Push(operator)
	}

	for !operators.Empty() {
		operator := operators.Pop().(byte)
		if operator == openBracket {
			return "", errors.New("invalid format")
		}
		postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operator))
	}

	postfixExpr = strings.TrimSpace(postfixExpr)

	return
}

func Calculate(infixExpr string) (result float64, err error) {
	postfixExpr, err := shuntingYard(infixExpr)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert from infix to postfix")
	}

	reader := strings.NewReader(postfixExpr)
	var operands Stack

	for reader.Len() != 0 {
		operand, err := sreadOperand(reader)
		if err == nil {
			operands.Push(operand)
			continue
		}

		// it must be an operator
		operator, _ := sreadOperator(reader)
		err = apply(&operands, operator)
		if err != nil {
			return 0, err
		}
	}

	if operands.Empty() {
		return 0, errors.New("empty expression")
	}
	result = operands.Pop().(float64)
	if !operands.Empty() {
		return 0, errors.New("invalid format")
	}
	return
}

func apply(operands *Stack, operator byte) error {
	if operands.Empty() {
		return errors.New("invalid format")
	}
	b := operands.Pop().(float64)
	if operands.Empty() {
		return errors.New("invalid format")
	}
	a := operands.Pop().(float64)
	c, err := execute(a, b, operator)
	if err != nil {
		return err
	}
	operands.Push(c)
	return nil
}

func execute(a float64, b float64, operator byte) (float64, error) {
	switch operator {
	case plus:
		return a + b, nil
	case minus:
		return a - b, nil
	case multiply:
		return a * b, nil
	case divide:
		if math.Abs(b) < 1e-9 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	return 0, errors.New("invalid operand")  // never
}
