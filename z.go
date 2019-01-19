package main

import (
	"fmt"
)

func main() {
	str := []rune("yeet")
	str = append(str[:3], str[3+1:]...)
	str = append(str[:2], str[2+1:]...)
	fmt.Println(string(str))
	fmt.Printf("%d\n", len(str))
}
