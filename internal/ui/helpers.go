package ui

// SetNumPlayers tells the UI how many seats the table has, so players without
// revealed cards are shown too.
func (ui *TextUI) SetNumPlayers(n int) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	ui.numPlayers = n
}
