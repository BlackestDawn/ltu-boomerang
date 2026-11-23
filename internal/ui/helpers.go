package ui

// SetNumPlayers tells the UI how many seats the table has, so players without
// revealed cards are shown too.
func (ui *TextUI) SetNumPlayers(n int) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	ui.numPlayers = n
}

// SetSelf tells the UI which seat the local player has. It defaults to 0,
// which is the host's seat.
func (ui *TextUI) SetSelf(seat int) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	ui.self = seat
}
