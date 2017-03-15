package main

import (
	"fmt"
)

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	ans := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			a := 0
			b := 0
			fmt.Scan(&a, &b)
			if a != 0 || b != 0 {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
