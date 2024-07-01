package render

import (
	"fmt"
	"strings"
	"time"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

const (
	colorRed   = "\033[31m"
	colorBlue  = "\033[34m"
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
	IsLeftKaHit   chan bool
	IsRightKaHit  chan bool
	IsLeftDonHit  chan bool
	IsRightDonHit chan bool
	CurrentBoard  string
	CurrentWidth  int
	CurrentHeight int
	Notes
}

func NewBoard() Board {
	return Board{
		IsLeftKaHit: make(chan bool),
		IsRightKaHit: make(chan bool),
		IsLeftDonHit: make(chan bool),
		IsRightDonHit: make(chan bool),
	}
}

func (b *Board) Update() {
	sb := &strings.Builder{}

	fmt.Fprintf(sb, "%s%s%s%s%s\n", "┏", strings.Repeat("━", b.CurrentWidth/5-1), "┳", strings.Repeat("━", b.CurrentWidth-(b.CurrentWidth/5)-2), "┓")
	fmt.Fprintf(sb, "%s%s%s%s%s\n", "┃", strings.Repeat(" ", b.CurrentWidth/5-1), "┃", strings.Repeat(" ", b.CurrentWidth-(b.CurrentWidth/5)-2), "┃")

	// This section updates 3 lines of the board
	// First line
	sb.WriteString("┃")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth/5-12))

	var isLeftKaHit, isRightKaHit, isLeftDonHit, isRightDonHit bool

	select {
	case kaInput := <-b.IsLeftKaHit:
		isLeftKaHit = kaInput
	default:
		isLeftKaHit = false
	}

	select {
	case kaInput := <-b.IsRightKaHit:
		isRightKaHit = kaInput
	default:
		isRightKaHit = false
	}

	select {
	case donInput := <-b.IsLeftDonHit:
		isLeftDonHit = donInput
	default:
		isLeftDonHit = false
	}

	select {
	case donInput := <-b.IsRightDonHit:
		isRightDonHit = donInput
	default:
		isRightDonHit = false
	}

	if isLeftKaHit {
		sb.WriteString("\x1b[34m▄▀▀\x1b[0m")
	} else {
		sb.WriteString("▄▀▀")
	}

	if isRightKaHit {
		sb.WriteString("\x1b[34m▀▀▄\x1b[0m")
	} else {
		sb.WriteString("▀▀▄")
	}

	sb.WriteString(strings.Repeat(" ", 5))
	sb.WriteString("┃  ╭───╮")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth-(b.CurrentWidth/5)-9))
	sb.WriteString("┃\n")

	// Second line
	sb.WriteString("┃")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth/5-12))

	if isLeftKaHit {
		sb.WriteString("\x1b[34m█\x1b[0m")
	} else {
		sb.WriteString("█")
	}

	if isLeftDonHit {
		sb.WriteString("\x1b[31m██\x1b[0m")
	} else {
		sb.WriteString(strings.Repeat(" ", 2))
	}

	if isRightDonHit {
		sb.WriteString("\x1b[31m██\x1b[0m")
	} else {
		sb.WriteString(strings.Repeat(" ", 2))
	}

	if isRightKaHit {
		sb.WriteString("\x1b[34m█\x1b[0m")
	} else {
		sb.WriteString("█")
	}

	sb.WriteString(strings.Repeat(" ", 5))
	sb.WriteString("┃  │   │")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth-(b.CurrentWidth/5)-9))
	sb.WriteString("┃\n")

	// Third line
	sb.WriteString("┃")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth/5-12))

	if isLeftKaHit {
		sb.WriteString("\x1b[34m▀▄▄\x1b[0m")
	} else {
		sb.WriteString("▀▄▄")
	}

	if isRightKaHit {
		sb.WriteString("\x1b[34m▄▄▀\x1b[0m")
	} else {
		sb.WriteString("▄▄▀")
	}

	sb.WriteString(strings.Repeat(" ", 5))
	sb.WriteString("┃  ╰───╯")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth-(b.CurrentWidth/5)-9))
	sb.WriteString("┃\n")

	fmt.Fprintf(sb, "%s%s%s%s%s\n", "┃", strings.Repeat(" ", b.CurrentWidth/5-1), "┃", strings.Repeat(" ", b.CurrentWidth-(b.CurrentWidth/5)-2), "┃")
	fmt.Fprintf(sb, "%s%s%s%s%s\n", "┗", strings.Repeat("━", b.CurrentWidth/5-1), "┻", strings.Repeat("━", b.CurrentWidth-(b.CurrentWidth/5)-2), "┛")

	// Update current board
	b.CurrentBoard = sb.String()
}

func (b Board) Render() {
	fmt.Printf("%s", b.CurrentBoard)
	cursor.StartOfLineUp(7)
}

func Start(board *Board) {
	if !term.IsTerminal(0) {
		panic("Not in a terminal!")
	}

	cursor.Hide()
	width, height, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	board.CurrentWidth = width
	board.CurrentHeight = height

	cursor.ClearLinesUp(height)
	cursor.StartOfLine()

	fmt.Print("TaTaCLI - Taiko no Tatsujin CLI\n")

	for {
		width, height, err = term.GetSize(0)
		if err != nil {
			panic(err)
		}

		board.CurrentWidth = width
		board.CurrentHeight = height

		board.Update()
		time.Sleep(50 * time.Millisecond)
		board.Render()
	}
}

func Finish() {
	fmt.Println("\033[H\033[2JThanks for playing, see ya!")
	cursor.Show()
}
