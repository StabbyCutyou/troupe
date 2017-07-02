package troupe

import (
	"sync/atomic"
	"time"
)

// These consts define the various states of the Actor
const (
	NOTBUSY int32 = iota
	BUSY
)

// Actor is an actor, who receives messages and acts over them
type Actor struct {
	errorHandler func(error)
	mailbox      chan Work
	quit         chan struct{}
	lastAccepted *int64
	busy         *int32
	lastFinished *int64
}

// ActorConfig is the configuration info needed to start a Actor
type ActorConfig struct {
	MailboxSize int
}

// NewActor returns a new Actor
func NewActor(c ActorConfig) (*Actor, error) {
	if c.MailboxSize < 0 {
		return nil, ActorConfigurationError("mailbox must be greater than or equal to 0")
	}
	a := &Actor{
		mailbox:      make(chan Work, c.MailboxSize),
		quit:         make(chan struct{}),
		busy:         new(int32),
		lastFinished: new(int64),
		lastAccepted: new(int64),
	}
	go a.loop()
	return a, nil
}

// Accept will push a Letter onto the Actors mailbox. If the box is full, this will
// block.
func (a *Actor) Accept(w Work) error {
	if a.IsShutdown() {
		return ActorShuttingDownError("actor is shutting down, cannot accept work")
	}
	select {
	case a.mailbox <- w:
		atomic.StoreInt64(a.lastAccepted, time.Now().Unix())
	default:
		return ActorFullError("actor is full, cannot accept work")
	}
	return nil
}

// LastAccepted returns the last time this Actor got a new letter. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (a *Actor) LastAccepted() int64 {
	return atomic.LoadInt64(a.lastAccepted)
}

// LastFinished returns the last time this Actor completed a job. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (a *Actor) LastFinished() int64 {
	return atomic.LoadInt64(a.lastFinished)

}

// IsBusy returns if the Actor is currently working or not. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (a *Actor) IsBusy() bool {
	if busy := atomic.LoadInt32(a.busy); busy == BUSY {
		return true
	}
	return false
}

// IsShutdown returns if the Actor is currently shutting down or not.
func (a *Actor) IsShutdown() bool {
	select {
	case <-a.quit:
		return true
	default:
		return false
	}
}

func (a *Actor) stop() {
	close(a.quit)
}

// isFinished is meant for internal use only, to be called only after shutdown is initiated so
// that the system knows when the actor has finished all of it's available work
func (a *Actor) isFinished() bool {
	if len(a.mailbox) == 0 && a.IsShutdown() {
		return true
	}
	return false
}

func (a *Actor) loop() {
	for {
		select {
		case w := <-a.mailbox:
			atomic.StoreInt32(a.busy, BUSY)
			if err := w(); err != nil && a.errorHandler != nil {
				a.errorHandler(err)
			}
			atomic.StoreInt64(a.lastFinished, time.Now().Unix())
			atomic.StoreInt32(a.busy, NOTBUSY)
		case <-a.quit:
			// there is a possibility that this goroutine picks up the quit signal
			// but something was in the middle of assigning work
			// so if len of agent isn't 0, we break to the top of the loop and try again
			if len(a.mailbox) > 0 {
				continue
			}
			return
		default:
			// TODO provide a means of configurable backoff
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (a *Actor) join() {
	for !a.isFinished() {
		time.Sleep(200 * time.Microsecond)
	}
}
