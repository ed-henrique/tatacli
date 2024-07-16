package main

import (
	"fmt"
	"os"
	"tatacli/cmd/cli/models/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(game.New(), tea.WithFPS(60))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
