package renderer

import (
	"ascii-art-web/internal/banner"
	"ascii-art-web/internal/justify"
	"ascii-art-web/internal/layout"
	"errors"
	"strings"
	"sync"
)

var (
	standardShapes map[rune][]string
	tinnkerShapes  map[rune][]string
	shadowShapes   map[rune][]string
	shapes         = new(sync.Once)
	shapesErr      error
)

// LoadShapes loads all banner styles once using sync.Once.
func LoadShapes() error {
	shapes.Do(func() {
		var err error
		standardShapes, err = banner.Load("banners/standard.txt")
		if err != nil {
			shapesErr = err
			return
		}
		shadowShapes, err = banner.Load("banners/shadow.txt")
		if err != nil {
			shapesErr = err
			return
		}
		tinnkerShapes, err = banner.Load("banners/thinkertoy.txt")
		if err != nil {
			shapesErr = err
			return
		}
	})
	return shapesErr
}

// GetShapes returns the banner map for the requested banner style by accessing the loaded variable.
func GetShapes(bannerName string) (map[rune][]string, error) {
	if err := LoadShapes(); err != nil {
		return nil, err
	}

	switch bannerName {
	case "standard":
		return standardShapes, nil
	case "shadow":
		return shadowShapes, nil
	case "thinkertoy":
		return tinnkerShapes, nil
	default:
		return nil, errors.New("unsupported banner")
	}
}

// ResetShapes resets the sync.Once and clears the loaded variables (useful for testing).
func ResetShapes() {
	shapes = new(sync.Once)
	shapesErr = nil
	standardShapes = nil
	shadowShapes = nil
	tinnkerShapes = nil
}

type Segment struct {
	Text    string
	Colored bool
}

type Row struct {
	Segments []Segment
}

type RenderedBlock struct {
	Rows    []Row
	Justify bool
}

var ErrUnsupportedCharacter = errors.New("unsupported character")
var ErrEmptyText = errors.New("text cannot be empty")
var ErrInvalidColoredRanges = errors.New("invalid colored ranges")

// RenderCharacter returns the glyph rows for one character.
func RenderCharacter(shapes map[rune][]string, char rune) ([]string, error) {
	shape, ok := shapes[char]

	if ok == false {
		return nil, ErrUnsupportedCharacter
	}

	return shape, nil
}

// RenderLine renders one non-empty logical line into 8 ASCII-art rows.
func RenderLine(shapes map[rune][]string, text string) ([]string, error) {
	if text == "" {
		return nil, ErrEmptyText
	}
	rows := make([]string, 8)

	for _, char := range text {
		shape, err := RenderCharacter(shapes, char)
		if err != nil {
			return nil, err
		}

		for i := 0; i < 8; i++ {
			rows[i] += shape[i]
		}
	}

	return rows, nil
}

// Render renders complete text, including logical newlines, into ASCII-art rows.
func Render(shapes map[rune][]string, text string) ([]string, error) {
	if text == "" {
		return nil, ErrEmptyText
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	rows := []string{}
	for _, line := range lines {
		if line == "" {
			for i := 0; i < 8; i++ {
				rows = append(rows, "")
			}
		} else {
			renderedRows, err := RenderLine(shapes, line)
			if err != nil {
				return nil, err
			}
			rows = append(rows, renderedRows...)

		}

	}
	return rows, nil
}

// if we need a string edit the function to be Render(...) (string, error)  and add test = strings.Join(rows, "\n") and return text

func GenerateArt(text, bannerName string) (string, error) {
	bannerShapes, err := GetShapes(bannerName)
	if err != nil {
		return "", err
	}

	rows, err := Render(bannerShapes, text)
	if err != nil {
		return "", err
	}

	return strings.Join(rows, "\n"), nil
}

func SegmentRow(row string, coloredRanges []layout.Boundary) ([]Segment, error) {
	segments := []Segment{}
	currentIndex := 0

	for _, coloredRange := range coloredRanges {
		if coloredRange.Start < currentIndex ||
			coloredRange.Start >= coloredRange.End ||
			coloredRange.End > len(row) {
			return nil, ErrInvalidColoredRanges
		}

		// Preserve uncolored text before the match.
		if currentIndex < coloredRange.Start {
			segments = append(segments, Segment{
				Text:    row[currentIndex:coloredRange.Start],
				Colored: false,
			})
		}

		// Add the selected colored text.
		segments = append(segments, Segment{
			Text:    row[coloredRange.Start:coloredRange.End],
			Colored: true,
		})

		// Advance past the processed colored text.
		currentIndex = coloredRange.End
	}

	// Preserve text after the last match.
	if currentIndex < len(row) {
		segments = append(segments, Segment{
			Text:    row[currentIndex:],
			Colored: false,
		})
	}

	return segments, nil
}

func RenderSelected(shapes map[rune][]string, text string, target string) ([]Row, error) {

	// 1. Normalize newlines for all operations.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// 2. Render the original text.
	rendered, err := Render(shapes, text)
	if err != nil {
		return nil, err
	}

	// 3. Calculate character boundaries for every row.
	boundaries, err := layout.CalculateMultilineBoundaries(shapes, text)
	if err != nil {
		return nil, err
	}

	// 4. Separate logical lines.
	lines := strings.Split(text, "\n")

	result := []Row{}

	// 5. Process each logical line independently.
	for lineIndex, line := range lines {

		ranges := layout.FindRanges(line, target)
		mergedRanges := layout.MergeRanges(ranges)

		// 6. Process the eight rendered rows of this line.
		for rowIndex := 0; rowIndex < 8; rowIndex++ {

			actualRowIndex := lineIndex*8 + rowIndex
			renderedRow := rendered[actualRowIndex]

			rowBoundaries := boundaries[lineIndex][rowIndex]

			coloredRanges := []layout.Boundary{}

			// Convert character ranges to rendered columns.
			for _, match := range mergedRanges {
				converted, err := layout.CharacterRangeToRenderedRange(
					rowBoundaries,
					match.Start,
					match.End,
				)
				if err != nil {
					return nil, err
				}

				coloredRanges = append(coloredRanges, converted)
			}

			// 7. Split the rendered row into segments.
			segments, err := SegmentRow(renderedRow, coloredRanges)
			if err != nil {
				return nil, err
			}

			result = append(result, Row{
				Segments: segments,
			})
		}
	}

	// 8. Return all rendered rows.
	return result, nil
}

func RenderBlocks(
	shapes map[rune][]string,
	blocks []justify.TextBlock,
) ([]RenderedBlock, error) {
	rendered := make([]RenderedBlock, 0, len(blocks))

	for _, block := range blocks {
		var text strings.Builder
		target := ""

		for _, segment := range block.Segments {
			text.WriteString(segment.Text)

			if target == "" && segment.Colored {
				target = segment.Text
			}
		}

		blockText := text.String()

		// Render each space independently.
		if len(blockText) > 0 && strings.Trim(blockText, " ") == "" {
			for _, segment := range block.Segments {
				for range segment.Text {
					spaceTarget := ""
					if segment.Colored {
						spaceTarget = " "
					}

					rows, err := RenderSelected(shapes, " ", spaceTarget)
					if err != nil {
						return nil, err
					}

					rendered = append(rendered, RenderedBlock{
						Rows:    rows,
						Justify: block.Justify,
					})
				}
			}
			continue
		}

		rows, err := RenderSelected(shapes, blockText, target)
		if err != nil {
			return nil, err
		}

		rendered = append(rendered, RenderedBlock{
			Rows:    rows,
			Justify: block.Justify,
		})
	}

	return rendered, nil
}
