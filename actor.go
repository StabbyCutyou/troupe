package troupe

import (
	"errors"
	"sync"
	"time"
)

// Actor is an actor, who receives messages and acts over them
type Actor struct {
	errorHandler  func(error)
	agent         chan Work
	quit          chan struct{}
	shutdownMutex sync.RWMutex
	shutdown      bool
	acceptMutex   sync.RWMutex
	lastAccepted  time.Time
	busyMutex     sync.RWMutex
	busy          bool
	lastFinished  time.Time
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
	b := &Actor{
		agent: make(chan Work, c.MailboxSize),
		quit:  make(chan struct{}),
	}
	go b.loop()
	return b, nil
}

// Accept will push a Letter onto the Actors mailbox. If the box is full, this will
// block.
func (b *Actor) Accept(w Work) error {
	if b.IsShutdown() {
		return errors.New("Actor is shutting down, no longer accepting work")
	}
	b.agent <- w
	b.acceptMutex.Lock()
	b.lastAccepted = time.Now()
	b.acceptMutex.Unlock()
	return nil
}

// LastAccepted returns the last time this Actor got a new letter. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (b *Actor) LastAccepted() time.Time {
	b.acceptMutex.RLock()
	defer b.acceptMutex.RUnlock()
	return b.lastAccepted
}

// LastFinished returns the last time this Actor completed a job. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (b *Actor) LastFinished() time.Time {
	b.acceptMutex.RLock()
	defer b.acceptMutex.RUnlock()
	return b.lastAccepted
}

// IsBusy returns if the Actor is currently working or not. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (b *Actor) IsBusy() bool {
	b.busyMutex.RLock()
	defer b.busyMutex.RUnlock()
	return b.busy
}

// IsShutdown returns if the Actor is currently shutting down or not. Note that while this is
// protected  by a Mutex, by the time you take action on the result of this value,
// it may have changed by another concurrent operation.
func (b *Actor) IsShutdown() bool {
	b.shutdownMutex.RLock()
	defer b.shutdownMutex.RUnlock()
	return b.shutdown
}

func (b *Actor) stop() {
	close(b.quit)
}

func (b *Actor) loop() {
	for {
		select {
		case w := <-b.agent:
			b.busyMutex.Lock()
			b.busy = true
			b.busyMutex.Unlock()
			if err := w(); err != nil && b.errorHandler != nil {
				b.errorHandler(err)
			}
			b.busyMutex.Lock()
			b.lastFinished = time.Now()
			b.busy = false
			b.busyMutex.Unlock()
		case <-b.quit:
			b.shutdownMutex.Lock()
			b.shutdown = true
			b.shutdownMutex.Unlock()
			return
		default:
		}
	}
}
