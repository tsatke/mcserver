package packet

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

var (
	packetInterfaceType    = reflect.TypeOf((*Packet)(nil)).Elem()
	serverboundPacketTypes = make(map[State]map[ID]reflect.Type)
	clientboundPacketTypes = make(map[State]map[ID]reflect.Type)
)

func RegisterPacket(state State, typ reflect.Type) {
	if !typ.Implements(packetInterfaceType) {
		panic(fmt.Sprintf("%s does not implement the Packet interface", typ.Name()))
	}
	created := reflect.New(typ).Interface()
	id := created.(Packet).ID()
	if _, ok := created.(Serverbound); ok {
		if serverboundPacketTypes[state] == nil {
			serverboundPacketTypes[state] = make(map[ID]reflect.Type)
		}
		if serverboundPacketTypes[state][id] != nil {
			panic(fmt.Sprintf("already registered serverbound packet %T in state %s", created, state))
		}
		serverboundPacketTypes[state][id] = typ
	} else if _, ok := created.(Clientbound); ok {
		if clientboundPacketTypes[state] == nil {
			clientboundPacketTypes[state] = make(map[ID]reflect.Type)
		}
		if clientboundPacketTypes[state][id] != nil {
			panic(fmt.Sprintf("already registered clientbound packet %T in state %s", created, state))
		}
		clientboundPacketTypes[state][id] = typ
	} else {
		panic(fmt.Sprintf("%s is a packet, but does neither implement Serverbound nor Clientbound", typ.Name()))
	}
}

// Packet is a packet that can either be sent to a client or the server.
type Packet interface {
	ID() ID
	Name() string
}

// Serverbound describes a packet that was sent from the client to the server.
// Every Serverbound packet has a DecodeFrom method that fills its values with
// data from a reader. The reader will always be uncompressed. Other than that,
// there are no guarantees about the reader. The packet implementation must make
// sure that the stream is read from correctly. The reader might be the network
// connection itself, so make sure to check how many bytes are actually read with
// each read.
type Serverbound interface {
	Packet
	// DecodeFrom will read fields according to this packet from the given reader.
	// The reader must only contain the packet fields, not the length or ID.
	// It's probably best to pass in an io.LimitedReader, however, the implementation
	// must not rely on that.
	DecodeFrom(io.Reader) error
}

type Clientbound interface {
	Packet
	// EncodeInto will only write the packet fields onto the given writer, NOT the length and/or ID.
	EncodeInto(io.Writer) error
}

func Encode(pkg Clientbound, w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	var buf bytes.Buffer
	enc := Encoder{&buf}
	enc.WriteVarInt("packet ID", int(pkg.ID()))
	panicIffErr("packet", pkg.EncodeInto(&buf))
	enc.WriteVarInt("packet length", buf.Len())
	if _, err := buf.WriteTo(w); err != nil {
		return fmt.Errorf("write to: %w", err)
	}
	return
}

func Decode(rd io.Reader, state State) (p Serverbound, err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	packetLen := dec.ReadVarInt("packet length")
	packetID := ID(dec.ReadVarInt("packet ID"))
	payloadLength := packetLen - 1 /* packetID.Len(), which seems to always be 1 */
	payloadReader := io.LimitReader(rd, int64(payloadLength))
	packetType := serverboundPacketTypes[state][packetID]
	if packetType == nil {
		return nil, fmt.Errorf("unknown ID %s in State %s", packetID, state)
	}
	packetInterface := reflect.New(packetType).Interface()
	packet := packetInterface.(Serverbound)
	if err := packet.DecodeFrom(payloadReader); err != nil {
		return nil, fmt.Errorf("decode %s: %w", packet.Name(), err)
	}
	if validator, ok := packetInterface.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return nil, fmt.Errorf("validate %s: %w", packet.Name(), err)
		}
	}
	return packet, nil
}
