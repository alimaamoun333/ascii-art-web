package layout

import (
	"reflect"
	"testing"
)

func TestFindRanges(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		target   string
		expected []Range
	}{
		{
			name:   "exact match",
			text:   "kit",
			target: "kit",
			expected: []Range{
				{Start: 0, End: 3},
			},
		},
		{
			name:     "no exact match",
			text:     "hello",
			target:   "kit",
			expected: []Range{},
		},
		{
			name:     "empty target",
			text:     "hello",
			target:   "",
			expected: []Range{},
		},
		{
			name:   "multiple occurrences",
			text:   "banana",
			target: "an",
			expected: []Range{
				{Start: 1, End: 3},
				{Start: 3, End: 5},
			},
		},
		{
			name:   "overlapping occurrences are skipped",
			text:   "aaaaa",
			target: "aa",
			expected: []Range{
				{Start: 0, End: 2},
				{Start: 2, End: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRanges(tt.text, tt.target)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "overlapping ranges",
			input: []Range{
				{Start: 1, End: 4},
				{Start: 3, End: 6},
			},
			expected: []Range{
				{Start: 1, End: 6},
			},
		},
		{
			name: "touching ranges",
			input: []Range{
				{Start: 1, End: 3},
				{Start: 3, End: 5},
			},
			expected: []Range{
				{Start: 1, End: 5},
			},
		},
		{
			name: "contained range",
			input: []Range{
				{Start: 1, End: 6},
				{Start: 2, End: 4},
			},
			expected: []Range{
				{Start: 1, End: 6},
			},
		},
		{
			name: "separate ranges",
			input: []Range{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
			expected: []Range{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeRanges(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCalculateCharacterBoundariesSingleCharacter(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]Boundary{
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
func TestCalculateCharacterBoundariesMultipleCharacters(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
		'B': {
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "AB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]Boundary{
		{{Start: 0, End: 2}, {Start: 2, End: 5}},
		{{Start: 0, End: 3}, {Start: 3, End: 4}},
		{{Start: 0, End: 1}, {Start: 1, End: 3}},
		{{Start: 0, End: 4}, {Start: 4, End: 7}},
		{{Start: 0, End: 2}, {Start: 2, End: 3}},
		{{Start: 0, End: 3}, {Start: 3, End: 5}},
		{{Start: 0, End: 1}, {Start: 1, End: 4}},
		{{Start: 0, End: 4}, {Start: 4, End: 5}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCalculateCharacterBoundariesUnsupportedCharacter(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "AZ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil boundaries, got %v", got)
	}
}

func TestCharacterRangeToRenderedRange(t *testing.T) {
	boundaries := []Boundary{
		{Start: 0, End: 4},
		{Start: 4, End: 9},
		{Start: 9, End: 15},
	}

	tests := []struct {
		name        string
		boundaries  []Boundary
		start       int
		end         int
		expected    Boundary
		expectError bool
	}{
		{
			name:        "first character",
			boundaries:  boundaries,
			start:       0,
			end:         1,
			expected:    Boundary{Start: 0, End: 4},
			expectError: false,
		},
		{
			name:        "middle character",
			boundaries:  boundaries,
			start:       1,
			end:         2,
			expected:    Boundary{Start: 4, End: 9},
			expectError: false,
		},
		{
			name:        "last character",
			boundaries:  boundaries,
			start:       2,
			end:         3,
			expected:    Boundary{Start: 9, End: 15},
			expectError: false,
		},
		{
			name:        "multiple characters span",
			boundaries:  boundaries,
			start:       0,
			end:         2,
			expected:    Boundary{Start: 0, End: 9},
			expectError: false,
		},
		{
			name:        "all characters",
			boundaries:  boundaries,
			start:       0,
			end:         3,
			expected:    Boundary{Start: 0, End: 15},
			expectError: false,
		},
		{
			name:        "start less than zero",
			boundaries:  boundaries,
			start:       -1,
			end:         2,
			expected:    Boundary{},
			expectError: true,
		},
		{
			name:        "end greater than length",
			boundaries:  boundaries,
			start:       0,
			end:         4,
			expected:    Boundary{},
			expectError: true,
		},
		{
			name:        "start equals end",
			boundaries:  boundaries,
			start:       1,
			end:         1,
			expected:    Boundary{},
			expectError: true,
		},
		{
			name:        "start greater than end",
			boundaries:  boundaries,
			start:       2,
			end:         1,
			expected:    Boundary{},
			expectError: true,
		},
		{
			name:        "empty boundaries",
			boundaries:  []Boundary{},
			start:       0,
			end:         1,
			expected:    Boundary{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CharacterRangeToRenderedRange(tt.boundaries, tt.start, tt.end)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if got != tt.expected {
					t.Fatalf("expected %v, got %v", tt.expected, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCalculateMultilineBoundaries(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
		'B': {
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
		},
	}

	expectedLineA := [][]Boundary{
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
	}
	expectedLineB := [][]Boundary{
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
	}
	expectedEmptyLine := [][]Boundary{
		{}, {}, {}, {}, {}, {}, {}, {},
	}

	t.Run("single line", func(t *testing.T) {
		got, err := CalculateMultilineBoundaries(shapes, "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := [][][]Boundary{expectedLineA}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		got, err := CalculateMultilineBoundaries(shapes, "A\nB")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := [][][]Boundary{expectedLineA, expectedLineB}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})

	t.Run("empty lines between characters", func(t *testing.T) {
		got, err := CalculateMultilineBoundaries(shapes, "A\n\nB")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := [][][]Boundary{expectedLineA, expectedEmptyLine, expectedLineB}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})

	t.Run("unsupported character", func(t *testing.T) {
		got, err := CalculateMultilineBoundaries(shapes, "A\nZ")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if got != nil {
			t.Fatalf("expected nil boundaries, got %v", got)
		}
	})
}
