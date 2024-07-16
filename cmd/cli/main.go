package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	noteRepresentation = "◎"
	hitDuration        = 150 * time.Millisecond
)

var (
	don = lipgloss.NewStyle().Foreground(lipgloss.Color("#e94d2d")).Render(noteRepresentation)
	ka  = lipgloss.NewStyle().Foreground(lipgloss.Color("#63b8b5")).Render(noteRepresentation)
)

type song struct {
	bpm          int
	notes        []string
}

func randomSong(bpm, length int) *song {
	s := song{
		bpm:          bpm,
		notes: make([]string, 0),
	}

	for range length {
		randint := rand.IntN(10)

		switch {
		case randint < 2:
			s.notes = append(s.notes, don)
		case randint < 4:
			s.notes = append(s.notes, ka)
		default:
			s.notes = append(s.notes, " ")
		}
	}

	return &s
}

func (s *song) Update() {
	if len(s.notes) == 0 {
		s.notes = []string{""}
		return
	}

	s.notes = s.notes[1:]
}

func (s song) View(width int) string {
	if width <= len(s.notes) {
		return strings.Join(s.notes[:width], "")
	}

	result := make([]string, width)
	copy(result, s.notes)

	for range width-len(s.notes) {
		result = append(result, " ")
	}

	return strings.Join(result, "")
}

type boardModel struct {
	width int

	timeLeftKaHit   time.Time
	timeRightKaHit  time.Time
	timeLeftDonHit  time.Time
	timeRightDonHit time.Time

	isLeftKaHit   bool
	isRightKaHit  bool
	isLeftDonHit  bool
	isRightDonHit bool

	song *song
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() boardModel {
	return boardModel{
		width: 0,
		song:  randomSong(240, 30),
	}
}

func (m boardModel) Init() tea.Cmd {
	ticker := time.NewTicker(time.Duration(60_000_000_000/m.song.bpm) * time.Nanosecond)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				m.song.Update()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return doTick()
}

func (m boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.isRightKaHit = !(time.Until(m.timeRightKaHit) <= 0)
	m.isRightDonHit = !(time.Until(m.timeRightDonHit) <= 0)
	m.isLeftKaHit = !(time.Until(m.timeLeftKaHit) <= 0)
	m.isLeftDonHit = !(time.Until(m.timeLeftDonHit) <= 0)

	switch msg := msg.(type) {
	case tickMsg:
		return m, doTick()
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			m.isRightDonHit = true
			m.timeRightDonHit = time.Now().Add(hitDuration)
		case "k":
			m.isRightKaHit = true
			m.timeRightKaHit = time.Now().Add(hitDuration)
		case "f":
			m.isLeftDonHit = true
			m.timeLeftDonHit = time.Now().Add(hitDuration)
		case "d":
			m.isLeftKaHit = true
			m.timeLeftKaHit = time.Now().Add(hitDuration)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m boardModel) View() string {
	sb := &strings.Builder{}

	sb.WriteString("TaTaCLI - Taiko no Tatsujin Command Line Interface\n\n")

	if m.width != 0 {
		fmt.Fprintf(sb, "%s%s%s\n", "┏", strings.Repeat("━", m.width-2), "┓")

		// First line
		sb.WriteString("┃ ")

		if m.isLeftKaHit && m.isLeftDonHit {
			sb.WriteString("\x1b[34m▄\x1b[0;34;41m▀▀\x1b[0m")
		} else if m.isLeftKaHit {
			sb.WriteString("\x1b[34m▄▀▀\x1b[0m")
		} else if m.isLeftDonHit {
			sb.WriteString("▄\x1b[41m▀▀\x1b[0m")
		} else {
			sb.WriteString("▄▀▀")
		}
		if m.isRightKaHit && m.isRightDonHit {
			sb.WriteString("\x1b[0;34;41m▀▀\x1b[0;34;49m▄\x1b[0m")
		} else if m.isRightKaHit {
			sb.WriteString("\x1b[34m▀▀▄\x1b[0m")
		} else if m.isRightDonHit {
			sb.WriteString("\x1b[41m▀▀\x1b[0m▄")
		} else {
			sb.WriteString("▀▀▄")
		}
		sb.WriteString(" ┃  ╭───╮")
		sb.WriteString(strings.Repeat(" ", m.width-18))
		sb.WriteString("┃\n")

		// Second line
		sb.WriteString("┃ ")

		if m.isLeftKaHit {
			sb.WriteString("\x1b[34m█\x1b[0m")
		} else {
			sb.WriteString("█")
		}

		if m.isLeftDonHit {
			sb.WriteString("\x1b[31m██\x1b[0m")
		} else {
			sb.WriteString(strings.Repeat(" ", 2))
		}

		if m.isRightDonHit {
			sb.WriteString("\x1b[31m██\x1b[0m")
		} else {
			sb.WriteString(strings.Repeat(" ", 2))
		}

		if m.isRightKaHit {
			sb.WriteString("\x1b[34m█\x1b[0m")
		} else {
			sb.WriteString("█")
		}

		sb.WriteString(" ┃  │   │")
		sb.WriteString(m.song.View(m.width - 18))
		sb.WriteString("┃\n")

		// Third line
		sb.WriteString("┃ ")

		if m.isLeftKaHit && m.isLeftDonHit {
			sb.WriteString("\x1b[34m▀\x1b[0;34;41m▄▄\x1b[0m")
		} else if m.isLeftKaHit {
			sb.WriteString("\x1b[34m▀▄▄\x1b[0m")
		} else if m.isLeftDonHit {
			sb.WriteString("▀\x1b[41m▄▄\x1b[0m")
		} else {
			sb.WriteString("▀▄▄")
		}

		if m.isRightKaHit && m.isRightDonHit {
			sb.WriteString("\x1b[0;34;41m▄▄\x1b[0;34;49m▀\x1b[0m")
		} else if m.isRightKaHit {
			sb.WriteString("\x1b[34m▄▄▀\x1b[0m")
		} else if m.isRightDonHit {
			sb.WriteString("\x1b[41m▄▄\x1b[0m▀")
		} else {
			sb.WriteString("▄▄▀")
		}

		sb.WriteString(" ┃  ╰───╯")
		sb.WriteString(strings.Repeat(" ", m.width-18))
		sb.WriteString("┃\n")

		fmt.Fprintf(sb, "%s%s%s\n", "┗", strings.Repeat("━", m.width-2), "┛")
	}

	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithFPS(60))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
