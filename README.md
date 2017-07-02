# Troupe

Troupe provides an implementation of an Actor based concurrency system in Go

It exposes a few basic structures: Actors, which act over a mailbox of work, and Troupes, which are collections of actors for a given purpose. Troupes and Actors accept Work, which is a closure that takes no arguments, and only returns an error. Troupes also allow you to define a per-Troupe error handler, so that you can make a decision about errors that happen to occur from a given unit of Work.

Nearly all behaviors in Troupe are non-blocking, meaning that if you go to Assign work to an available Actor, and none are currently able to accept work (because their mailboxes are full, for example), you will receive an error and will be required to re-submit the job based on your own heuristics for delivery (ie: you can keep trying until it finally is accepted, or define a retry-limit after which you no longer care to do the work.)

## Installation

`go get github.com/StabbyCutyou/troupe`

Troupe requires no additional dependencies outside of the Go standard library

## Creating a Troupe

The first step is to determine your configuration options, provided by `troupe.Config{}`

```golang
cfg := Config{
    Mode: Dynamic, // Troupes can operate as a Fixed pool of actors, or Dynamic with resisizng
    MailboxSize: 0, // How big the buffer of work an Actor can hold onto before it is "full"
    Min: 0, // The minimum number of Actors in a pool. Ingored for Fixed Troupes
    Initial: 0, // The initial number of Actors to pre-boot. Ignored for Fixed Troupes
    Max: 1000 // The maximum number of Actors for the Troupe. This is what is used for Fixed Troupes
}
```