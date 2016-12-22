package troupe

// Work is a closures that invokes whatever code you define.
// While you can pass any function directly that matches the signature, you're intended
// to pass closures over anything more complicated. This is to avoid relying on interface{}
type Work func() error
