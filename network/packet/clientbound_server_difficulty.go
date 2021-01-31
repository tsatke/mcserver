package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundServerDifficulty{}))
}

type ClientboundServerDifficulty struct {
	Difficulty       byte
	DifficultyLocked bool
}

func (ClientboundServerDifficulty) ID() ID       { return IDClientboundServerDifficulty }
func (ClientboundServerDifficulty) Name() string { return "Server Difficulty" }

func (c ClientboundServerDifficulty) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteUbyte("difficulty", c.Difficulty)
	enc.WriteBoolean("difficulty locked", c.DifficultyLocked)

	return
}
