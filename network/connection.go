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

// Phase returns the packet.Phase that this connection is currently in.
func (c Conn) Phase() packet.Phase {
	return c.phase
}

// IP returns the IP address of the remote end of this connection.
func (c Conn) IP() net.IP {
	return c.underlying.RemoteAddr().(*net.TCPAddr).IP
}

// TransitionTo sets the phase of this connection to the given phase.
// This method will panic if the transition is invalid.
func (c *Conn) TransitionTo(phase packet.Phase) {
	if !transitionValid(c.phase, phase) {
		panic(fmt.Sprintf("cannot transition from %s to %s", c.phase, phase))
	}
	c.log.Debug().
		Stringer("to", phase).
		Msg("transition connection")
	c.phase = phase
}

func transitionValid(from, to packet.Phase) bool {
	switch from {
	case packet.PhaseHandshaking:
		return to == packet.PhaseLogin || to == packet.PhaseStatus
	case packet.PhaseLogin:
		return to == packet.PhasePlay
	}
	return false
}

// ReadPacket will attempt to read and decode a serverbound packet from this connection.
// This method will block until either the connection is closed, in which case this method
// will return an error, or until a serverbound packet is received. If the packet is malformed
// or contains invalid values, then a respective error that is describing the issue is returned.
// A packet is considered as containing invalid values if the packet definition implements the
// packet.Validator interface and the validation fails.
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

// WritePacket writes the given clientbound packet to the client, returning any error if any
// returns.
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
