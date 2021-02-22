package game

type Sentinel string

func (s Sentinel) Error() string { return string(s) }

const (
	// ErrPlayerNotExist indicates that the player that wanted to join does not exist and
	// needs to be created.
	ErrPlayerNotExist Sentinel = "player does not exist"
)
