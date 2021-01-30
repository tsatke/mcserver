package network

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog"

	"github.com/tsatke/mcserver/network/packet"
)

type Conn struct {
	log        zerolog.Logger
	underlying net.Conn
	state      packet.State
	closed     bool
}

func NewConn(log zerolog.Logger, underlying net.Conn) *Conn {
	return &Conn{
		log:        log,
		underlying: underlying,
		state:      packet.StateHandshaking,
	}
}

func (c Conn) State() packet.State {
	return c.state
}

func (c Conn) IP() net.IP {
	return c.underlying.RemoteAddr().(*net.TCPAddr).IP
}

func (c *Conn) SetState(state packet.State) {
	c.state = state
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	if c.closed {
		return nil, ErrClosed
	}

	packet, err := packet.Decode(c.underlying, c.state)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = c.Close()
			return nil, ErrClosed
		}
		return nil, err
	}
	c.log.Trace().
		Str("name", packet.Name()).
		Msg("received packet")
	return packet, nil
}

func (c *Conn) WritePacket(p packet.Packet) error {
	if c.closed {
		return ErrClosed
	}

	clientbound, ok := p.(packet.Clientbound)
	if !ok {
		return fmt.Errorf("packet %s does not implement Clientbound", clientbound.Name())
	}

	if err := packet.Encode(clientbound, c.underlying); err != nil {
		return fmt.Errorf("encode into: %w", err)
	}
	c.log.Trace().
		Str("name", p.Name()).
		Msg("sent packet")
	return nil
}

// Close closes this connection. This method is idempotent.
func (c *Conn) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true
	return c.underlying.Close()
}
