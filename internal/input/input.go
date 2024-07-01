package input

import (
	"os"
	"os/signal"
	"syscall"

	"tatacli/internal/render"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func Start(board *render.Board) {
	// Handle OS signals to finish the program
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		render.Finish()
		os.Exit(0)
	}()

	// Handle user input
	go func() {
		err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			switch key.Code {
			case keys.CtrlC, keys.Escape:
				return true, nil
			case keys.RuneKey:
				switch key.String() {
				case "d":
					board.IsLeftKaHit <- true
				case "f":
					board.IsLeftDonHit <- true
				case "j":
					board.IsRightDonHit <- true
				case "k":
					board.IsRightKaHit <- true
				case "q":
					return true, nil
				}
			}

			return false, nil
		})

		if err != nil {
			panic("Could not start keyboard input listener")
		}

		c <- os.Interrupt
	}()
}
