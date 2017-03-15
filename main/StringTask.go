package main

import (
	"bytes"
	"fmt"
	"strings"
)

var vowels = map[string]bool{
	"a": true,
	"o": true,
	"y": true,
	"e": true,
	"u": true,
	"i": true,
}

func main() {
	var input string
	var buffer bytes.Buffer
	fmt.Scan(&input)
	for _, s := range strings.ToLower(input) {
		_, exist := vowels[string(s)]
		if !exist {
			buffer.WriteString("." + string(s))
		}
	}
	fmt.Println(buffer.String())
}
