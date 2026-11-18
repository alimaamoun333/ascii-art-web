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

func TestMarkColored(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		target   string
		expected []bool
	}{
		{
			name:     "empty text and empty target",
			text:     "",
			target:   "",
			expected: []bool{},
		},
		{
			name:     "empty text and non-empty target",
			text:     "",
			target:   "a",
			expected: []bool{},
		},
		{
			name:     "non-empty text and empty target",
			text:     "hello",
			target:   "",
			expected: []bool{false, false, false, false, false},
		},
		{
			name:     "target longer than text",
			text:     "hi",
			target:   "hello",
			expected: []bool{false, false},
		},
		{
			name:     "target not found",
			text:     "hello",
			target:   "world",
			expected: []bool{false, false, false, false, false},
		},
		{
			name:     "exact full match",
			text:     "cat",
			target:   "cat",
			expected: []bool{true, true, true},
		},
		{
			name:     "match at beginning",
			text:     "cat dog",
			target:   "cat",
			expected: []bool{true, true, true, false, false, false, false},
		},
		{
			name:     "match at end",
			text:     "dog cat",
			target:   "cat",
			expected: []bool{false, false, false, false, true, true, true},
		},
		{
			name:     "match in middle",
			text:     "a cat b",
			target:   "cat",
			expected: []bool{false, false, true, true, true, false, false},
		},
		{
			name:     "multiple occurrences disjoint",
			text:     "cat and cat",
			target:   "cat",
			expected: []bool{true, true, true, false, false, false, false, false, true, true, true},
		},
		{
			name:     "adjacent occurrences",
			text:     "catcat",
			target:   "cat",
			expected: []bool{true, true, true, true, true, true},
		},
		{
			name:     "overlapping pattern advances non-overlapping",
			text:     "aaaa",
			target:   "aa",
			expected: []bool{true, true, true, true},
		},
		{
			name:     "overlapping pattern with odd length",
			text:     "aaaaa",
			target:   "aa",
			expected: []bool{true, true, true, true, false},
		},
		{
			name:     "case-sensitive matching",
			text:     "Hello hello",
			target:   "hello",
			expected: []bool{false, false, false, false, false, false, true, true, true, true, true},
		},
		{
			name:     "target with spaces and punctuation",
			text:     "hello, world!",
			target:   ", ",
			expected: []bool{false, false, false, false, false, true, true, false, false, false, false, false, false},
		}, {
			name:     "match across words and spaces",
			text:     "Go   is",
			target:   "o   i",
			expected: []bool{false, true, true, true, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarkColored(tt.text, tt.target)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("MarkColored(%q, %q) = %v, expected %v", tt.text, tt.target, got, tt.expected)
			}
		})
	}
}

func TestSplitSegments(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		colored  []bool
		expected []Segment
	}{
		{
			name:     "empty text",
			text:     "",
			colored:  []bool{},
			expected: []Segment{},
		},
		{
			name:    "single character uncolored",
			text:    "a",
			colored: []bool{false},
			expected: []Segment{
				{Text: "a", Colored: false},
			},
		},
		{
			name:    "single character colored",
			text:    "a",
			colored: []bool{true},
			expected: []Segment{
				{Text: "a", Colored: true},
			},
		},
		{
			name:    "all uncolored",
			text:    "hello",
			colored: []bool{false, false, false, false, false},
			expected: []Segment{
				{Text: "hello", Colored: false},
			},
		},
		{
			name:    "all colored",
			text:    "hello",
			colored: []bool{true, true, true, true, true},
			expected: []Segment{
				{Text: "hello", Colored: true},
			},
		},
		{
			name:    "middle colored (user example)",
			text:    "Hello",
			colored: []bool{false, true, true, true, false},
			expected: []Segment{
				{Text: "H", Colored: false},
				{Text: "ell", Colored: true},
				{Text: "o", Colored: false},
			},
		},
		{
			name:    "prefix colored",
			text:    "hello",
			colored: []bool{true, true, false, false, false},
			expected: []Segment{
				{Text: "he", Colored: true},
				{Text: "llo", Colored: false},
			},
		},
		{
			name:    "suffix colored",
			text:    "hello",
			colored: []bool{false, false, false, true, true},
			expected: []Segment{
				{Text: "hel", Colored: false},
				{Text: "lo", Colored: true},
			},
		},
		{
			name:    "alternating colors",
			text:    "abc",
			colored: []bool{true, false, true},
			expected: []Segment{
				{Text: "a", Colored: true},
				{Text: "b", Colored: false},
				{Text: "c", Colored: true},
			},
		},
		{
			name:    "with spaces",
			text:    "a b",
			colored: []bool{false, true, false},
			expected: []Segment{
				{Text: "a", Colored: false},
				{Text: " ", Colored: true},
				{Text: "b", Colored: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitSegments(tt.text, tt.colored)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("SplitSegments(%q, %v) = %#v, expected %#v", tt.text, tt.colored, got, tt.expected)
			}
		})
	}
}

func TestSplitTextBlocks(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		colored   []bool
		expected  []TextBlock
		expectErr bool
	}{
		{
			name:     "empty text",
			text:     "",
			colored:  []bool{},
			expected: []TextBlock{},
		},
		{
			name:    "single word without spaces",
			text:    "hello",
			colored: make([]bool, len("hello")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "single space only",
			text:    " ",
			colored: make([]bool, 1),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: " ", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "multiple spaces only",
			text:    "   ",
			colored: make([]bool, 3),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "   ", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "two words with single space gap",
			text:    "a b",
			colored: make([]bool, len("a b")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "a", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: " ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "b", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "two words with multiple spaces gap",
			text:    "hello   world",
			colored: make([]bool, len("hello   world")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "   ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "world", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "single word with leading spaces",
			text:    "  hello",
			colored: make([]bool, len("  hello")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "  ", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "single word with trailing spaces",
			text:    "hello  ",
			colored: make([]bool, len("hello  ")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "  ", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "multiple words with leading, internal, and trailing spaces",
			text:    "  hello   world  ",
			colored: make([]bool, len("  hello   world  ")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "  ", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "   ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "world", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "  ", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "multiple words with multiple varying gaps",
			text:    "one  two   three",
			colored: make([]bool, len("one  two   three")),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "one", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "  ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "two", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: "   ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "three", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:    "single word with middle colored (user example)",
			text:    "Hello",
			colored: []bool{false, true, true, true, false},
			expected: []TextBlock{
				{
					Segments: []Segment{
						{Text: "H", Colored: false},
						{Text: "ell", Colored: true},
						{Text: "o", Colored: false},
					},
					Justify: false,
				},
			},
		},
		{
			name:    "two words with second word colored",
			text:    "hello world",
			colored: MarkColored("hello world", "world"),
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "hello", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []Segment{{Text: " ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []Segment{{Text: "world", Colored: true}},
					Justify:  false,
				},
			},
		},
		{
			name:    "match across word and internal spaces",
			text:    "Go   is",
			colored: MarkColored("Go   is", "o   i"),
			expected: []TextBlock{
				{
					Segments: []Segment{
						{Text: "G", Colored: false},
						{Text: "o", Colored: true},
					},
					Justify: false,
				},
				{
					Segments: []Segment{
						{Text: "   ", Colored: true},
					},
					Justify: true,
				},
				{
					Segments: []Segment{
						{Text: "i", Colored: true},
						{Text: "s", Colored: false},
					},
					Justify: false,
				},
			},
		},
		{
			name:    "partially colored internal space",
			text:    "a  b",
			colored: []bool{true, true, false, false},
			expected: []TextBlock{
				{
					Segments: []Segment{{Text: "a", Colored: true}},
					Justify:  false,
				},
				{
					Segments: []Segment{
						{Text: " ", Colored: true},
						{Text: " ", Colored: false},
					},
					Justify: true,
				},
				{
					Segments: []Segment{{Text: "b", Colored: false}},
					Justify:  false,
				},
			},
		},
		{
			name:      "mismatched length returns error",
			text:      "hello",
			colored:   []bool{false},
			expected:  []TextBlock{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitTextBlocks(tt.text, tt.colored)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("SplitTextBlocks(%q) expected error, got nil", tt.text)
				}
				return
			}
			if err != nil {
				t.Fatalf("SplitTextBlocks(%q) unexpected error: %v", tt.text, err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("SplitTextBlocks(%q) = %#v, expected %#v", tt.text, got, tt.expected)
			}
		})
	}
}
