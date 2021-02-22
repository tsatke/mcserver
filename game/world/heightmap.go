package world

type Heightmapper interface {
	Chunk
	MotionBlocking() []int64
}

func ChunkHeightmap(ch Chunk) []int64 {
	if hm, ok := ch.(Heightmapper); ok {
		return hm.MotionBlocking()
	}
	panic("implement me")
}
