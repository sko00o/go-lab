package benchmark

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplit(t *testing.T) {
	tests := map[string]struct {
		input string
		sep   string
		want  []string
	}{
		"simple":       {"a/b/c", "/", []string{"a", "b", "c"}},
		"wrong sep":    {"a/b/c", ",", []string{"a/b/c"}},
		"trailing sep": {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
	}

	for name, tc := range tests {
		// sub test add in go1.7
		t.Run(name, func(t *testing.T) {
			got := split(tc.input, tc.sep)
			/* if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			} */
			// use go-cmp is better
			// show diff when pointers in struct
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

// Test one func
// go test -run="^TestSplit$" -v

// Individual sub test cases can be executed directly
// go test -run=.*/trailing -v
