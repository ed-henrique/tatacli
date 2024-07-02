package board

import "time"

type Board struct {
	timeLeftKaHit   time.Time
	timeRightKaHit  time.Time
	timeLeftDonHit  time.Time
	timeRightDonHit time.Time

	IsLeftKaHit   chan bool
	IsRightKaHit  chan bool
	IsLeftDonHit  chan bool
	IsRightDonHit chan bool

	CurrentBoard  string
	CurrentWidth  int
	CurrentHeight int
	Notes
}
