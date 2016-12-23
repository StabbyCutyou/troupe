package troupe

// These are the errors that are returned from using a Bobbin

// MailboxSizeTooSmallError is returned when you've misconfigured mailboxes
type MailboxSizeTooSmallError string

// Error implements the error interface
func (e MailboxSizeTooSmallError) Error() string {
	return string(e)
}

// ActorShuttingDownError is returned when you've misconfigured mailboxes
type ActorShuttingDownError string

// Error implements the error interface
func (e ActorShuttingDownError) Error() string {
	return string(e)
}
