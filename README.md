# Troupe [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/StabbyCutyou/troupe) [![Build Status](https://api.travis-ci.org/StabbyCutyou/troupe.svg)](https://travis-ci.org/StabbyCutyou/troupe)

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
    ErrorHandler: f(error), // A function which takes an error from a Work item, and determines how to handle it
    MailboxSize: 0, // How big the buffer of work an Actor can hold onto before it is "full"
    Min: 0, // The minimum number of Actors in a pool. Ingored for Fixed Troupes
    Initial: 0, // The initial number of Actors to pre-boot. Ignored for Fixed Troupes
    Max: 1000 // The maximum number of Actors for the Troupe. This is what is used for Fixed Troupes
}
```

Then, create a new Troupe

```golang
t, err := NewTroupe(c.cfg)
if err != nil {
    return err
}
```

Now, you can pass the Troupe around and have work assigned to it from any number of concurrent go routines.

## Assigning work

Once you have a troupe, you can assign work to it via the Assign method. Assigning work is done by passing in a closure, which contains all the state needed to perform the work.

```golang
w := func() error {
    // returns nil, or error
    return closedOverService.DoThing(closedOverDataA, closedOverDataB)
}

err := t.Assign(w)
// These are the three types of errors that an Assign call could respond with
switch x := err.(type):
case ShuttingDownError:
    // Means that you attempted to assign work during shutdown
case ActorShuttingDownError:
    // Should not be possible to be returned, but means you assigned work 
    // while shutting down but was caught by the actor
case ActorFullError:
    // Means that the actors mailbox was full, and you can try again as much
    // as you like, or consider the message a lost cause after a number of failed
    // attempts to assign
case nil:
    // No error, proceed
```

## Shutting down

Once your application is ready to terminate, simply call Shutdown, like so

```golang
t.Shutdown()
```

This will begin signalling all Actors that they should no longer accept work, which 
will prevent new work from being assigned.

It's important to note that this will not pre-empty any inflight work, which could still
be going on.

If you need to ensure current work finishes, you can use the Join method, to block
until all work has been completed.

```golang
t.Join()
```

If you don't care if the inflight work is finished, simply calling Shutdown is enough
to safely terminate, for your value of safety.

## Retry

Troupe has no built-in method of retry. It relies on you to define a way via the ErrorHandler to provide enough context to know when you need to re-assign a job, and how to do so. You should return a custom error that has enough context about the job being performed that the ErrorHandler can take appropriate action.

When shutting down, if you're feeding messages into Troupe from a broker that does not automatically enable retry after some timeout, you'll need to track the messages in flight, and signal the message broker that those messages should be eligible for retry.

## Example Implementation

Check out the `test/rpc/{client,server}` packages to see a basic implementation and a live test of the concept.

It includes a randomized error simulation, to demonstrate how the error handler could work.
