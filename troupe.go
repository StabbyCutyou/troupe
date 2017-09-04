// Package troupe defines a Troupe - a collection of Actors. Actors are actors
// who receive work on a channel. Troupes manage swarms of Actors, spinning them
// up and tearing them down depending on load.
package troupe

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Mode determines if the troupe uses a fixed size, or a dynamic size
// fixed size will use a blind random assignment of work strategy and a fixed number
// of resources, while a dynamic size will use a priority based assignment strategy
// which can grow the list of workers as needed.
type Mode int

const (
	// Dynamic is for the priority assignment
	Dynamic Mode = iota
	// Fixed is for the random assignment
	Fixed
)

// Troupe represents a swarm of Actors
type Troupe struct {
	minActors          int
	maxActors          int
	ActorMutex         sync.Mutex
	Actors             []*Actor
	shutdown           bool
	defaultActorConfig ActorConfig
	r                  *rand.Rand
	mode               Mode
}

// Config is
type Config struct {
	Min              int
	Max              int
	Initial          int
	IdleActorTimeout time.Duration
	MailboxSize      int
	Mode             Mode
	ErrorHandler     ErrorHandler
}

// ActorConfig maps the Troupe Config struct into a ActorConfig
func (c Config) ActorConfig() ActorConfig {
	return ActorConfig{
		MailboxSize:  c.MailboxSize,
		ErrorHandler: c.ErrorHandler,
	}
}

// NewTroupe returns a new Troupe
func NewTroupe(cfg Config) (*Troupe, error) {
	if cfg.Max < cfg.Initial {
		return nil, ConfigurationError(fmt.Sprintf("cannot create Troupe with Max (%d) < Inital (%d) size", cfg.Max, cfg.Initial))
	}
	if cfg.Min > cfg.Max {
		return nil, ConfigurationError(fmt.Sprintf("cannot create Troupe with Min (%d) > Max (%d) size", cfg.Min, cfg.Max))
	}
	if cfg.Max == 0 {
		return nil, ConfigurationError(fmt.Sprintf("max must be greater than 0"))
	}
	if cfg.MailboxSize == 0 {
		cfg.MailboxSize = 1
	}
	// For fixed mode, we need to allocate a fixed pool since the assignment
	// will not attempt to grow or shrink the pool
	if cfg.Mode == Fixed {
		cfg.Min = cfg.Max
		cfg.Initial = cfg.Max
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
		r:                  rand.New(rand.NewSource(time.Now().UnixNano())),
		mode:               cfg.Mode,
	}, nil
}

// Shutdown shuts down the Troupe
func (t *Troupe) Shutdown() error {
	t.ActorMutex.Lock()
	defer t.ActorMutex.Unlock()
	if t.shutdown {
		return ShuttingDownError("you cannot call shutdown more than once on a troupe")
	}
	t.shutdown = true
	// Stop all of them right away, to shut off their ability to accept work
	// Is this necessary to break into 2 steps?
	for _, a := range t.Actors {
		a.stop()
	}
	return nil
}

// Join is something i'm experimenting with to call after Shutdown, so you know when all work has ceased
func (t *Troupe) Join() {
	t.ActorMutex.Lock()
	// wait until they all finish their backlogs of work.
	for _, a := range t.Actors {
		a.join()
	}
	t.ActorMutex.Unlock()
}

// Assign will distribute a Work object to the nearest available actor as defined by the
// Troupes configuration: Either a Priority assignment approach which allows the pool
// to grow and shrink dynamically, or a Random assignment, which keeps the pool at a fixed size.
// Both have tradeoffs: The priority one is able to resize and be throttled dynamically, while
// the random one is overall faster for performance at the cost of utilizing a fixed cost
// of resources.
func (t *Troupe) Assign(w Work) error {
	if t.mode == Dynamic {
		return t.assignPriority(w)
	}
	return t.assignRand(w)
}

// assignPriority will distribute a Letter to the first available Actor. If there are no available Actors (that is, no Actors
// currently free from work, it grow the pool of Actors by 1. If the pool is already full, it will assign the work to the
// Actor who has least-recently been assigned work. If all Actors have their mailboxes full, Assign will block until
// one becomes free.
func (t *Troupe) assignPriority(w Work) error {
	t.ActorMutex.Lock()
	defer t.ActorMutex.Unlock()
	if t.shutdown {
		return ShuttingDownError("unable to assign work, shutting down")
	}
	var item *Actor
	// First, do a best-effort attempt to find any Actors currently doing 0 work
	// This will make sure the assigned work is handled more quickly.
	for i, a := range t.Actors {
		if !a.IsBusy() {
			// Remove it, assign to it, push it back on the stack
			item = t.Actors[i]
			if err := item.Accept(w); err != nil {
				return err
			}
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
		if err = item.Accept(w); err != nil {
			return err
		}
		t.Actors = append(t.Actors, item)
		return nil
	}
	// If theres only 1, nothing to rotate on the list
	if len(t.Actors) == 1 {
		return t.Actors[0].Accept(w)
	}
	// The list was already at capacity, take the first one which we must assume is the
	// oldest waiting Actor, and assign to it. There are atleast 2 items on the list.
	item, t.Actors = t.Actors[0], t.Actors[1:]
	if err := item.Accept(w); err != nil {
		return err
	}
	t.Actors = append(t.Actors, item)
	return nil
}

// assignRand skips priority, and attempts to assign work randomly
// I've tested crypto rand for a better random distribution, but in all cases it was worse
// than either the priority assign, or the raw random assign. It's worse than priority
// assign for smaller sized pools, it's worse than random assign for larger sized pools
// however it's so poor in general that it would not make a good middle ground option.
func (t *Troupe) assignRand(w Work) error {
	// Rand isn't threadsafe, womp womp
	t.ActorMutex.Lock()
	defer t.ActorMutex.Unlock()
	return t.Actors[t.r.Int()%(len(t.Actors))].Accept(w)
}
