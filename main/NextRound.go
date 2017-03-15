package main

import (
	"fmt"
)

func main() {
	var n, k, kth int
	result := -1
	fmt.Scan(&n, &k)
	//var result int = 0
	for i := 0; i < n; i++ {
		var temp int
		fmt.Scan(&temp)
		if i == k-1 {
			// record kth number
			kth = temp
		}
		if temp > 0 && ((i < k) || (i >= k && kth == temp)) {
			// continue
		} else {
			result = i
			break
		}
	}
	if result == -1 {
		result = n
	}
	fmt.Println(result)
}
