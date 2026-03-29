package notes

// Lines reserved above/below the list scroll area (title + paging banner + spacing + nav hint).
const ListViewFrameLines = 5

// DetailViewFrameLines: title + spacing + nav hint below the viewport.
const DetailViewFrameLines = 5

const DefaultListTake = 2

const MaxListTake = 100

func DetailScrollViewportHeight(termHeight int) int {
	h := termHeight - DetailViewFrameLines
	if h < 4 {
		h = 4
	}
	return h
}
