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
	fmt.Println(text)

	vm.Run(text)
}
