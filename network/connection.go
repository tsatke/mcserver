package network

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"

	"github.com/tsatke/mcserver/network/packet"
)

type Conn struct {
	log        zerolog.Logger
	underlying net.Conn
	phase      packet.Phase
	closed     bool
}

func NewConn(log zerolog.Logger, underlying net.Conn) *Conn {
	return &Conn{
		log:        log,
		underlying: underlying,
		phase:      packet.PhaseHandshaking,
	}
}

func (c Conn) Phase() packet.Phase {
	return c.phase
}

func (c Conn) IP() net.IP {
	return c.underlying.RemoteAddr().(*net.TCPAddr).IP
}

func (c *Conn) TransitionTo(phase packet.Phase) {
	c.log.Debug().
		Stringer("to", phase).
		Msg("transition connection")
	c.phase = phase
}

func (c *Conn) ReadPacket() (packet.Serverbound, error) {
	if c.closed {
		return nil, ErrClosed
	}

	packet, err := packet.Decode(c.underlying, c.phase)
	if err != nil {
		return nil, err
	}
	c.log.Trace().
		Str("name", packet.Name()).
		Msg("received packet")
	return packet, nil
}

func (c *Conn) WritePacket(p packet.Clientbound) error {
	if c.closed {
		return ErrClosed
	}

	if err := packet.Encode(p, c.underlying); err != nil {
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
