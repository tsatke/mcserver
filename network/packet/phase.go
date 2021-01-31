package packet

//go:generate stringer -linecomment -trimprefix=Phase -type=Phase

type Phase uint8

const (
	PhaseInvalid Phase = iota
	PhaseHandshaking
	PhaseStatus
	PhaseLogin
	PhasePlay
)
