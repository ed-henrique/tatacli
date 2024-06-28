package render

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

const (
	colorRed = "\033[31m"
	colorBlue = "\033[34m"
	colorReset = "\033[0m"
)

type Note string

func NewNote() Note {
	return Note("◉")
}

func Red() Note {
	return Note(fmt.Sprintf("%s◉%s", colorRed, colorReset))
}

func Blue() Note {
	return Note(fmt.Sprintf("%s◉%s", colorBlue, colorReset))
}

type Notes []Note

func (n Notes) Render(width int) string {
	return strings.Repeat(" ", width)
}

type Board struct {
	Notes
}

func NewBoard() Board {
	return Board{}
}

func (b Board) Render(height, width int) {
	fmt.Printf("%s%s%s%s%s\n", "┏", strings.Repeat("━", width/5-1), "┳", strings.Repeat("━", width-(width/5)-2), "┓")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃", strings.Repeat(" ", width-(width/5)-2), "┃")
	fmt.Printf("%s%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-14), "  ▄▀▀▀▀▄     ", "┃  ╭───╮", strings.Repeat(" ", width-(width/5)-9), "┃")
	fmt.Printf("%s%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-14), "  █    █     ", "┃  │   │", b.Notes.Render(width-(width/5)-9), "┃")
	fmt.Printf("%s%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-14), "  ▀▄▄▄▄▀     ", "┃  ╰───╯", strings.Repeat(" ", width-(width/5)-9), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃", strings.Repeat(" ", width-(width/5)-2), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┗", strings.Repeat("━", width/5-1), "┻", strings.Repeat("━", width-(width/5)-2), "┛")
	cursor.StartOfLineUp(7)
}

func Start() {
	if !term.IsTerminal(0) {
		panic("Not in a terminal!")
	}

	cursor.Hide()
	board := NewBoard()

	width, height, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	cursor.ClearLinesUp(height)
	cursor.StartOfLine()

	fmt.Print("TaTaCLI - Taiko no Tatsujin CLI\n")

	for {
		width, height, err = term.GetSize(0)
		if err != nil {
			panic(err)
		}

		board.Render(height, width)
	}
}

func Finish() {
	fmt.Println("\033[H\033[2JThanks for playing, see ya!")
	cursor.Show()
}
