package justify

import "errors"

type LineLayout struct {
	Words    []string
	Leading  int
	Gaps     []int
	Trailing int
}

type Segment struct {
	Text    string
	Colored bool
}

type TextBlock struct {
	Segments []Segment
	Justify  bool
}

func ParseLine(text string) LineLayout {
	layout := LineLayout{}
	index := 0

	for index < len(text) {
		if text[index] == ' ' {
			start := index

			for index < len(text) && text[index] == ' ' {
				index++
			}

			count := index - start

			switch {
			case len(layout.Words) == 0:
				layout.Leading = count
			case index == len(text):
				layout.Trailing = count
			default:
				layout.Gaps = append(layout.Gaps, count)
			}

			continue
		}

		start := index

		for index < len(text) && text[index] != ' ' {
			index++
		}

		layout.Words = append(layout.Words, text[start:index])
	}

	return layout
}

func MarkColored(text, target string) []bool {
	colored := make([]bool, len(text))
	targetLen := len(target)

	for i := 0; i <= len(text)-targetLen; {
		if targetLen == 0 {
			return colored
		}
		if text[i:i+targetLen] == target {
			for j := 0; j < targetLen; j++ {
				colored[i+j] = true
			}
			i = i + targetLen
			continue
		}
		i++

	}

	return colored
}

func SplitTextBlocks(text string, colored []bool) ([]TextBlock, error) {
	blocks := []TextBlock{}
	if len(text) != len(colored) {
		return blocks, errors.New("text and colored must have the same length")
	}
	start := 0

	for i := 1; i < len(text); i++ {
		currentIsSpace := text[i] == ' '
		previousIsSpace := text[i-1] == ' '

		if currentIsSpace != previousIsSpace {
			isInternalSpace := start > 0 && text[start] == ' '
			segments := SplitSegments(text[start:i], colored[start:i])
			blocks = append(blocks, TextBlock{
				Segments: segments,
				Justify:  isInternalSpace,
			})

			start = i
		}
	}
	if start < len(text) {
		segments := SplitSegments(text[start:], colored[start:])
		blocks = append(blocks, TextBlock{
			Segments: segments,
			Justify:  false,
		})
	}
	return blocks, nil
}

func SplitSegments(text string, colored []bool) []Segment {
	segments := []Segment{}

	if len(text) == 0 {
		return segments
	}

	segmentStart := 0

	for i := 1; i < len(text); i++ {
		if colored[i] != colored[i-1] {
			segments = append(segments, Segment{
				Text:    text[segmentStart:i],
				Colored: colored[segmentStart],
			})

			segmentStart = i
		}
	}

	if segmentStart < len(text) {
		segments = append(segments, Segment{
			Text:    text[segmentStart:],
			Colored: colored[segmentStart],
		})
	}

	return segments
}
