package main

import (
	"fmt"
	"strings"
)

func main() {

	str1 := "hello"
	str2 := "HELLO"
	if strings.EqualFold(str1,str2) {
        fmt.Println("str1 and str2 are equal")
	}  else {
		fmt.Println("str1 and str2 are not equal")
	}
}
