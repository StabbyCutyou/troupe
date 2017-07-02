package troupe

import (
	"testing"
	"time"
)

// This is to cover an issue found with an older implementation of the shutdown logic
func xTestShutdown(t *testing.T) {
	a, _ := NewActor(ActorConfig{5})

	go func() {
		a.Accept(func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()
	a.stop()

	if !a.IsShutdown() {
		t.FailNow()
	}
}

func xTestIsFinished(t *testing.T) {
	a, _ := NewActor(ActorConfig{5})

	go func() {
		a.Accept(func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()
	a.stop()
	time.Sleep(200 * time.Millisecond)
	if !a.isFinished() {
		t.FailNow()
	}
}
