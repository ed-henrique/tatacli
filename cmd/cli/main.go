package main

import (
	"tatacli/internal/input"
	"tatacli/internal/render"
)

func main() {
	board := render.NewBoard()

	input.Start(&board)
	render.Start(&board)
}
