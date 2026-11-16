package justify

type LineLayout struct {
	Words    []string
	Leading  int
	Gaps     []int
	Trailing int
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
