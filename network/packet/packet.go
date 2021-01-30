package packet

import (
	"fmt"
	"io"
	"reflect"
)

var (
	packetInterfaceType    = reflect.TypeOf((*Packet)(nil)).Elem()
	serverboundPacketTypes = make(map[State]map[ID]reflect.Type)
	clientboundPacketTypes = make(map[State]map[ID]reflect.Type)
)

func registerPacket(state State, typ reflect.Type) {
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
