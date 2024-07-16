package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	noteRepresentation = "◎"
)

var (
	don = lipgloss.NewStyle().Foreground(lipgloss.Color("#e94d2d")).SetString(noteRepresentation).String()
	ka  = lipgloss.NewStyle().Foreground(lipgloss.Color("#63b8b5")).SetString(noteRepresentation).String()
)

func main() {
	for _, r := range don {
		fmt.Printf("%c ", r)
	}
	fmt.Printf("%s%d\n", don, len(don))

	for _, r := range ka {
		fmt.Printf("%c ", r)
	}
	fmt.Printf("%s%d\n", ka, len(ka))

	fmt.Printf("%d\n", len([]rune{' '}))
}
