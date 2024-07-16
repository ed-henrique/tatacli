package main

import (
	"fmt"
	"testing"
)

func TestSong(t *testing.T) {
	t.Run("Single don", func(t *testing.T) {
		s := song{notes: [][]byte{don}}

		totalLen := 0
		for _, ss := range s.notes {
			totalLen += len(ss)
		}

		expectedLen := len(don)

		if totalLen != expectedLen {
			t.Errorf("expected %d got %d", expectedLen, totalLen)
		}
	})

	t.Run("Single ka", func(t *testing.T) {
		s := song{notes: [][]byte{ka}}

		totalLen := 0
		for _, ss := range s.notes {
			totalLen += len(ss)
		}

		expectedLen := len(ka)

		if totalLen != expectedLen {
			t.Errorf("expected %d got %d", expectedLen, totalLen)
		}
	})

	t.Run("Sample song", func(t *testing.T) {
		s := song{notes: [][]byte{don, don, don, don, don, ka, ka, ka, ka, ka}}

		totalLen := 0
		for _, ss := range s.notes {
			totalLen += len(ss)
		}

		expectedLen := 5*len(don) + 5*len(ka)

		if totalLen != expectedLen {
			t.Errorf("expected %d got %d", expectedLen, totalLen)
		}
	})

	t.Run("Sample song view", func(t *testing.T) {
		s := song{notes: [][]byte{don, don, don, don, don, ka, ka, ka, ka, ka}}
		view := s.View(5)

		expectedLen := 5 * len(don)

		if len(view) != expectedLen {
			fmt.Println(view)
			t.Errorf("expected %d got %d", expectedLen, len(view))
		}
	})
}
