package game

type Sentinel string

func (s Sentinel) Error() string { return string(s) }

const (
	// ErrChunkNotGenerated indicates that the requested chunk has not been generated yet.
	ErrChunkNotGenerated Sentinel = "chunk not generated"
	// ErrPlayerNotExist indicates that the player that wanted to join does not exist and
	// needs to be created.
	ErrPlayerNotExist Sentinel = "player does not exist"
)
