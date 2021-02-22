package packet

//go:generate stringer -trimprefix=Phase -type=Phase

// Phase is a type for constant phases that a connection can be in.
type Phase uint8

// Known phases.
const (
	PhaseInvalid Phase = iota
	PhaseHandshaking
	PhaseStatus
	PhaseLogin
	PhasePlay
)
