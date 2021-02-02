package network

import (
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/tsatke/mcserver/network/packet"
)

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionSuite))
}

type ConnectionSuite struct {
	suite.Suite
}

func (suite *ConnectionSuite) TestWritePacket() {
	net1, net2 := net.Pipe()
	source := NewConn(zerolog.Nop(), net1)
	uuid := uuid.New()
	go func() {
		suite.NoError(source.WritePacket(packet.ClientboundLoginSuccess{
			UUID:     uuid,
			Username: "aUsername",
		}))
	}()
	defer func() {
		rec := recover()
		suite.Nilf(rec, "%v", rec)
	}()
	dec := packet.Decoder{net2}
	suite.EqualValues(27, dec.ReadVarInt("packet length"))
	suite.EqualValues(packet.IDClientboundLoginSuccess, dec.ReadVarInt("packet id"))
	suite.EqualValues(uuid, dec.ReadUUID("uuid"))
	suite.EqualValues("aUsername", dec.ReadString("username"))
}

func (suite *ConnectionSuite) TestReadPacket() {
	net1, net2 := net.Pipe()
	sink := NewConn(zerolog.Nop(), net1)
	go func() {
		defer func() {
			rec := recover()
			suite.Nilf(rec, "%v", rec)
		}()
		enc := packet.Encoder{net2}
		enc.WriteVarInt("packet length", 19) // remember to change this when you change the values below
		enc.WriteVarInt("packet id", int(packet.IDServerboundHandshake))
		enc.WriteVarInt("protocol version", 111)
		enc.WriteString("server address", "my.server.com")
		enc.WriteUshort("server port", 65432)
		enc.WriteVarInt("next state", int(packet.NextStateStatus))
	}()
	p, err := sink.ReadPacket()
	suite.NoError(err)
	suite.Equal(&packet.ServerboundHandshake{
		ProtocolVersion: 111,
		ServerAddress:   "my.server.com",
		ServerPort:      65432,
		NextState:       1,
	}, p)
}

func (suite *ConnectionSuite) TestClose() {
	net, _ := net.Pipe()
	conn := NewConn(zerolog.Nop(), net)

	suite.False(conn.closed)
	suite.NoError(conn.Close())

	suite.True(conn.closed)
	suite.NoError(conn.Close())
	suite.True(conn.closed)
}
