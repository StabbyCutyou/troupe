package troupe

import (
	"sync/atomic"
	"time"
)

// Actor is an actor, who receives messages and acts over them
type Actor struct {
	errorHandler func(error)
	agent        chan Work
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
	if c.MailboxSize < 1 {
		return nil, MailboxSizeTooSmallError("mailbox must be greater than 0")
	}
	a := &Actor{
		agent:        make(chan Work, c.MailboxSize),
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
	a.agent <- w
	atomic.StoreInt64(a.lastAccepted, time.Now().Unix())
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
	if busy := atomic.LoadInt32(a.busy); busy == 1 {
		return true
	}
	return false
}

// IsShutdown returns if the Actor is currently shutting down or not. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
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

func (a *Actor) loop() {
	for {
		select {
		case w := <-a.agent:
			atomic.StoreInt32(a.busy, 1)
			if err := w(); err != nil && a.errorHandler != nil {
				a.errorHandler(err)
			}
			atomic.StoreInt64(a.lastFinished, time.Now().Unix())
			atomic.StoreInt32(a.busy, 1)
		case <-a.quit:
			return
		default:
			// TODO provide a means of configurable backoff
			time.Sleep(50 * time.Microsecond)
		}
	}
}
