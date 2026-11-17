package ui

import (
	"fmt"
	"strconv"

	"github.com/BlackestDawn/ltu-boomerang/internal/engine"
)

type TextUI struct {
}

func (ui *TextUI) OnEvent(e engine.Event) {}

func (ui *TextUI) ShowResults(result engine.Result) {}

func (ui *TextUI) AskForNumber(max, min int) int {
	var okInput bool = false
	var choice int
	for !okInput {
		fmt.Printf("Choose a number between %d and %d: ", min, max)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Could not properly read input:", err)
			continue
		}
		choice, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Could not convert input:", err)
			continue
		}
		if choice < min || choice > max {
			fmt.Println("Input outisde of range:", choice)
			continue
		}
		okInput = true
	}
	return choice
}
