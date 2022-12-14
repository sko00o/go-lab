package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
        return s, errors.New("input is not valid UTF-8")
    }

	return reverse(s), nil
}

func reverse(s string) string {
	// b := []byte(s)
	b := []rune(s)

	// debug
	fmt.Printf("input: %q\nrunes: %q\n", s, b)

	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev, revErr := Reverse(input)
	doubleRev, doubleErr := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reverse: %q, err: %v\n", rev, revErr)
	fmt.Printf("reverse again: %q, err: %v\n", doubleRev, doubleErr)
}
