package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')

	value, err := vm.Run(text)
	if !value.IsUndefined() {
		fmt.Println(value)
	}
	if err != nil {
		fmt.Println(err)
	}
}
