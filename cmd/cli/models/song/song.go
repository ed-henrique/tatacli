package song

import (
	"bytes"
	"math/rand/v2"
	"tatacli/cmd/cli/note"
	"tatacli/cmd/cli/render"
)

type Song struct {
	BPM   int
	Notes []note.Note
}

func Random(bpm, length int) *Song {
	s := Song{
		BPM:   bpm,
		Notes: make([]note.Note, 0),
	}

	for range length {
		randint := rand.IntN(100)

		switch {
		case randint < 10:
			s.Notes = append(s.Notes, note.New(note.Don))
		case randint < 20:
			s.Notes = append(s.Notes, note.New(note.BigDon))
		case randint < 30:
			s.Notes = append(s.Notes, note.New(note.Ka))
		case randint < 40:
			s.Notes = append(s.Notes, note.New(note.BigKa))
		case randint < 50:
			s.Notes = append(s.Notes, note.New(note.Balloon))
		case randint < 60:
			s.Notes = append(s.Notes, note.New(note.DrumrollStart))
		case randint < 70:
			s.Notes = append(s.Notes, note.New(note.DrumrollExtension))
		default:
			s.Notes = append(s.Notes, note.New(note.Pause))
		}
	}

	return &s
}

func (s Song) Head() note.Type {
	if len(s.Notes) == 0 {
		return note.Pause
	}

	return s.Notes[0].Type
}

func (s *Song) Update() {
	if len(s.Notes) == 0 {
		s.Notes = []note.Note{}
		return
	}

	s.Notes = s.Notes[1:]
}

func (s Song) View(width int) string {
	if width <= len(s.Notes) {
		notes := s.Notes[:width]
		notesRep := make([][]byte, len(notes)) 

		for i, n := range notes {
			notesRep[i] = n.Representation
		}

		return string(bytes.Join(notesRep, []byte{}))
	}

	notesRep := make([][]byte, len(s.Notes)) 

	for i, n := range s.Notes {
		notesRep[i] = n.Representation
	}

	result := make([][]byte, width)
	copy(result, notesRep)

	for range width - len(s.Notes) {
		result = append(result, render.Pause)
	}

	return string(bytes.Join(result, []byte{}))
}
