package packet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

var (
	packetInterfaceType    = reflect.TypeOf((*Packet)(nil)).Elem()
	serverboundPacketTypes = make(map[Phase]map[ID]reflect.Type)
	clientboundPacketTypes = make(map[Phase]map[ID]reflect.Type)
)

// RegisterPacket will associate the given state with the given type.
// The given type must implement Packet and either Serverbound
// or Clientbound. If a packet implements both (which wouldn't make much
// sense, but it's possible), Serverbound takes precedence.
//
// Registering a packet with a state will allow Decode to automatically
// create and decode an incoming packet depending on a Phase and the
// decoded packet ID.
func RegisterPacket(state Phase, typ reflect.Type) {
	if !typ.Implements(packetInterfaceType) {
		panic(fmt.Sprintf("%s does not implement the Packet interface", typ.Name()))
	}
	created := reflect.New(typ).Interface()
	id := created.(Packet).ID()
	var target map[Phase]map[ID]reflect.Type
	if _, ok := created.(Serverbound); ok {
		// initialize state map if necessary
		if serverboundPacketTypes[state] == nil {
			serverboundPacketTypes[state] = make(map[ID]reflect.Type)
		}
		target = serverboundPacketTypes
	} else if _, ok := created.(Clientbound); ok {
		// initialize state map if necessary
		if clientboundPacketTypes[state] == nil {
			clientboundPacketTypes[state] = make(map[ID]reflect.Type)
		}
		target = clientboundPacketTypes
	} else {
		panic(fmt.Sprintf("%s is a packet, but does neither implement Serverbound nor Clientbound", typ.Name()))
	}

	if target[state][id] != nil {
		panic(fmt.Sprintf("already registered packet %T in state %s", created, state))
	}
	target[state][id] = typ
}

// Packet is a packet that can either be sent to a client or the server.
type Packet interface {
	ID() ID
	Name() string
}

// Serverbound describes a packet that was sent from the client to the server.
// Every Serverbound packet has a DecodeFrom method that fills its values with
// data from a reader. The reader will always be uncompressed. The reader will
// not contain the packet length and/or ID. The reader will return EOF if the
// packet reaches an end. That means, that the reader can not over-read.
type Serverbound interface {
	Packet
	// DecodeFrom will read fields according to this packet from the given reader.
	// The reader must only contain the packet fields, not the length or ID.
	// It's probably best to pass in an io.LimitedReader, however, the implementation
	// must not rely on that.
	DecodeFrom(io.Reader) error
}

// Clientbound describes a packet that can be sent from the server to a client.
// Each Clientbound packet has an EncodeInto method that will encode only
// the packet fields onto the passed in writer.
type Clientbound interface {
	Packet
	// EncodeInto will only write the packet fields onto the given writer, NOT the length and/or ID.
	EncodeInto(io.Writer) error
}

// Encode will write the given packet onto the given writer. As opposed to the
// Clientbound.EncodeInto method, this will also write the packet length and ID.
func Encode(pkg Clientbound, w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	var buf bytes.Buffer
	enc := Encoder{&buf}
	enc.WriteVarInt("packet ID", int(pkg.ID()))
	panicIffErr("packet", pkg.EncodeInto(&buf))
	Encoder{w}.WriteVarInt("packet length", buf.Len())
	if _, err := buf.WriteTo(w); err != nil {
		return fmt.Errorf("write to: %w", err)
	}
	return
}

// Decode decodes a Serverbound packet from the given reader, depending on the
// given state.
func Decode(rd io.Reader, state Phase) (p Serverbound, err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	packetLen := dec.ReadVarInt("packet length")
	packetID := ID(dec.ReadVarInt("packet ID"))
	payloadLength := packetLen - 1 /* packetID.Len(), which seems to always be 1 */
	payloadReader := io.LimitReader(rd, int64(payloadLength))
	packetType := serverboundPacketTypes[state][packetID]
	if packetType == nil {
		_, _ = io.CopyN(ioutil.Discard, rd, int64(payloadLength)) // discard remaining bytes of the packet
		return nil, fmt.Errorf("unknown ID %s in Phase %s, discarding", packetID, state)
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
