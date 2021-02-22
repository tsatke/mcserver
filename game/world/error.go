package world

type sentinel string

func (s sentinel) Error() string { return string(s) }

const (
	ErrChunkNotExist sentinel = "chunk does not exist"
	// ErrChunkNotGenerated indicates that the requested chunk has not been generated yet.
	ErrChunkNotGenerated sentinel = "chunk not generated"
)
