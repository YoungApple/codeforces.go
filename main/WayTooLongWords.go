package main

import (
	"bytes"
	"fmt"
)

func main() {
	var n int
	var input string
	fmt.Scan(&n)

	var output bytes.Buffer
	for n > 0 {
		fmt.Scan(&input)
		size := len(input)
		if size <= 10 {
			output.WriteString(input + "\n")
		} else {
			output.WriteString(fmt.Sprintf("%c%d%c\n", input[0], size-2, input[size-1]))
		}
		n--
	}
	fmt.Println(output.String())
}
