package ui

import "github.com/BlackestDawn/ltu-boomerang/internal/engine"

type TextUI struct {
}

func (ui *TextUI) OnEvent(e engine.Event) {}

func (ui *TextUI) ShowResults(result engine.Result) {}
