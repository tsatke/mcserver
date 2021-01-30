package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundServerDifficulty{}))
}

type Difficulty uint8

const (
	DifficultyPeaceful Difficulty = iota
	DifficultyEasy
	DifficultyNormal
	DifficultyHard
)

type ClientboundServerDifficulty struct {
	Difficulty       Difficulty
	DifficultyLocked bool
}

func (ClientboundServerDifficulty) ID() ID       { return IDClientboundServerDifficulty }
func (ClientboundServerDifficulty) Name() string { return "Server Difficulty" }

func (c ClientboundServerDifficulty) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeUbyte("difficulty", uint8(c.Difficulty))
	enc.writeBoolean("difficulty locked", c.DifficultyLocked)

	return
}
