package mcserver

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"golang.org/x/net/nettest"

	"github.com/tsatke/mcserver/config"
	"github.com/tsatke/mcserver/network/packet"
)

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

type ServerSuite struct {
	suite.Suite

	listener  net.Listener
	openConns []net.Conn
	server    *MCServer
	cancelFn  func()
}

func testConfig() config.Config {
	vp := viper.New()
	vp.Set(config.KeyGameWorld, "game/testdata/maps/world01")

	cfg := config.New(vp)
	return cfg
}

func (suite *ServerSuite) SetupTest() {
	lis, err := nettest.NewLocalListener("tcp")
	suite.Require().NoError(err)
	suite.listener = lis

	srv, err := New(
		zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger(),
		testConfig(),
		WithListener(lis),
	)
	suite.Require().NoError(err)
	suite.server = srv

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		suite.NoError(srv.Start(ctx))
	}()
	suite.cancelFn = cancel
}

func (suite *ServerSuite) TearDownTest() {
	for _, conn := range suite.openConns {
		_ = conn.Close()
	}
	suite.openConns = nil

	if suite.cancelFn != nil {
		suite.cancelFn()
		suite.cancelFn = nil
	}
	if suite.listener != nil {
		_ = suite.listener.Close()
	}
	if suite.server != nil {
		suite.server.Stop()
	}
}

// DialServer will return a connection to the test server, which will automatically be
// closed after the test finished. This means, that there's no need for the tester
// to explicitly close the connection. It is allowed however.
func (suite *ServerSuite) DialServer() net.Conn {
	conn, err := net.Dial("tcp", suite.listener.Addr().String())
	suite.Require().NoError(err)
	suite.openConns = append(suite.openConns, conn)
	return conn
}

// DoSend allows the tester to provide a encoding callback to encode packet values.
// The encoded values will be sent as a valid packet (prefixed with length and the
// provided ID) to the given writer.
func (suite *ServerSuite) DoSend(to io.Writer, id packet.ID, fn func(packet.Encoder)) {
	var buf bytes.Buffer
	bufferEnc := packet.Encoder{&buf}
	bufferEnc.WriteVarInt("packet id", int(id))
	suite.NotPanics(func() {
		fn(bufferEnc)
	})

	packetEnc := packet.Encoder{to}
	packetEnc.WriteVarInt("packet length", buf.Len())
	_, err := buf.WriteTo(to)
	suite.NoError(err)
}

func (suite *ServerSuite) DoReceive(from io.Reader, fn func(packet.ID, packet.Decoder)) {
	dec := packet.Decoder{from}
	packetLength := dec.ReadVarInt("packet length")
	packetID := packet.ID(dec.ReadVarInt("packet id"))
	rd := io.LimitReader(from, int64(packetLength-1))
	suite.NotPanics(func() {
		fn(packetID, packet.Decoder{rd})
	})
}

// EOF assumes that the given reader is closed, and will fail the test if it is not.
// A reader is considered closed if a read returns an io.EOF. Please note that this
// method reads one byte from the given reader.
func (suite *ServerSuite) EOF(rd io.Reader) {
	_, err := rd.Read([]byte{0})
	suite.ErrorIs(err, io.EOF)
}
