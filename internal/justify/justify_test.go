package justify

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LineLayout
	}{
		{
			name:     "empty string",
			input:    "",
			expected: LineLayout{},
		},
		{
			name:  "spaces only",
			input: "   ",
			expected: LineLayout{
				Leading: 3,
			},
		},
		{
			name:  "single space only",
			input: " ",
			expected: LineLayout{
				Leading: 1,
			},
		},
		{
			name:  "single word without spaces",
			input: "hello",
			expected: LineLayout{
				Words: []string{"hello"},
			},
		},
		{
			name:  "single word with leading space",
			input: " hello",
			expected: LineLayout{
				Words:   []string{"hello"},
				Leading: 1,
			},
		},
		{
			name:  "single word with multiple leading spaces",
			input: "   hello",
			expected: LineLayout{
				Words:   []string{"hello"},
				Leading: 3,
			},
		},
		{
			name:  "single word with trailing space",
			input: "hello ",
			expected: LineLayout{
				Words:    []string{"hello"},
				Trailing: 1,
			},
		},
		{
			name:  "single word with multiple trailing spaces",
			input: "hello   ",
			expected: LineLayout{
				Words:    []string{"hello"},
				Trailing: 3,
			},
		},
		{
			name:  "single word with leading and trailing spaces",
			input: "  hello   ",
			expected: LineLayout{
				Words:    []string{"hello"},
				Leading:  2,
				Trailing: 3,
			},
		},
		{
			name:  "two words with single space gap",
			input: "hello world",
			expected: LineLayout{
				Words: []string{"hello", "world"},
				Gaps:  []int{1},
			},
		},
		{
			name:  "two words with multiple spaces gap",
			input: "hello   world",
			expected: LineLayout{
				Words: []string{"hello", "world"},
				Gaps:  []int{3},
			},
		},
		{
			name:  "multiple words with single space gaps",
			input: "one two three",
			expected: LineLayout{
				Words: []string{"one", "two", "three"},
				Gaps:  []int{1, 1},
			},
		},
		{
			name:  "multiple words with varying gaps",
			input: "a   b  c    d",
			expected: LineLayout{
				Words: []string{"a", "b", "c", "d"},
				Gaps:  []int{3, 2, 4},
			},
		},
		{
			name:  "multiple words with leading, varying gaps, and trailing spaces",
			input: "   hello   brave   new   world  ",
			expected: LineLayout{
				Words:    []string{"hello", "brave", "new", "world"},
				Leading:  3,
				Gaps:     []int{3, 3, 3},
				Trailing: 2,
			},
		},
		{
			name:  "single character words",
			input: "x y z",
			expected: LineLayout{
				Words: []string{"x", "y", "z"},
				Gaps:  []int{1, 1},
			},
		},
		{
			name:  "punctuation and special characters in words",
			input: "hello, world! 123",
			expected: LineLayout{
				Words: []string{"hello,", "world!", "123"},
				Gaps:  []int{1, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("ParseLine(%q) = %#v, expected %#v", tt.input, got, tt.expected)
			}
		})
	}
}
