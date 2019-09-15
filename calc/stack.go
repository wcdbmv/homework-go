package main

type StackNode struct {
	data interface{}
	next *StackNode
}

type Stack struct {
	top *StackNode
}

func (stack *Stack) Push(data interface{}) {
	stack.top = &StackNode{data, stack.top}
}

func (stack *Stack) Pop() (data interface{}) {
	if !stack.Empty() {
		data, stack.top = stack.top.data, stack.top.next
		return data
	}
	return nil
}

func (stack *Stack) Top() interface{} {
	if !stack.Empty() {
		return stack.top.data
	}
	return nil
}

func (stack *Stack) Empty() bool {
	return stack.top == nil
}
