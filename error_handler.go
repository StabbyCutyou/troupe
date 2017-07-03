package troupe

// ErrorHandler is a way to define a generic action to accept an error, and make a
// decision about what to do. The idea is, your Work function would return an error
// containing enough context that you can type-switch on it in an ErrorHandler, and
// take the appropriate action with the additional information on the typed-error
type ErrorHandler func(error)
