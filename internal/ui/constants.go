package ui

// seat is an edge of the table, as seen from above.
type seat int

// Layout settings
const (
	cardInnerWidth = 15
	cardGap        = "  "
)

// Table layout settings
const (
	localPlayer = 0 // the player sitting at the bottom edge of the table
	sideColumns = 2 // card columns for players at the left/right edges
	handColumns = 4 // card columns for the hand shown in the middle
	tableGap    = 4 // spaces between the side drafts and the middle
)

// Formatting and colors
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBlue   = "\033[34m"
)

const clearScreen = "\033[H\033[2J"

const (
	seatBottom seat = iota
	seatLeft
	seatTop
	seatRight
)
