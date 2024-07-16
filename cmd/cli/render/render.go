package render

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	note    = "◎"
	bigNote = "◎"
	balloon = "⃝~"
	space   = "⃝ "
	drumrollExtension = "▪"
)

var (
	None              = []byte{}
	Pause             = []byte{' '}
	Don               = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#e94d2d")).Render(note))
	BigDon            = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#e94d2d")).Render(bigNote))
	Ka                = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#63b8b5")).Render(note))
	BigKa             = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#63b8b5")).Render(bigNote))
	Balloon           = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4504")).Render(balloon))
	DrumrollStart     = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb500")).Render(note))
	DrumrollExtension = []byte(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb500")).Render(drumrollExtension))
)
