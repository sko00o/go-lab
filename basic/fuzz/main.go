package main

import (
	"strings"
)

func MyIndexAny(s, chars string) int {
	for i, c := range s {
		if strings.ContainsRune(chars, c) {
			return i
		}
	}
	return -1
}
