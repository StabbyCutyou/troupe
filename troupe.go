// Package troupe defines a Troupe - a collection of Actors. Actors are actors
// who receive work on a channel. Troupes manage swarms of Actors, spinning them
// up and tearing them down depending on load.
package troupe

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Troupe represents a swarm of Actors
type Troupe struct {
	minActors          int
	maxActors          int
	ActorMutex         sync.Mutex
	Actors             []*Actor
	shutdown           bool
	defaultActorConfig ActorConfig
}

// Config is
type Config struct {
	Min              int
	Max              int
	Initial          int
	IdleActorTimeout time.Duration
	MailboxSize      int
}

// ActorConfig maps the Troupe Config struct into a ActorConfig
func (c Config) ActorConfig() ActorConfig {
	return ActorConfig{
		MailboxSize: c.MailboxSize,
	}
}

// NewTroupe returns a new Troupe
func NewTroupe(cfg Config) (*Troupe, error) {
	if cfg.Max < cfg.Initial {
		return nil, fmt.Errorf("Cannot create Troupe with Max (%d) < Inital (%d) size", cfg.Max, cfg.Initial)
	}
	if cfg.Min > cfg.Max {
		return nil, fmt.Errorf("Cannot create Troupe with Min (%d) > Max (%d) size", cfg.Min, cfg.Max)
	}
	if cfg.Max == 0 {
		return nil, fmt.Errorf("Max must be greater than 0")
	}
	if cfg.MailboxSize == 0 {
		cfg.MailboxSize = 1
	}
	bCfg := cfg.ActorConfig()
	Actors := make([]*Actor, 0)
	var err error
	var b *Actor
	for i := 0; i < cfg.Initial; i++ {
		if b, err = NewActor(bCfg); err != nil {
			return nil, err
		}
		Actors = append(Actors, b)
	}
	return &Troupe{
		Actors:             Actors,
		minActors:          cfg.Min,
		maxActors:          cfg.Max,
		defaultActorConfig: bCfg,
	}, nil
}

// Shutdown shuts down the Troupe
func (t *Troupe) Shutdown() {
	t.ActorMutex.Lock()
	t.shutdown = true
	for _, a := range t.Actors {
		a.stop()
	}
	t.ActorMutex.Unlock()
}

// Assign will distribute a Letter to the first available Actor. If there are no available Actors (that is, no Actors
// currently free from work, it grow the pool of Actors by 1. If the pool is already full, it will assign the work to the
// Actor who has least-recently been assigned work. If all Actors have their mailboxes full, Assign will block until
// one becomes free.
func (t *Troupe) Assign(w Work) error {
	t.ActorMutex.Lock()
	defer t.ActorMutex.Unlock()
	if t.shutdown {
		return errors.New("Unable to assign work - shutting down")
	}
	var item *Actor
	// First, do a best-effort attempt to find any Actors currently doing 0 work
	// This will make sure the assigned work is handled more quickly.
	for i, a := range t.Actors {
		if !a.IsBusy() {
			// Remove it, assign to it, push it back on the stack
			item = t.Actors[i]
			item.Accept(w)
			// Remove the item from where we found it
			t.Actors = append(t.Actors[:i], t.Actors[i+1:]...)
			// Reinsert it at the end
			t.Actors = append(t.Actors, item)
			return nil
		}
	}
	// We couldn't find one that wasn't busy
	// If the list is not full, make a new one
	if len(t.Actors) < t.maxActors {
		var err error
		item, err = NewActor(t.defaultActorConfig)
		if err != nil {
			return err
		}
		item.Accept(w)
		t.Actors = append(t.Actors, item)
		return nil
	}

	// If theres only 1, nothing to rotate on the list
	if len(t.Actors) == 1 {
		t.Actors[0].Accept(w)
		return nil
	}
	// The list was already at capacity, take the first one which we must assume is the
	// oldest waiting Actor, and assign to it. There are atleast 2 items on the list.
	item, t.Actors = t.Actors[0], t.Actors[1:]
	item.Accept(w)
	t.Actors = append(t.Actors, item)
	return nil
}
