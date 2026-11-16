package layout

import (
	"fmt"
	"strings"
)

type Range struct {
	Start int
	End   int
}

type Boundary struct {
	Start int
	End   int
}

func FindRanges(text string, target string) []Range {
	ranges := []Range{}
	if target == "" {
		return ranges
	}
	searchStart := 0
	for {
		if searchStart >= len(text) {
			return ranges
		}
		relativeIndex := strings.Index(text[searchStart:], target)
		if relativeIndex == -1 {
			return ranges
		}
		matchIndex := searchStart + relativeIndex
		rangeFound := Range{
			Start: matchIndex,
			End:   matchIndex + len(target),
		}
		ranges = append(ranges, rangeFound)
		searchStart = matchIndex + len(target)
	}

}

func MergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{}
	}

	merged := []Range{}
	current := ranges[0]

	for _, next := range ranges[1:] {
		if next.Start <= current.End {
			if next.End > current.End {
				current.End = next.End
			}
			continue
		}

		merged = append(merged, current)
		current = next
	}

	merged = append(merged, current)

	return merged
}

func CalculateCharacterBoundaries(shapes map[rune][]string, text string) ([][]Boundary, error) {
	boundaries := [][]Boundary{}

	for row := 0; row < 8; row++ {
		currentColumn := 0
		rowBoundaries := []Boundary{}

		for _, char := range text {
			shape, ok := shapes[char]
			if !ok {
				return nil, fmt.Errorf("unsupported character: %q", char)
			}

			width := len(shape[row])

			boundary := Boundary{
				Start: currentColumn,
				End:   currentColumn + width,
			}

			rowBoundaries = append(rowBoundaries, boundary)
			currentColumn = boundary.End
		}

		boundaries = append(boundaries, rowBoundaries)
	}

	return boundaries, nil
}

func CalculateMultilineBoundaries(shapes map[rune][]string, text string) ([][][]Boundary, error) {
	lines := strings.Split(text, "\n")
	boundaries := [][][]Boundary{}

	for _, line := range lines {
		lineBoundaries, err := CalculateCharacterBoundaries(shapes, line)
		if err != nil {
			return nil, err
		}
		boundaries = append(boundaries, lineBoundaries)
	}
	return boundaries, nil
}

func CharacterRangeToRenderedRange(rowBoundaries []Boundary, start int, end int) (Boundary, error) {
	if start < 0 || end > len(rowBoundaries) || start >= end {
		return Boundary{}, fmt.Errorf("invalid character range")
	}

	return Boundary{
		Start: rowBoundaries[start].Start,
		End:   rowBoundaries[end-1].End,
	}, nil
}
