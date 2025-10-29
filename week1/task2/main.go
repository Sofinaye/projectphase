package main

import (
	"regexp"
	"strings"
)

func WordFrequency(input string) map[string]int {
	reg, _ := regexp.Compile("[^a-zA-Z0-9\\s]+")
	cleaned := reg.ReplaceAllString(input, "")
	cleaned = strings.ToLower(cleaned)

	words := strings.Fields(cleaned)

	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}

	return frequency
}
