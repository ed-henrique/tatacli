package input

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
	"tatacli/internal/render"
)

func Start() {
	// Handle OS signals to finish the program
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		render.Finish()
		os.Exit(0)
	}()

	// Handle user input
	cInput, err := keyboard.GetKeys(1)
	if err != nil {
		panic("Could not start input listener!")
	}

	go func() {
		var key keyboard.KeyEvent

		for {
			key = <-cInput
			switch key.Rune {
			case 'z':
			case 'x':
			case 'q':
				_ = keyboard.Close()
				c <- os.Interrupt
			}
		}
	}()
}
