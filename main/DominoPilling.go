package main

import (
	"fmt"
	"math"
)

func main() {
	var m, n int
	fmt.Scan(&m, &n)
	fmt.Println(int(math.Floor(float64(m) * float64(n) / 2.0)))
}
