package mcserver

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

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

	// listener is the listener that the server of this suite uses.
	listener net.Listener

	// openConnsLock protects openConns.
	openConnsLock sync.Mutex
	// openConns holds all connections that were opened using the suite
	// to connect to the test server.
	openConns []net.Conn

	// ctx is the context with which the test server was started.
	ctx context.Context
	// server is the test server.
	server *MCServer
	// cancelFn is the cancel function of the test context.
	cancelFn func()
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

	srv, err := New(testConfig(),
		WithListener(lis),
		WithLogger(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.TraceLevel).With().Timestamp().Logger()),
	)
	suite.Require().NoError(err)
	suite.server = srv

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	suite.ctx = ctx
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
}

// DialServer will return a connection to the test server, which will automatically be
// closed after the test finished. This means, that there's no need for the tester
// to explicitly close the connection. It is allowed however.
func (suite *ServerSuite) DialServer() net.Conn {
	conn, err := net.Dial("tcp", suite.listener.Addr().String())
	suite.Require().NoError(err)

	suite.openConnsLock.Lock()
	suite.openConns = append(suite.openConns, conn)
	suite.openConnsLock.Unlock()

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

// DoReceive attempts to receive a message from the server. It attempts to read and
// call the given function with a timeout of 5 seconds.
func (suite *ServerSuite) DoReceive(from io.Reader, fn func(packet.ID, packet.Decoder)) {
	ch := make(chan struct{})
	var err error
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				if recErr, ok := rec.(error); ok {
					err = recErr
				} else {
					panic(rec)
				}
			}
			close(ch)
		}()
		dec := packet.Decoder{from}
		packetLength := dec.ReadVarInt("packet length")
		packetID := packet.ID(dec.ReadVarInt("packet id"))
		rd := io.LimitReader(from, int64(packetLength-1))

		fn(packetID, packet.Decoder{rd})
	}()

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		/*
			Might happen because the connection is still open but the server didn't
			send a message. This can not happen when the connection is closed; then
			io.EOF will be returned.
		*/
		suite.FailNow("timeout while receiving")
	}
	suite.NoError(err)
}

// ClosedOrEOF assumes that the given reader is closed, and will fail the test if it is not.
// A reader is considered closed if a read returns an io.EOF. Please note that this
// method reads one byte from the given reader.
func (suite *ServerSuite) ClosedOrEOF(rd io.Reader) {
	ch := make(chan struct{})
	go func() {
		_, err := rd.Read([]byte{0})
		if !errors.Is(err, io.EOF) {
			if netErr, ok := err.(net.Error); !ok && netErr.Temporary() {
				suite.Fail("error is not not EOF or closed")
			}
		}
		close(ch)
	}()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		suite.FailNow("connection not closed")
	}
}
