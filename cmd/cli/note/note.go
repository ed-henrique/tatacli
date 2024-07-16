package note

import "tatacli/cmd/cli/render"

type Type uint8

const (
	Pause Type = iota
	Don
	BigDon
	Ka
	BigKa
	Balloon
	DrumrollStart
	DrumrollExtension
)

type Note struct {
	Type           Type
	Representation []byte
}

func New(t Type) Note {
	n := Note{Type: t}

	switch t {
	case Pause:
		n.Representation = render.Pause
	case Don:
		n.Representation = render.Don
	case BigDon:
		n.Representation = render.BigDon
	case Ka:
		n.Representation = render.Ka
	case BigKa:
		n.Representation = render.BigKa
	case Balloon:
		n.Representation = render.Balloon
	case DrumrollStart:
		n.Representation = render.DrumrollStart
	case DrumrollExtension:
		n.Representation = render.DrumrollExtension
	}

	return n
}

