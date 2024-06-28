package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

const colorRed = "\033[31m"
const colorBlue = "\033[34m"
const colorReset = "\033[0m"

type note string

func newNote() note {
	return note("◉")
}

func red() note {
	return note(fmt.Sprintf("%s◉%s", colorRed, colorReset))
}

func blue() note {
	return note(fmt.Sprintf("%s◉%s", colorBlue, colorReset))
}

type notes []note

func (n notes) render(width int) string {
	return strings.Repeat(" ", width)
}

type board struct {
	notes
}

func newBoard() board {
	return board{}
}

func (b board) render(width int) {
	fmt.Printf("%s%s%s%s%s\n", "┏", strings.Repeat("━", width/5-1), "┳", strings.Repeat("━", width-(width/5)-2), "┓")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃", strings.Repeat(" ", width-(width/5)-2), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃ ╭───╮", strings.Repeat(" ", width-(width/5)-8), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃ │   │", b.notes.render(width-(width/5)-8), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃ ╰───╯", strings.Repeat(" ", width-(width/5)-8), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┃", strings.Repeat(" ", width/5-1), "┃", strings.Repeat(" ", width-(width/5)-2), "┃")
	fmt.Printf("%s%s%s%s%s\n", "┗", strings.Repeat("━", width/5-1), "┻", strings.Repeat("━", width-(width/5)-2), "┛")
	cursor.StartOfLineUp(7)
}

func main() {
	if !term.IsTerminal(0) {
		panic("Not in a terminal!")
	}

	cursor.Hide()
	c := make(chan os.Signal, 1)
	cHeight := make(chan int, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cursor.Show()
		cursor.UpAndClear(<-cHeight)
		cursor.StartOfLine()
		fmt.Println("Thanks for playing, see ya!")
		os.Exit(0)
	}()

	board := newBoard()
	// smallNote := newNote()
	// smallNote = red()

	width, height, err := term.GetSize(0)
	cursor.UpAndClear(height)
	fmt.Print("TaTaCLI - Taiko no Tatsujin CLI\n")

	for {
		width, height, err = term.GetSize(0)
		cHeight <- height
		if err != nil {
			panic(err)
		}

		board.render(width)
	}
}
