package ui

import "strings"

// cursorMark is a zero-width placeholder put where the input cursor should go
// once a screen is printed. It is removed before printing.
const cursorMark = "\uE000"

// block is a rectangle of text lines that can be composed with other blocks.
type block []string

// width is the visible width of the widest line.
func (b block) width() int {
	w := 0
	for _, l := range b {
		w = max(w, visibleWidth(l))
	}
	return w
}

// centered returns the block with every line centred as a whole in width w,
// keeping the block's own left alignment intact.
func (b block) centered(w int) block {
	bw := b.width()
	left := strings.Repeat(" ", max(w-bw, 0)/2)
	out := make(block, len(b))
	for i, l := range b {
		out[i] = padTo(left+padTo(l, bw), w)
	}
	return out
}

// visibleWidth counts the characters of s that take up space on screen,
// skipping ANSI escape sequences and the cursor mark.
func visibleWidth(s string) int {
	n := 0
	inEscape := false
	for _, r := range s {
		switch {
		case inEscape:
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
		case r == '\033':
			inEscape = true
		case string(r) == cursorMark:
		default:
			n++
		}
	}
	return n
}

// padTo right-pads s with spaces to the visible width w.
func padTo(s string, w int) string {
	return s + strings.Repeat(" ", max(w-visibleWidth(s), 0))
}

// centerIn centres s within the visible width w.
func centerIn(s string, w int) string {
	free := max(w-visibleWidth(s), 0)
	return strings.Repeat(" ", free/2) + s + strings.Repeat(" ", free-free/2)
}

// hjoin puts blocks next to each other, separated by gap spaces and each
// vertically centred against the tallest one. Empty blocks are skipped.
func hjoin(gap int, blocks ...block) block {
	var parts []block
	height := 0
	for _, b := range blocks {
		if len(b) > 0 {
			parts = append(parts, b)
			height = max(height, len(b))
		}
	}
	out := make(block, height)
	sep := strings.Repeat(" ", gap)
	for i, b := range parts {
		bw := b.width()
		top := (height - len(b)) / 2
		for line := range height {
			s := ""
			if line >= top && line < top+len(b) {
				s = b[line-top]
			}
			if i > 0 {
				out[line] += sep
			}
			out[line] += padTo(s, bw)
		}
	}
	return out
}

// vjoin stacks blocks on top of each other, centred in a common width, with
// gap blank lines between them. Empty blocks are skipped.
func vjoin(gap int, blocks ...block) block {
	w := 0
	for _, b := range blocks {
		w = max(w, b.width())
	}
	var out block
	for _, b := range blocks {
		if len(b) == 0 {
			continue
		}
		if len(out) > 0 {
			for range gap {
				out = append(out, "")
			}
		}
		out = append(out, b.centered(w)...)
	}
	return out
}

// findCursor locates and removes the cursor mark, returning the cleaned lines
// and the mark's line and visible column. ok is false if there is no mark.
func findCursor(lines block) (out block, row, col int, ok bool) {
	out = make(block, len(lines))
	copy(out, lines)
	for i, l := range out {
		if idx := strings.Index(l, cursorMark); idx >= 0 {
			out[i] = l[:idx] + l[idx+len(cursorMark):]
			return out, i, visibleWidth(l[:idx]), true
		}
	}
	return out, 0, 0, false
}
