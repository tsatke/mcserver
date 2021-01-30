package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tsatke/mcserver/network/packet/types"
)

func Encode(pkg Clientbound, w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	var buf bytes.Buffer
	panicIffErr("packet ID", types.NewVarInt(int(pkg.ID())).EncodeInto(&buf))
	panicIffErr("packet", pkg.EncodeInto(&buf))
	panicIffErr("packet length", types.NewVarInt(buf.Len()).EncodeInto(w))
	if _, err := buf.WriteTo(w); err != nil {
		return fmt.Errorf("write to: %w", err)
	}
	return
}
