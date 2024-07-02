package note

const (
	noteText = "◉\x1b[0m"
)

type Size uint8

const (
	Small Size = iota
	Big
)

type Variant uint8

const (
	Any Variant = iota
	Don
	Ka
	Pause
)

type Note struct {
	Size           Size
	Variant        Variant
	Representation string
}

func New(size Size, variant Variant) Note {
	rep := noteText

	switch variant {
	case Don:
		rep = "\x1b[31m" + rep
	case Ka:
		rep = "\x1b[34m" + rep
	case Pause:
		rep = " "
	}

	return Note{
		Size:           size,
		Variant:        variant,
		Representation: rep,
	}
}

