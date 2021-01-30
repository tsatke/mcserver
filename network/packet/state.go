package packet

//go:generate stringer -linecomment -type=State

type State uint8

const (
	StateInvalid     State = iota // Invalid
	StateHandshaking              // Handshaking
	StateStatus                   // Status
	StateLogin                    // Login
	StatePlay                     // Play
)
