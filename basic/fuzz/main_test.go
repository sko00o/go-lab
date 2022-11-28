package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMyIndexAny(t *testing.T) {
	type args struct {
		s, chars string
		want     int
	}
	tests := []args{
		{"abc", "xyz", -1},
		{"\x70\x71\x72", "\x73", -1},
		{"ab@c", "x@yz", 2},
		{"a\x20c", "xyz\x20", 1},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("(%q, %q)", tc.s, tc.chars), func(t *testing.T) {
			if got := MyIndexAny(tc.s, tc.chars); got != tc.want {
				t.Errorf("%d, want %d", got, tc.want)
			}
		})
	}
}

func FuzzMyIndexAny(f *testing.F) {
	tests := []struct {
		s, chars string
	}{
		{"", ""},
		{"", "a"},
		{"", "abc"},
		{"a", ""},
		{"a", "a"},
		{"aaa", "a"},
		{"abc", "xyz"},
		{"ab@c", "x@yz"},
		{"\x70\x71\x72", "\x73"},
	}
	for _, tc := range tests {
		// add to seed corpus
		f.Add(tc.s, tc.chars)
	}
	f.Fuzz(func(t *testing.T, s, chars string) {
		if got, want := MyIndexAny(s, chars), strings.IndexAny(s, chars); got != want {
			t.Errorf("(%q, %q) got %d, want %d", s, chars, got, want)
		}
	})
}
