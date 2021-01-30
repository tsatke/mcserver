package network

type Error string

func (e Error) Error() string { return string(e) }

const (
	// ErrClosed indicates, that the component is already
	// closed, and required resources might already have
	// been released.
	ErrClosed Error = "already closed"
)
