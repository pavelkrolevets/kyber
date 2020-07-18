package dkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPhaserFunc(t *testing.T) {
	var newPhase = make(chan Phase)
	sleep := func(p Phase) {
		newPhase <- p
	}
	phaser := NewTimePhaserFunc(sleep)
	go phaser.Start()
	require.Equal(t, DealPhase, drain(t, newPhase))
	require.Equal(t, ResponsePhase, drain(t, newPhase))
	require.Equal(t, JustifPhase, drain(t, newPhase))
	checkPhaserChannel(t, phaser.NextPhase())
}

func TestPhaserCopy(t *testing.T) {
	phaser := NewTimePhaserFunc(func(Phase) {})
	cop := make(chan Phase, 4)
	copyPhaser := NewCopyPhaser(phaser, cop)
	go phaser.Start()
	checkPhaserChannel(t, copyPhaser.NextPhase())
	checkPhaserChannel(t, cop)
}

func checkPhaserChannel(t *testing.T, phaser chan Phase) {
	require.Equal(t, DealPhase, drain(t, phaser))
	require.Equal(t, ResponsePhase, drain(t, phaser))
	require.Equal(t, JustifPhase, drain(t, phaser))
	require.Equal(t, FinishPhase, drain(t, phaser))
}

func drain(t *testing.T, d chan Phase) Phase {
	select {
	case phase := <-d:
		return phase
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout draining channel")
	}
	return 0
}
