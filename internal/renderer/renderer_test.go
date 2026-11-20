package renderer

import (
	"ascii-art-web/internal/justify"
	"ascii-art-web/internal/layout"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// If running from internal/renderer directory, change working directory
	// to project root so that banners/ can be found.
	if _, err := os.Stat("banners/standard.txt"); os.IsNotExist(err) {
		if err := os.Chdir("../.."); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func buildTestShapes() map[rune][]string {
	return map[rune][]string{
		'A': {"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7"},
		'B': {"B0", "B1", "B2", "B3", "B4", "B5", "B6", "B7"},
		' ': {"S0", "S1", "S2", "S3", "S4", "S5", "S6", "S7"},
	}
}

func TestRenderCharacterExistingCharacter(t *testing.T) {
	shapes := buildTestShapes()

	shape, err := RenderCharacter(shapes, 'A')
	if err != nil {
		t.Fatalf("the character not found %v", err)
	}
	if len(shape) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(shape))
	}
	if shape[0] != "A0" {
		t.Fatalf("expected first row A0, got %q", shape[0])
	}
}

func TestRenderCharacterNotExistingCharacter(t *testing.T) {
	shapes := buildTestShapes()
	shape, err := RenderCharacter(shapes, 'Z')
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if shape != nil {
		t.Fatalf("expected nil shape, got %v", shape)
	}

}

func TestRenderLineMultipleCharacters(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := RenderLine(shapes, "AB")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(rows))
	}

	if rows[0] != "A0B0" {
		t.Fatalf("expected first row A0B0, got %q", rows[0])
	}
}

func TestRenderLineUnsupportedCharacter(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := RenderLine(shapes, "AZ")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderLineEmptyCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := RenderLine(shapes, "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderLineSpaceCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := RenderLine(shapes, " ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(rows))
	}

	if rows[0] != "S0" {
		t.Fatalf("expected first row S0, got %q", rows[0])
	}
}

func TestRenderMultipleLines(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "A\nB")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 16 {
		t.Fatalf("expected 16 rows, got %d", len(rows))
	}

	if rows[0] != "A0" {
		t.Fatalf("expected row 0 to be A0, got %q", rows[0])
	}

	if rows[8] != "B0" {
		t.Fatalf("expected row 8 to be B0, got %q", rows[8])
	}
}

func TestRenderConsecutiveNewlines(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "A\n\nB")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(rows) != 24 {
		t.Fatalf("expected 24 rows, got %d", len(rows))
	}

	if rows[0] != "A0" {
		t.Fatalf("expected row 0 to be A0, got %q", rows[0])
	}

	if rows[8] != "" {
		t.Fatalf("expected row 8 to be empty, got %q", rows[8])
	}

	if rows[15] != "" {
		t.Fatalf("expected row 15 to be empty, got %q", rows[15])
	}

	if rows[16] != "B0" {
		t.Fatalf("expected row 16 to be B0, got %q", rows[16])
	}
}

func TestRenderUnsupportedCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := Render(shapes, "A\nZ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderEmptyText(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderCharacterUnsupported(t *testing.T) {
	shapes := map[rune][]string{}

	_, err := RenderCharacter(shapes, '😀')

	if !errors.Is(err, ErrUnsupportedCharacter) {
		t.Fatalf(
			"expected ErrUnsupportedCharacter, got %v",
			err,
		)
	}
}

func TestSegmentRow(t *testing.T) {
	tests := []struct {
		name          string
		row           string
		coloredRanges []layout.Boundary
		expected      []Segment
		expectError   bool
	}{
		{
			name:          "empty row and empty ranges",
			row:           "",
			coloredRanges: []layout.Boundary{},
			expected:      []Segment{},
			expectError:   false,
		},
		{
			name:          "nil ranges returns uncolored row",
			row:           "Hello World",
			coloredRanges: nil,
			expected: []Segment{
				{Text: "Hello World", Colored: false},
			},
			expectError: false,
		},
		{
			name:          "no colored ranges returns entire row uncolored",
			row:           "Hello World",
			coloredRanges: []layout.Boundary{},
			expected: []Segment{
				{Text: "Hello World", Colored: false},
			},
			expectError: false,
		},
		{
			name: "entire row colored",
			row:  "Hello",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 5},
			},
			expected: []Segment{
				{Text: "Hello", Colored: true},
			},
			expectError: false,
		},
		{
			name: "prefix colored",
			row:  "HelloWorld",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 5},
			},
			expected: []Segment{
				{Text: "Hello", Colored: true},
				{Text: "World", Colored: false},
			},
			expectError: false,
		},
		{
			name: "suffix colored",
			row:  "HelloWorld",
			coloredRanges: []layout.Boundary{
				{Start: 5, End: 10},
			},
			expected: []Segment{
				{Text: "Hello", Colored: false},
				{Text: "World", Colored: true},
			},
			expectError: false,
		},
		{
			name: "middle segment colored",
			row:  "PreMidPost",
			coloredRanges: []layout.Boundary{
				{Start: 3, End: 6},
			},
			expected: []Segment{
				{Text: "Pre", Colored: false},
				{Text: "Mid", Colored: true},
				{Text: "Post", Colored: false},
			},
			expectError: false,
		},
		{
			name: "multiple disjoint colored segments",
			row:  "ABCDEFGHIJ",
			coloredRanges: []layout.Boundary{
				{Start: 1, End: 3},
				{Start: 6, End: 8},
			},
			expected: []Segment{
				{Text: "A", Colored: false},
				{Text: "BC", Colored: true},
				{Text: "DEF", Colored: false},
				{Text: "GH", Colored: true},
				{Text: "IJ", Colored: false},
			},
			expectError: false,
		},
		{
			name: "adjacent colored segments without uncolored gap",
			row:  "ABCDEF",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 3},
				{Start: 3, End: 6},
			},
			expected: []Segment{
				{Text: "ABC", Colored: true},
				{Text: "DEF", Colored: true},
			},
			expectError: false,
		},
		{
			name: "multiple adjacent segments covering whole row",
			row:  "123456",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 2},
				{Start: 2, End: 4},
				{Start: 4, End: 6},
			},
			expected: []Segment{
				{Text: "12", Colored: true},
				{Text: "34", Colored: true},
				{Text: "56", Colored: true},
			},
			expectError: false,
		},
		{
			name: "single character colored in row",
			row:  "12345",
			coloredRanges: []layout.Boundary{
				{Start: 2, End: 3},
			},
			expected: []Segment{
				{Text: "12", Colored: false},
				{Text: "3", Colored: true},
				{Text: "45", Colored: false},
			},
			expectError: false,
		},
		{
			name: "negative start index",
			row:  "Hello",
			coloredRanges: []layout.Boundary{
				{Start: -1, End: 3},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "overlapping ranges",
			row:  "ABCDEF",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 4},
				{Start: 2, End: 5},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "unsorted ranges",
			row:  "ABCDEF",
			coloredRanges: []layout.Boundary{
				{Start: 3, End: 5},
				{Start: 0, End: 2},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "zero length range start equals end",
			row:  "ABCDEF",
			coloredRanges: []layout.Boundary{
				{Start: 2, End: 2},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "inverted range start greater than end",
			row:  "ABCDEF",
			coloredRanges: []layout.Boundary{
				{Start: 4, End: 2},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "end index exceeds row length",
			row:  "ABC",
			coloredRanges: []layout.Boundary{
				{Start: 1, End: 4},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "non-empty range on empty row",
			row:  "",
			coloredRanges: []layout.Boundary{
				{Start: 0, End: 1},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "preserve spaces around colored text",
			row:  "  ABC  ",
			coloredRanges: []layout.Boundary{
				{Start: 2, End: 5},
			},
			expected: []Segment{
				{Text: "  ", Colored: false},
				{Text: "ABC", Colored: true},
				{Text: "  ", Colored: false},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SegmentRow(tt.row, tt.coloredRanges)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidColoredRanges) {
					t.Fatalf("expected ErrInvalidColoredRanges, got %v", err)
				}
				if got != nil {
					t.Fatalf("expected nil segments on error, got %v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %#v, got %#v", tt.expected, got)
			}
		})
	}
}

func TestRenderSelected(t *testing.T) {
	shapes := buildTestShapes()

	t.Run("partial match at start of line", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AB", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: true},
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("partial match at end of line", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AB", "B")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
					{Text: fmt.Sprintf("B%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("full match", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AB", "AB")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%dB%d", i, i), Colored: true},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("target not found", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AB", "C")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%dB%d", i, i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("empty target", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AB", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%dB%d", i, i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("target matching space", func(t *testing.T) {
		got, err := RenderSelected(shapes, "A B", " ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
					{Text: fmt.Sprintf("S%d", i), Colored: true},
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("multiple occurrences on single line", func(t *testing.T) {
		got, err := RenderSelected(shapes, "ABA", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: true},
					{Text: fmt.Sprintf("B%d", i), Colored: false},
					{Text: fmt.Sprintf("A%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}
	})

	t.Run("multiline text with CRLF normalization", func(t *testing.T) {
		got, err := RenderSelected(shapes, "A\r\nB", "B")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 16 {
			t.Fatalf("expected 16 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expectedA := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expectedA) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expectedA, got[i])
			}
		}
		for i := 0; i < 8; i++ {
			expectedB := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(got[8+i], expectedB) {
				t.Fatalf("row %d: expected %#v, got %#v", 8+i, expectedB, got[8+i])
			}
		}
	})

	t.Run("multiline text with consecutive newlines", func(t *testing.T) {
		got, err := RenderSelected(shapes, "A\n\nB", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 24 {
			t.Fatalf("expected 24 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expectedA := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(got[i], expectedA) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expectedA, got[i])
			}
		}
		for i := 0; i < 8; i++ {
			expectedEmpty := Row{
				Segments: []Segment{},
			}
			if !reflect.DeepEqual(got[8+i], expectedEmpty) {
				t.Fatalf("row %d: expected %#v, got %#v", 8+i, expectedEmpty, got[8+i])
			}
		}
		for i := 0; i < 8; i++ {
			expectedB := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[16+i], expectedB) {
				t.Fatalf("row %d: expected %#v, got %#v", 16+i, expectedB, got[16+i])
			}
		}
	})

	t.Run("case-sensitive matching", func(t *testing.T) {
		caseShapes := buildTestShapes()
		caseShapes['a'] = []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7"}

		// Target "A" should only color 'A' and not 'a'
		got, err := RenderSelected(caseShapes, "Aa", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: true},
					{Text: fmt.Sprintf("a%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}

		// Target "a" should only color 'a' and not 'A'
		gotLower, err := RenderSelected(caseShapes, "Aa", "a")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotLower) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(gotLower))
		}
		for i := 0; i < 8; i++ {
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
					{Text: fmt.Sprintf("a%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(gotLower[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, gotLower[i])
			}
		}
	})

	t.Run("adjacent matches merging into one segment", func(t *testing.T) {
		got, err := RenderSelected(shapes, "BAAB", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(got))
		}
		for i := 0; i < 8; i++ {
			if len(got[i].Segments) != 3 {
				t.Fatalf("row %d: expected 3 segments (merged adjacent matches), got %d", i, len(got[i].Segments))
			}
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: false},
					{Text: fmt.Sprintf("A%dA%d", i, i), Colored: true},
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(got[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, got[i])
			}
		}

		// Entire line of adjacent matches merges into single segment
		gotAll, err := RenderSelected(shapes, "AA", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for i := 0; i < 8; i++ {
			if len(gotAll[i].Segments) != 1 {
				t.Fatalf("row %d: expected 1 merged segment, got %d", i, len(gotAll[i].Segments))
			}
			expected := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%dA%d", i, i), Colored: true},
				},
			}
			if !reflect.DeepEqual(gotAll[i], expected) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expected, gotAll[i])
			}
		}
	})

	t.Run("matching independently across logical lines, including targets containing newlines", func(t *testing.T) {
		// 1. Independent matching across logical lines:
		// Target "A" matches at index 0 on the first line and index 1 on the second line.
		gotIndependent, err := RenderSelected(shapes, "AB\nBA", "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotIndependent) != 16 {
			t.Fatalf("expected 16 rows, got %d", len(gotIndependent))
		}
		for i := 0; i < 8; i++ {
			expectedLine1 := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: true},
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(gotIndependent[i], expectedLine1) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expectedLine1, gotIndependent[i])
			}
			expectedLine2 := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: false},
					{Text: fmt.Sprintf("A%d", i), Colored: true},
				},
			}
			if !reflect.DeepEqual(gotIndependent[8+i], expectedLine2) {
				t.Fatalf("row %d: expected %#v, got %#v", 8+i, expectedLine2, gotIndependent[8+i])
			}
		}

		// 2. Target containing newline:
		// Matching is per logical line, so a target spanning across newlines ("A\nB") matches no individual line.
		gotNewlineTarget, err := RenderSelected(shapes, "A\nB", "A\nB")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotNewlineTarget) != 16 {
			t.Fatalf("expected 16 rows, got %d", len(gotNewlineTarget))
		}
		for i := 0; i < 8; i++ {
			expectedA := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(gotNewlineTarget[i], expectedA) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expectedA, gotNewlineTarget[i])
			}
			expectedB := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(gotNewlineTarget[8+i], expectedB) {
				t.Fatalf("row %d: expected %#v, got %#v", 8+i, expectedB, gotNewlineTarget[8+i])
			}
		}

		// 3. Target is only a newline "\n":
		gotOnlyNewline, err := RenderSelected(shapes, "A\nB", "\n")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotOnlyNewline) != 16 {
			t.Fatalf("expected 16 rows, got %d", len(gotOnlyNewline))
		}
		for i := 0; i < 8; i++ {
			expectedA := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("A%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(gotOnlyNewline[i], expectedA) {
				t.Fatalf("row %d: expected %#v, got %#v", i, expectedA, gotOnlyNewline[i])
			}
			expectedB := Row{
				Segments: []Segment{
					{Text: fmt.Sprintf("B%d", i), Colored: false},
				},
			}
			if !reflect.DeepEqual(gotOnlyNewline[8+i], expectedB) {
				t.Fatalf("row %d: expected %#v, got %#v", 8+i, expectedB, gotOnlyNewline[8+i])
			}
		}
	})

	t.Run("empty text returns error", func(t *testing.T) {
		got, err := RenderSelected(shapes, "", "A")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, ErrEmptyText) {
			t.Fatalf("expected ErrEmptyText, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil rows, got %v", got)
		}
	})

	t.Run("unsupported character in text returns error", func(t *testing.T) {
		got, err := RenderSelected(shapes, "AZ", "A")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, ErrUnsupportedCharacter) {
			t.Fatalf("expected ErrUnsupportedCharacter, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil rows, got %v", got)
		}
	})
}

func TestGetShapes_Success(t *testing.T) {
	ResetShapes()
	defer ResetShapes()

	std, err := GetShapes("standard")
	if err != nil {
		t.Fatalf("unexpected error getting standard shapes: %v", err)
	}
	if len(std) != 95 {
		t.Errorf("expected 95 shapes, got %d", len(std))
	}

	shadow, err := GetShapes("shadow")
	if err != nil {
		t.Fatalf("unexpected error getting shadow shapes: %v", err)
	}
	if len(shadow) != 95 {
		t.Errorf("expected 95 shapes, got %d", len(shadow))
	}

	thinker, err := GetShapes("thinkertoy")
	if err != nil {
		t.Fatalf("unexpected error getting thinkertoy shapes: %v", err)
	}
	if len(thinker) != 95 {
		t.Errorf("expected 95 shapes, got %d", len(thinker))
	}
}

func TestGetShapes_LoadedOnce(t *testing.T) {
	ResetShapes()
	defer ResetShapes()

	// First call loads all shapes into variables
	std1, err := GetShapes("standard")
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	// Verify standardShapes, shadowShapes, and tinnkerShapes are populated
	if standardShapes == nil || shadowShapes == nil || tinnkerShapes == nil {
		t.Fatalf("expected all shape variables to be populated after first load")
	}

	// Subsequent calls return the variables directly
	std2, err := GetShapes("standard")
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	// Verify it returns the exact same map variable
	if reflect.ValueOf(std1).Pointer() != reflect.ValueOf(std2).Pointer() {
		t.Errorf("expected identical map reference from standardShapes variable")
	}
}

func TestGetShapes_Unsupported(t *testing.T) {
	ResetShapes()
	defer ResetShapes()

	_, err := GetShapes("unknown")
	if err == nil {
		t.Fatalf("expected error for unsupported banner, got nil")
	}
}

func TestGenerateArt_Success(t *testing.T) {
	ResetShapes()
	defer ResetShapes()

	art, err := GenerateArt("Hello", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if art == "" {
		t.Errorf("expected non-empty ASCII art")
	}
}

func TestRenderBlocks(t *testing.T) {
	shapes := buildTestShapes()

	makeExpectedRows := func(fn func(i int) []Segment) []Row {
		rows := make([]Row, 8)
		for i := 0; i < 8; i++ {
			rows[i] = Row{Segments: fn(i)}
		}
		return rows
	}

	tests := []struct {
		name      string
		blocks    []justify.TextBlock
		customMap map[rune][]string
		expected  []RenderedBlock
		expectErr bool
	}{
		{
			name:     "empty blocks list",
			blocks:   []justify.TextBlock{},
			expected: []RenderedBlock{},
		},
		{
			name: "single uncolored word block",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "A", Colored: false}},
					Justify:  false,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("A%d", i), Colored: false}}
					}),
					Justify: false,
				},
			},
		},
		{
			name: "single colored word block",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "A", Colored: true}},
					Justify:  false,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("A%d", i), Colored: true}}
					}),
					Justify: false,
				},
			},
		},
		{
			name: "word block with partial coloring",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{
						{Text: "A", Colored: false},
						{Text: "B", Colored: true},
					},
					Justify: false,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{
							{Text: fmt.Sprintf("A%d", i), Colored: false},
							{Text: fmt.Sprintf("B%d", i), Colored: true},
						}
					}),
					Justify: false,
				},
			},
		},
		{
			name: "single space block with justify true",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: " ", Colored: false}},
					Justify:  true,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
			},
		},
		{
			name: "multiple spaces split into individual space blocks",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "  ", Colored: false}},
					Justify:  true,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
			},
		},
		{
			name: "multiple spaces with mixed coloring",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{
						{Text: " ", Colored: true},
						{Text: " ", Colored: false},
					},
					Justify: true,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: true}}
					}),
					Justify: true,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
			},
		},
		{
			name: "sentence sequence of words and internal spaces",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "A", Colored: false}},
					Justify:  false,
				},
				{
					Segments: []justify.Segment{{Text: "  ", Colored: false}},
					Justify:  true,
				},
				{
					Segments: []justify.Segment{{Text: "B", Colored: true}},
					Justify:  false,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("A%d", i), Colored: false}}
					}),
					Justify: false,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: true,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("B%d", i), Colored: true}}
					}),
					Justify: false,
				},
			},
		},
		{
			name: "leading spaces with justify false",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "  ", Colored: false}},
					Justify:  false,
				},
			},
			expected: []RenderedBlock{
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: false,
				},
				{
					Rows: makeExpectedRows(func(i int) []Segment {
						return []Segment{{Text: fmt.Sprintf("S%d", i), Colored: false}}
					}),
					Justify: false,
				},
			},
		},
		{
			name: "unsupported character in word returns error",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "Z", Colored: false}},
					Justify:  false,
				},
			},
			expectErr: true,
		},
		{
			name: "unsupported space character returns error",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: " ", Colored: false}},
					Justify:  true,
				},
			},
			customMap: map[rune][]string{
				'A': {"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7"},
			},
			expectErr: true,
		},
		{
			name: "empty text block returns error",
			blocks: []justify.TextBlock{
				{
					Segments: []justify.Segment{{Text: "", Colored: false}},
					Justify:  false,
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetShapes := shapes
			if tt.customMap != nil {
				targetShapes = tt.customMap
			}

			got, err := RenderBlocks(targetShapes, tt.blocks)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("RenderBlocks() = %#v, expected %#v", got, tt.expected)
			}
		})
	}
}

func TestGenerateBlocks(t *testing.T) {
	ResetShapes()
	defer ResetShapes()

	t.Run("empty text returns single line with empty blocks", func(t *testing.T) {
		lines, err := GenerateBlocks("", "standard", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		if len(lines[0].Blocks) != 0 {
			t.Fatalf("expected 0 blocks, got %d", len(lines[0].Blocks))
		}
	})

	t.Run("single word uncolored", func(t *testing.T) {
		lines, err := GenerateBlocks("Hello", "standard", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		blocks := lines[0].Blocks
		if len(blocks) != 1 {
			t.Fatalf("expected 1 block, got %d", len(blocks))
		}
		if blocks[0].Justify {
			t.Errorf("expected Justify to be false for word block")
		}
		if len(blocks[0].Rows) != 8 {
			t.Fatalf("expected 8 rows, got %d", len(blocks[0].Rows))
		}
		for r, row := range blocks[0].Rows {
			if len(row.Segments) != 1 {
				t.Fatalf("row %d: expected 1 segment, got %d", r, len(row.Segments))
			}
			if row.Segments[0].Colored {
				t.Errorf("row %d: expected segment to be uncolored", r)
			}
		}
	})

	t.Run("multiple words with justify space", func(t *testing.T) {
		lines, err := GenerateBlocks("Hi there", "standard", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		blocks := lines[0].Blocks
		// "Hi" (false), " " (true), "there" (false) -> 3 blocks
		if len(blocks) != 3 {
			t.Fatalf("expected 3 blocks, got %d", len(blocks))
		}
		if blocks[0].Justify != false {
			t.Errorf("block 0: expected Justify false")
		}
		if blocks[1].Justify != true {
			t.Errorf("block 1: expected Justify true for internal space")
		}
		if blocks[2].Justify != false {
			t.Errorf("block 2: expected Justify false")
		}
	})

	t.Run("multiple spaces expand to individual justify blocks", func(t *testing.T) {
		lines, err := GenerateBlocks("A   B", "standard", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		blocks := lines[0].Blocks
		// "A" (false), 3 spaces (each true), "B" (false) -> 5 blocks
		if len(blocks) != 5 {
			t.Fatalf("expected 5 blocks, got %d", len(blocks))
		}
		if blocks[0].Justify != false || blocks[4].Justify != false {
			t.Errorf("expected words to have Justify false")
		}
		for i := 1; i <= 3; i++ {
			if !blocks[i].Justify {
				t.Errorf("block %d: expected Justify true", i)
			}
		}
	})

	t.Run("target coloring word", func(t *testing.T) {
		lines, err := GenerateBlocks("Hello World", "standard", "World")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		blocks := lines[0].Blocks
		if len(blocks) != 3 {
			t.Fatalf("expected 3 blocks, got %d", len(blocks))
		}
		// "Hello" uncolored
		for _, row := range blocks[0].Rows {
			for _, seg := range row.Segments {
				if seg.Colored {
					t.Errorf("expected 'Hello' to be uncolored")
				}
			}
		}
		// "World" colored
		for _, row := range blocks[2].Rows {
			for _, seg := range row.Segments {
				if !seg.Colored {
					t.Errorf("expected 'World' to be colored")
				}
			}
		}
	})

	t.Run("partial word coloring", func(t *testing.T) {
		lines, err := GenerateBlocks("Hello", "standard", "ell")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
		blocks := lines[0].Blocks
		if len(blocks) != 1 {
			t.Fatalf("expected 1 block, got %d", len(blocks))
		}
		for r, row := range blocks[0].Rows {
			if len(row.Segments) != 3 {
				t.Fatalf("row %d: expected 3 segments (H, ell, o), got %d", r, len(row.Segments))
			}
			if row.Segments[0].Colored {
				t.Errorf("row %d: expected 'H' to be uncolored", r)
			}
			if !row.Segments[1].Colored {
				t.Errorf("row %d: expected 'ell' to be colored", r)
			}
			if row.Segments[2].Colored {
				t.Errorf("row %d: expected 'o' to be uncolored", r)
			}
		}
	})

	t.Run("multiline text and CRLF normalization", func(t *testing.T) {
		lines, err := GenerateBlocks("Hello\r\nWorld", "standard", "World")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}
		if len(lines[0].Blocks) != 1 || len(lines[1].Blocks) != 1 {
			t.Fatalf("expected 1 block per line, got %d and %d", len(lines[0].Blocks), len(lines[1].Blocks))
		}

		// First line ("Hello") should have 8 rows and be uncolored
		if len(lines[0].Blocks[0].Rows) != 8 {
			t.Fatalf("line 1: expected 8 rows, got %d", len(lines[0].Blocks[0].Rows))
		}
		for r, row := range lines[0].Blocks[0].Rows {
			for _, seg := range row.Segments {
				if seg.Colored {
					t.Errorf("line 1 row %d: expected 'Hello' to be uncolored", r)
				}
			}
		}

		// Second line ("World") should have 8 rows and be colored
		if len(lines[1].Blocks[0].Rows) != 8 {
			t.Fatalf("line 2: expected 8 rows, got %d", len(lines[1].Blocks[0].Rows))
		}
		for r, row := range lines[1].Blocks[0].Rows {
			for _, seg := range row.Segments {
				if !seg.Colored {
					t.Errorf("line 2 row %d: expected 'World' to be colored", r)
				}
			}
		}
	})

	t.Run("different banners shadow and thinkertoy", func(t *testing.T) {
		for _, banner := range []string{"shadow", "thinkertoy"} {
			lines, err := GenerateBlocks("Go", banner, "")
			if err != nil {
				t.Fatalf("banner %s failed: %v", banner, err)
			}
			if len(lines) != 1 {
				t.Fatalf("banner %s: expected 1 line, got %d", banner, len(lines))
			}
			if len(lines[0].Blocks) != 1 {
				t.Fatalf("banner %s: expected 1 block, got %d", banner, len(lines[0].Blocks))
			}
		}
	})

	t.Run("invalid banner returns error", func(t *testing.T) {
		_, err := GenerateBlocks("Hello", "nonexistent", "")
		if err == nil {
			t.Fatalf("expected error for invalid banner, got nil")
		}
	})

	t.Run("unsupported character returns error", func(t *testing.T) {
		_, err := GenerateBlocks("Hello \x01", "standard", "")
		if !errors.Is(err, ErrUnsupportedCharacter) {
			t.Fatalf("expected ErrUnsupportedCharacter, got %v", err)
		}
	})

}
