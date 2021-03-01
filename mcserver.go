package mcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"

	"github.com/tsatke/mcserver/config"
	"github.com/tsatke/mcserver/game"
	"github.com/tsatke/mcserver/game/chat"
	"github.com/tsatke/mcserver/game/world"
	"github.com/tsatke/mcserver/network"
	"github.com/tsatke/mcserver/network/packet"
)

const (
	// ServerVersion is the vanilla server version that is mostly implemented by this server.
	ServerVersion = "1.16.5"
	// ProtocolVersion is the minecraft protocol version that this server implements.
	ProtocolVersion = 754
)

// MCServer is a minecraft server. It holds things like a logger, a net.Listener, a game.Game and a config.Config,
// and coordinates all the components.
type MCServer struct {
	log      zerolog.Logger
	addr     string
	listener net.Listener

	game   *game.Game
	config config.Config
}

// New creates a new MCServer with the given config. This server will not use a logger. Use WithLogger if you
// want the server to generate log output.
func New(config config.Config, opts ...Option) (*MCServer, error) {
	srv := &MCServer{
		log:    zerolog.Nop(),
		addr:   config.ServerAddr(),
		config: config,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.listener == nil {
		lis, err := net.Listen("tcp", srv.addr)
		if err != nil {
			return nil, fmt.Errorf("listen: %w", err)
		}
		srv.listener = lis
	}

	return srv, nil
}

// Start will prepare the game. The given context will be respected. This method will go into an infinite
// loop, accepting incoming connections. This method terminates only if an error occurs or if the context
// was cancelled.
func (s *MCServer) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		// wait for context cancellation and close listener
		_ = s.listener.Close()
	}()

	s.log.Info().
		Msg("preparing game")
	if err := s.prepareGame(ctx); err != nil {
		return fmt.Errorf("prepare game: %w", err)
	}

	s.log.Info().
		Str("addr", s.addr).
		Msg("waiting for incoming connection")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// check if context was cancelled
			select {
			case <-ctx.Done():
				// if the context was cancelled, ignore the error
				s.log.Debug().
					Msg("stopped waiting for incoming connections")
				// listener is closed in separate goroutine
				return nil
			default:
				// otherwise, return the error, interrupting the loop
				return fmt.Errorf("accept: %w", err)
			}
		}
		s.log.Info().
			IPAddr("from", conn.RemoteAddr().(*net.TCPAddr).IP).
			Msg("incoming connection")
		go s.handleRequest(conn)
	}
}

func (s *MCServer) prepareGame(ctx context.Context) error {
	start := time.Now()
	w, err := world.LoadVanilla(afero.NewBasePathFs(afero.NewOsFs(), s.config.GameWorld()))
	if err != nil {
		return fmt.Errorf("load world: %w", err)
	}
	s.log.Info().
		Stringer("took", time.Since(start)).
		Msg("loaded world")

	g, err := game.New(
		w,
		game.WithLogger(s.log.With().
			Str("component", "game").
			Logger()),
	)
	if err != nil {
		return fmt.Errorf("create game: %w", err)
	}

	s.game = g
	go s.game.Start(ctx)
	s.log.Debug().
		Msg("wait for game to be ready")
	select {
	case <-ctx.Done():
		s.log.Info().
			Msg("aborted start")
	case <-s.game.Ready():
		s.log.Info().
			Msg("game ready")
	}
	return nil
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
