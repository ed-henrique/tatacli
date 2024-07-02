package render

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"tatacli/internal/note"
	"time"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

const (
	colorRed   = "\033[31m"
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"

	hitDuration = 100 * time.Millisecond
)

type Song struct {
	Counter int
	Notes   []note.Note
}

func (s *Song) Update() {
	s.Counter++
}

func (s Song) Render(width int) string {
	sb := &strings.Builder{}

	for i := s.Counter; i < width && i < len(s.Notes); i++ {
		sb.WriteString(s.Notes[i].Representation)
	}

	if s.Counter >= len(s.Notes) {
		sb.WriteString(strings.Repeat(" ", width))
	} else if s.Counter < len(s.Notes) && width > len(s.Notes[s.Counter:]) {
		sb.WriteString(strings.Repeat(" ", width-len(s.Notes)+s.Counter))
	}

	return sb.String()
}

type Board struct {
	timeLeftKaHit   time.Time
	timeRightKaHit  time.Time
	timeLeftDonHit  time.Time
	timeRightDonHit time.Time

	IsLeftKaHit   chan bool
	IsRightKaHit  chan bool
	IsLeftDonHit  chan bool
	IsRightDonHit chan bool
	CurrentBoard  string
	CurrentWidth  int
	CurrentHeight int
	Song
}

func NewBoard() Board {
	song := Song{}
	for i := 0; i < 300; i++ {
		rng := rand.Float64()

		if rng < 0.33 {
			song.Notes = append(song.Notes, note.New(note.Small, note.Don))
		} else if rng < 0.66 {
			song.Notes = append(song.Notes, note.New(note.Small, note.Ka))
		} else {
			song.Notes = append(song.Notes, note.New(note.Small, note.Pause))
		}
	}

	return Board{
		IsLeftKaHit:   make(chan bool),
		IsRightKaHit:  make(chan bool),
		IsLeftDonHit:  make(chan bool),
		IsRightDonHit: make(chan bool),
		Song:          song,
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
		b.timeLeftKaHit = time.Now().Add(hitDuration)
	default:
		isLeftKaHit = !(time.Until(b.timeLeftKaHit) <= 0)
	}

	select {
	case kaInput := <-b.IsRightKaHit:
		isRightKaHit = kaInput
		b.timeRightKaHit = time.Now().Add(hitDuration)
	default:
		isRightKaHit = !(time.Until(b.timeRightKaHit) <= 0)
	}

	select {
	case donInput := <-b.IsLeftDonHit:
		isLeftDonHit = donInput
		b.timeLeftDonHit = time.Now().Add(hitDuration)
	default:
		isLeftDonHit = !(time.Until(b.timeLeftDonHit) <= 0)
	}

	select {
	case donInput := <-b.IsRightDonHit:
		isRightDonHit = donInput
		b.timeRightDonHit = time.Now().Add(hitDuration)
	default:
		isRightDonHit = !(time.Until(b.timeRightDonHit) <= 0)
	}

	if isLeftKaHit && isLeftDonHit {
		sb.WriteString("\x1b[34m▄\x1b[0;34;41m▀▀\x1b[0m")
	} else if isLeftKaHit {
		sb.WriteString("\x1b[34m▄▀▀\x1b[0m")
	} else if isLeftDonHit {
		sb.WriteString("▄\x1b[41m▀▀\x1b[0m")
	} else {
		sb.WriteString("▄▀▀")
	}

	if isRightKaHit && isRightDonHit {
		sb.WriteString("\x1b[0;34;41m▀▀\x1b[0;34;49m▄\x1b[0m")
	} else if isRightKaHit {
		sb.WriteString("\x1b[34m▀▀▄\x1b[0m")
	} else if isRightDonHit {
		sb.WriteString("\x1b[41m▀▀\x1b[0m▄")
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
	sb.WriteString("┃  │")

	song := b.Song.Render(b.CurrentWidth - (b.CurrentWidth / 5) - 5)

	for i := 0; i < 4; i++ {
		if song[i] == ' ' && i == 3 {
			sb.WriteRune('│')
		} else {
			sb.WriteByte(song[i])
		}
	}

	sb.WriteString(song[5:])

	sb.WriteString("┃\n")

	// Third line
	sb.WriteString("┃")
	sb.WriteString(strings.Repeat(" ", b.CurrentWidth/5-12))

	if isLeftKaHit && isLeftDonHit {
		sb.WriteString("\x1b[34m▀\x1b[0;34;41m▄▄\x1b[0m")
	} else if isLeftKaHit {
		sb.WriteString("\x1b[34m▀▄▄\x1b[0m")
	} else if isLeftDonHit {
		sb.WriteString("▀\x1b[41m▄▄\x1b[0m")
	} else {
		sb.WriteString("▀▄▄")
	}

	if isRightKaHit && isRightDonHit {
		sb.WriteString("\x1b[0;34;41m▄▄\x1b[0;34;49m▀\x1b[0m")
	} else if isRightKaHit {
		sb.WriteString("\x1b[34m▄▄▀\x1b[0m")
	} else if isRightDonHit {
		sb.WriteString("\x1b[41m▄▄\x1b[0m▀")
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
	select {
	case <-time.After(hitDuration):
		b.Song.Update()
	}
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
		board.Render()
	}
}

func Finish() {
	// fmt.Println("\033[H\033[2JThanks for playing, see ya!")
	cursor.Show()
}
