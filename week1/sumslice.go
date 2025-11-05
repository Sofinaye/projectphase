package main

import "fmt"

func Sum_slice(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}

	return sum
}

func main() {
	fmt.Println(Sum_slice([]int{1, 2, 3, 4, 5}))
}
