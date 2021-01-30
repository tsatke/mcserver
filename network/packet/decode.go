package packet

import (
	"fmt"
	"io"
	"reflect"
)

func Decode(rd io.Reader, state State) (p Packet, err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	packetLen := dec.readVarInt("packet length")
	packetID := ID(dec.readVarInt("packet ID"))
	payloadLength := packetLen - 1 /* packetID.Len(), which seems to always be 1 */
	payloadReader := io.LimitReader(rd, int64(payloadLength))
	packetType := serverboundPacketTypes[state][packetID]
	if packetType == nil {
		return nil, fmt.Errorf("unknown ID %s in State %s", packetID, state)
	}
	packet := reflect.New(packetType).Interface().(Serverbound)
	if err := packet.DecodeFrom(payloadReader); err != nil {
		return nil, fmt.Errorf("decode %s: %w", packet.Name(), err)
	}
	return packet, nil
}
