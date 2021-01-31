package mcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"

	"github.com/tsatke/mcserver/game"
	"github.com/tsatke/mcserver/game/chat"
	"github.com/tsatke/mcserver/network"
	"github.com/tsatke/mcserver/network/packet"
)

const (
	ServerVersion   = "1.16.5"
	ProtocolVersion = 754
)

type MCServer struct {
	log      zerolog.Logger
	addr     string
	listener net.Listener

	game *game.Game
}

func New(log zerolog.Logger, addr string) *MCServer {
	return &MCServer{
		log:  log,
		addr: addr,
	}
}

func (s *MCServer) Start(ctx context.Context) error {
	s.log.Info().
		Msg("preparing game")
	s.prepareGame(ctx)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	s.listener = lis
	s.log.Info().
		Str("addr", s.addr).
		Msg("waiting for incoming connection")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				s.log.Debug().
					Msg("stopped waiting for incoming connections")
				return nil
			default:
				return fmt.Errorf("accept: %w", err)
			}
		}
		s.log.Info().
			IPAddr("from", conn.RemoteAddr().(*net.TCPAddr).IP).
			Msg("incoming connection")
		go s.handleRequest(conn)
	}
}

func (s *MCServer) Stop() {
	_ = s.listener.Close()
}

func (s *MCServer) prepareGame(ctx context.Context) {
	s.game = game.New(
		s.log.With().
			Str("component", "game").
			Logger(),
		afero.NewBasePathFs(afero.NewOsFs(), "world"),
	)
	s.game.Start(ctx)
	select {
	case <-ctx.Done():
		s.log.Info().
			Msg("aborted start")
	case <-s.game.Ready():
		s.log.Info().
			Msg("game ready")
	}
}

func (s *MCServer) handleRequest(c net.Conn) {
	conn := network.NewConn(
		s.log.With().
			IPAddr("remote", c.RemoteAddr().(*net.TCPAddr).IP).
			Logger(),
		c,
	)
	p, err := conn.ReadPacket()
	if err != nil {
		s.log.Debug().
			Err(err).
			Msg("read packet failed, closing connection")
		_ = c.Close()
		return
	}

	handshake, ok := p.(*packet.ServerboundHandshake)
	if !ok {
		s.log.Debug().
			Int("gotID", int(p.ID())).
			Int("wantID", int(packet.IDServerboundHandshake)).
			Msg("require handshake packet")
		_ = c.Close()
		return
	}

	switch handshake.NextState {
	case packet.NextStateStatus:
		conn.TransitionTo(packet.PhaseStatus)
		s.handleStatusRequest(conn)
	case packet.NextStateLogin:
		conn.TransitionTo(packet.PhaseLogin)
		s.handleLoginRequest(conn)
	default:
		s.log.Error().
			Stringer("nextstate", handshake.NextState).
			Msg("invalid next state")
		_ = c.Close()
	}
}

func (s *MCServer) handleStatusRequest(conn *network.Conn) {
	_, err := conn.ReadPacket()
	if err != nil {
		s.log.Debug().
			Err(err).
			Msg("read request failed, closing connection")
		_ = conn.Close()
		return
	}

	if err := conn.WritePacket(packet.ClientboundResponse{
		JSONResponse: packet.Response{
			Version: packet.ResponseVersion{
				Name:     ServerVersion,
				Protocol: ProtocolVersion,
			},
			Players: packet.ResponsePlayers{
				Max:    100,
				Online: s.game.AmountOfConnectedPlayers(),
				Sample: []packet.ResponsePlayersSample{},
			},
			Description: chat.Chat{
				ChatFragment: chat.ChatFragment{
					Text: "Timi loves Tanni ",
				},
				Extra: []chat.ChatFragment{
					{
						Text:  "❤",
						Color: "red",
					},
				},
			},
		},
	}); err != nil {
		s.log.Debug().
			Err(err).
			Msg("write response failed, closing connection")
		_ = conn.Close()
		return
	}

	p, err := conn.ReadPacket()
	if err != nil {
		s.log.Debug().
			Err(err).
			Msg("receive ping failed, closing connection")
		_ = conn.Close()
		return
	}

	ping, ok := p.(*packet.ServerboundPing)
	if !ok {
		s.log.Debug().
			Int("gotID", int(p.ID())).
			Int("wantID", int(packet.IDServerboundPing)).
			Msg("require ping packet")
		_ = conn.Close()
		return
	}

	timestamp := ping.Payload
	if err := conn.WritePacket(packet.ClientboundPong{
		Payload: timestamp,
	}); err != nil {
		s.log.Debug().
			Err(err).
			Msg("write pong failed, closing connection")
		_ = conn.Close()
		return
	}

	_ = conn.Close()
}

func (s *MCServer) handleLoginRequest(conn *network.Conn) {
	p, err := conn.ReadPacket()
	if err != nil {
		s.log.Debug().
			Err(err).
			Msg("read packet failed, closing connection")
		_ = conn.Close()
		return
	}

	loginStart, ok := p.(*packet.ServerboundLoginStart)
	if !ok {
		s.log.Debug().
			Int("gotID", int(p.ID())).
			Int("wantID", int(packet.IDServerboundLoginStart)).
			Msg("require login start packet")
		return
	}

	username := loginStart.Username
	s.log.Info().
		Str("username", username).
		Msg("player trying to connect")

	uid := uuid.New()
	if err := conn.WritePacket(packet.ClientboundLoginSuccess{
		UUID:     uid, // TODO: get uuid from game data
		Username: username,
	}); err != nil {
		s.log.Debug().
			Err(err).
			Msg("write packet failed, closing connection")
		_ = conn.Close()
		return
	}

	conn.TransitionTo(packet.PhasePlay)
	s.game.AddPlayer(game.NewPlayer(uid, username, conn))
	// game handles the connection as of here, nothing more to do
}
