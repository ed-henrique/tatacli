package game

import (
	"fmt"
	"strings"
	"tatacli/cmd/cli/models/song"
	"tatacli/cmd/cli/note"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const HitDuration = 150 * time.Millisecond

type Game struct {
	width int

	timeLeftKaHit   time.Time
	timeRightKaHit  time.Time
	timeLeftDonHit  time.Time
	timeRightDonHit time.Time

	isLeftKaHit   bool
	isRightKaHit  bool
	isLeftDonHit  bool
	isRightDonHit bool

	song *song.Song

	points              int
	lastSongNotesLength int
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func New() Game {
	b := Game{song: song.Random(240, 50)}
	b.lastSongNotesLength = len(b.song.Notes)
	return b
}

func (g Game) Init() tea.Cmd {
	ticker := time.NewTicker(time.Duration(60_000_000_000/g.song.BPM) * time.Nanosecond)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				g.song.Update()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return doTick()
}

func (g *Game) AddPoints(t note.Type) {
	if g.lastSongNotesLength != len(g.song.Notes) && g.song.Head() == t {
		g.points++
		g.lastSongNotesLength = len(g.song.Notes)
	}
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	g.isRightKaHit = !(time.Until(g.timeRightKaHit) <= 0)
	g.isRightDonHit = !(time.Until(g.timeRightDonHit) <= 0)
	g.isLeftKaHit = !(time.Until(g.timeLeftKaHit) <= 0)
	g.isLeftDonHit = !(time.Until(g.timeLeftDonHit) <= 0)

	switch msg := msg.(type) {
	case tickMsg:
		return g, doTick()
	case tea.WindowSizeMsg:
		g.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			g.isRightDonHit = true
			g.timeRightDonHit = time.Now().Add(HitDuration)
			g.AddPoints(note.Don)
		case "k":
			g.isRightKaHit = true
			g.timeRightKaHit = time.Now().Add(HitDuration)
			g.AddPoints(note.Ka)
		case "f":
			g.isLeftDonHit = true
			g.timeLeftDonHit = time.Now().Add(HitDuration)
			g.AddPoints(note.Don)
		case "d":
			g.isLeftKaHit = true
			g.timeLeftKaHit = time.Now().Add(HitDuration)
			g.AddPoints(note.Ka)
		case "ctrl+c", "q":
			return g, tea.Quit
		}
	}

	return g, nil
}

func (g Game) View() string {
	sb := &strings.Builder{}

	sb.WriteString("TaTaCLI - Taiko no Tatsujin Command Line Interface\n\n")

	if g.width != 0 {
		fmt.Fprintf(sb, "%s%s%s\n", "┏", strings.Repeat("━", g.width-2), "┓")

		// First line
		sb.WriteString("┃ ")

		if g.isLeftKaHit && g.isLeftDonHit {
			sb.WriteString("\x1b[34m▄\x1b[0;34;41m▀▀\x1b[0m")
		} else if g.isLeftKaHit {
			sb.WriteString("\x1b[34m▄▀▀\x1b[0m")
		} else if g.isLeftDonHit {
			sb.WriteString("▄\x1b[41m▀▀\x1b[0m")
		} else {
			sb.WriteString("▄▀▀")
		}
		if g.isRightKaHit && g.isRightDonHit {
			sb.WriteString("\x1b[0;34;41m▀▀\x1b[0;34;49m▄\x1b[0m")
		} else if g.isRightKaHit {
			sb.WriteString("\x1b[34m▀▀▄\x1b[0m")
		} else if g.isRightDonHit {
			sb.WriteString("\x1b[41m▀▀\x1b[0m▄")
		} else {
			sb.WriteString("▀▀▄")
		}
		sb.WriteString(" ┃  ╭───╮")
		sb.WriteString(strings.Repeat(" ", g.width-18))
		sb.WriteString("┃\n")

		// Second line
		sb.WriteString("┃ ")

		if g.isLeftKaHit {
			sb.WriteString("\x1b[34m█\x1b[0m")
		} else {
			sb.WriteString("█")
		}

		if g.isLeftDonHit {
			sb.WriteString("\x1b[31m██\x1b[0m")
		} else {
			sb.WriteString(strings.Repeat(" ", 2))
		}

		if g.isRightDonHit {
			sb.WriteString("\x1b[31m██\x1b[0m")
		} else {
			sb.WriteString(strings.Repeat(" ", 2))
		}

		if g.isRightKaHit {
			sb.WriteString("\x1b[34m█\x1b[0m")
		} else {
			sb.WriteString("█")
		}

		sb.WriteString(" ┃  │   │")
		sb.WriteString(g.song.View(g.width - 18))
		sb.WriteString("┃\n")

		// Third line
		sb.WriteString("┃ ")

		if g.isLeftKaHit && g.isLeftDonHit {
			sb.WriteString("\x1b[34m▀\x1b[0;34;41m▄▄\x1b[0m")
		} else if g.isLeftKaHit {
			sb.WriteString("\x1b[34m▀▄▄\x1b[0m")
		} else if g.isLeftDonHit {
			sb.WriteString("▀\x1b[41m▄▄\x1b[0m")
		} else {
			sb.WriteString("▀▄▄")
		}

		if g.isRightKaHit && g.isRightDonHit {
			sb.WriteString("\x1b[0;34;41m▄▄\x1b[0;34;49m▀\x1b[0m")
		} else if g.isRightKaHit {
			sb.WriteString("\x1b[34m▄▄▀\x1b[0m")
		} else if g.isRightDonHit {
			sb.WriteString("\x1b[41m▄▄\x1b[0m▀")
		} else {
			sb.WriteString("▄▄▀")
		}

		sb.WriteString(" ┃  ╰───╯")
		sb.WriteString(strings.Repeat(" ", g.width-18))
		sb.WriteString("┃\n")

		fmt.Fprintf(sb, "%s%s%s\n", "┗", strings.Repeat("━", g.width-2), "┛")

		// Points
		fmt.Fprintf(sb, "%012d\n", g.points)
	}

	return sb.String()
}
