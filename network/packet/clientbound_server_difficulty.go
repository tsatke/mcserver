package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundServerDifficulty{}))
}

// ClientboundServerDifficulty is used by the server to tell
// the client the current server difficulty.
type ClientboundServerDifficulty struct {
	// Difficulty is the current server difficulty, where
	// 0=peaceful, 1=easy, 2=normal, 3=hard. Other values
	// are invalid.
	Difficulty byte
	// DifficultyLocked indicates whether the client can
	// change the difficulty value. If this is true,
	// the difficulty button in the client's settings
	// is disabled.
	DifficultyLocked bool
}

// ID returns the constant packet ID.
func (ClientboundServerDifficulty) ID() ID { return IDClientboundServerDifficulty }

// Name returns the constant packet name.
func (ClientboundServerDifficulty) Name() string { return "Server Difficulty" }

// EncodeInto writes this packet into the given writer.
func (c ClientboundServerDifficulty) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteUbyte("difficulty", c.Difficulty)
	enc.WriteBoolean("difficulty locked", c.DifficultyLocked)

	return
}
