package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s1 := "mom"
	s2 := "topspot"
	s3 := "tada"

	fmt.Printf("is %v a palindrome? \t %t\n", s1, isPalidrome(s1))
	fmt.Printf("is %v a palindrome? \t %t\n", s2, isPalidrome(s2))
	fmt.Printf("is %v a palindrome? \t %t\n", s3, isPalidrome(s3))
}

func isPalidrome(s string) bool {
	r := utf8.RuneCountInString(s)
	mid := r / 2
	for i := 0; i < mid; i++ {
		if s[i] != s[r-i-1] {

			return false
		}

	}
	return true

}
