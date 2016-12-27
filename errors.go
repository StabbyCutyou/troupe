package troupe

// These are the errors that are returned from using a Bobbin

// ConfigurationError is returned when you've misconfigured an actor.
// inspect the message for the specific reason
type ConfigurationError string

// Error implements the error interface
func (e ConfigurationError) Error() string {
	return string(e)
}

// ShuttingDownError is returned when you go to assign work, but the Troupe has begun to shut down
type ShuttingDownError string

// Error implements the error interface
func (e ShuttingDownError) Error() string {
	return string(e)
}

// ActorConfigurationError is returned when you've misconfigured an actor.
// inspect the message for the specific reason
type ActorConfigurationError string

// Error implements the error interface
func (e ActorConfigurationError) Error() string {
	return string(e)
}

// ActorShuttingDownError is returned when a Troupe has attempted to assign work to an Actor, but it was scheduled for shutdown
type ActorShuttingDownError string

// Error implements the error interface
func (e ActorShuttingDownError) Error() string {
	return string(e)
}

// ActorFullError is returned when an Actor is too full to accept work
type ActorFullError string

// Error implements the error interface
func (e ActorFullError) Error() string {
	return string(e)
}
