package dkg

import (
	"time"
)

// Phaser must signal on its channel when the protocol should move to a next
// phase. Phase must be sequential: DealPhase (start), ResponsePhase,
// JustifPhase and then FinishPhase.
// Note that if the dkg protocol finishes before the phaser sends the
// FinishPhase, the protocol will not listen on the channel anymore. This can
// happen if there is no complaints, or if using the "FastSync" mode.
// Most of the times, user should use the TimePhaser when using the network, but
// if one wants to use a smart contract as a board, then the phaser can tick at
// certain blocks, or when the smart contract tells it.
type Phaser interface {
	NextPhase() chan Phase
}

// TimePhaser is a phaser that sleeps between the different phases and send the
// signal over its channel.
type TimePhaser struct {
	out   chan Phase
	sleep func(Phase)
}

func NewTimePhaser(p time.Duration) *TimePhaser {
	return NewTimePhaserFunc(func(Phase) { time.Sleep(p) })
}

func NewTimePhaserFunc(sleepPeriod func(Phase)) *TimePhaser {
	return &TimePhaser{
		out:   make(chan Phase, 4),
		sleep: sleepPeriod,
	}
}

func (t *TimePhaser) Start() {
	t.out <- DealPhase
	t.sleep(DealPhase)
	t.out <- ResponsePhase
	t.sleep(ResponsePhase)
	t.out <- JustifPhase
	t.sleep(JustifPhase)
	t.out <- FinishPhase
	close(t.out)
}

func (t *TimePhaser) NextPhase() chan Phase {
	return t.out
}

type copyPhaser struct {
	phaser   Phaser
	original chan Phase
	cop      chan Phase
}

// NewCopyPhaser returns a phaser that copies the phase returned from the
// orignal to the given channel. The copy channel must not be blocking otherwise
// it delays the original phaser's phase signals.
func NewCopyPhaser(original Phaser, cop chan Phase) Phaser {
	c := &copyPhaser{
		phaser:   original,
		original: make(chan Phase, 1),
		cop:      cop,
	}
	go c.copy()
	return c
}

func (c *copyPhaser) copy() {
	for phase := range c.phaser.NextPhase() {
		c.original <- phase
		c.cop <- phase
	}
	close(c.original)
	close(c.cop)
}

func (c *copyPhaser) NextPhase() chan Phase {
	return c.original
}
