package network

import (
	"io"
	"net"
	"testing"
	"testing/iotest"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
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
		err := source.WritePacket(packet.ClientboundLoginSuccess{
			UUID:     uuid,
			Username: "aUsername",
		})
		if err != nil {
			panic(err)
		}
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
func (suite *ConnectionSuite) TestReadPacketMalformedWrite() {
	net1, net2 := net.Pipe()
	sink := NewConn(zerolog.Nop(), net1)
	go func() {
		wr := iotest.TruncateWriter(net2, 13)
		enc := packet.Encoder{wr}
		enc.WriteVarInt("packet length", 19) // remember to change this when you change the values below
		enc.WriteVarInt("packet id", int(packet.IDServerboundHandshake))
		enc.WriteVarInt("protocol version", 111)
		enc.WriteString("server address", "my.server.com") // the writer aborts somewhere in the middle of the string content
		enc.WriteUshort("server port", 65432)
		enc.WriteVarInt("next state", int(packet.NextStateStatus))

		/*
			We close the connection here. In prod, the read would be stuck until either
			the rest of the packet is written or the connection is closed. We close the
			connection and expect the correct error.
		*/
		_ = net2.Close()
	}()
	p, err := sink.ReadPacket()
	suite.NotEqual(io.EOF, err) // make sure that the error is not simply EOF, since we stop writing in the middle of a string
	// the error message may change, however, it must not be simply io.EOF
	suite.EqualError(err, "decode Handshake: server address: prefix indicated length of 13, but only got 9 bytes")
	suite.Nil(p)
}

func (suite *ConnectionSuite) TestReadPacketClosedConnection() {
	net1, net2 := net.Pipe()
	sink := NewConn(zerolog.Nop(), net1)
	go func() {
		// instead of writing a packet, we're going to close the pipe
		_ = net2.Close()
	}()
	p, err := sink.ReadPacket()
	suite.ErrorIs(err, io.EOF)
	suite.Nil(p)
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

func Test_transitionValid(t *testing.T) {
	type args struct {
		from packet.Phase
		to   packet.Phase
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// from invalid
		{
			"invalid invalid",
			args{packet.PhaseInvalid, packet.PhaseInvalid},
			false,
		},
		{
			"invalid handshaking",
			args{packet.PhaseInvalid, packet.PhaseHandshaking},
			false,
		},
		{
			"invalid status",
			args{packet.PhaseInvalid, packet.PhaseStatus},
			false,
		},
		{
			"invalid login",
			args{packet.PhaseInvalid, packet.PhaseLogin},
			false,
		},
		{
			"invalid play",
			args{packet.PhaseInvalid, packet.PhasePlay},
			false,
		},
		// from handshaking
		{
			"handshaking handshaking",
			args{packet.PhaseHandshaking, packet.PhaseHandshaking},
			false,
		},
		{
			"handshaking status",
			args{packet.PhaseHandshaking, packet.PhaseStatus},
			true,
		},
		{
			"handshaking login",
			args{packet.PhaseHandshaking, packet.PhaseLogin},
			true,
		},
		{
			"handshaking play",
			args{packet.PhaseHandshaking, packet.PhasePlay},
			false,
		},
		// from status
		{
			"status handshaking",
			args{packet.PhaseStatus, packet.PhaseHandshaking},
			false,
		},
		{
			"status status",
			args{packet.PhaseStatus, packet.PhaseStatus},
			false,
		},
		{
			"status login",
			args{packet.PhaseStatus, packet.PhaseLogin},
			false,
		},
		{
			"status play",
			args{packet.PhaseStatus, packet.PhasePlay},
			false,
		},
		// from login
		{
			"login handshaking",
			args{packet.PhaseLogin, packet.PhaseHandshaking},
			false,
		},
		{
			"login status",
			args{packet.PhaseLogin, packet.PhaseStatus},
			false,
		},
		{
			"login login",
			args{packet.PhaseLogin, packet.PhaseLogin},
			false,
		},
		{
			"login play",
			args{packet.PhaseLogin, packet.PhasePlay},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := transitionValid(tt.args.from, tt.args.to)
			assert.Equal(tt.want, got)
		})
	}
}
