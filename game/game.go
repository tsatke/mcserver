package game

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/chunk"
	"github.com/tsatke/mcserver/game/entity"
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
	"github.com/tsatke/mcserver/network/packet"
	"github.com/tsatke/mcserver/network/packet/types"
)

const (
	TickDuration = 50 * time.Millisecond

	defaultQueueBufferSize = 100
)

type incomingMessage struct {
	source *Player
	pkg    packet.Serverbound
}

type Game struct {
	log   zerolog.Logger
	ready chan struct{}

	fs    afero.Fs
	world *World

	currentTick int64

	loadedChunks map[voxel.V2]*chunk.Chunk

	connectedPlayers     map[uuid.UUID]*Player
	incomingMessageQueue chan incomingMessage
}

func New(log zerolog.Logger, world afero.Fs) *Game {
	return &Game{
		log:                  log,
		ready:                make(chan struct{}),
		fs:                   world,
		loadedChunks:         make(map[voxel.V2]*chunk.Chunk),
		connectedPlayers:     make(map[uuid.UUID]*Player),
		incomingMessageQueue: make(chan incomingMessage, defaultQueueBufferSize), // TODO: check if 100 is too large, too little or whatever
	}
}

// Start starts the main game loop. To stop, you have to cancel the context passed into this method.
func (g *Game) Start(ctx context.Context) {
	go g.start(ctx)
}

func (g *Game) Ready() <-chan struct{} {
	return g.ready
}

func (g *Game) AmountOfConnectedPlayers() int {
	return len(g.connectedPlayers)
}

func (g *Game) WritePacket(p *Player, pkg packet.Packet) {
	if err := p.conn.WritePacket(pkg); err != nil {
		g.log.Debug().
			Err(err).
			IPAddr("to", p.conn.IP()).
			Str("player", p.name).
			Stringer("uuid", p.UUID).
			Msg("write packet failed, disconnecting player")
		g.Disconnect(p)
	}
}

func (g *Game) DisconnectWithReason(p *Player, reason types.Chat) {
	// we absolutely don't care what happens on the connection anymore, so
	// if the write fails - ok, if it doesn't - ok.
	_ = p.conn.WritePacket(packet.ClientboundDisconnectPlay{
		Reason: reason,
	})
	g.Disconnect(p)
}

func (g *Game) Disconnect(p *Player) {
	p.Disconnect()
	delete(g.connectedPlayers, p.UUID)
}

func (g *Game) AddPlayer(p *Player) {
	p.Player = &entity.Player{
		Mob: entity.Mob{
			Data: entity.Data{
				UUID: uuid.UUID(p.tempUUID),
			},
		},
	}

	// if err := g.loadPlayerEntity(p); err != nil {
	// 	if errors.Is(err, ErrPlayerNotExist) {
	// 		// TODO: create player
	// 	} else {
	// 		g.log.Error().
	// 			Err(err).
	// 			Stringer("uuid", p.UUID).
	// 			Msg("loading player data failed, disconnecting")
	// 		g.DisconnectWithReason(p, types.Chat{
	// 			ChatFragment: types.ChatFragment{
	// 				Text: "Chances are the ser",
	// 			},
	// 			Extra: []types.ChatFragment{
	// 				{
	// 					Text:       "v",
	// 					Obfuscated: true,
	// 				},
	// 				{
	// 					Text: "er is broken.\n",
	// 				},
	// 				{
	// 					Text: "We weren't able to load your player profile.\n",
	// 					Bold: true,
	// 				},
	// 				{
	// 					Text: "Sorry!",
	// 				},
	// 			},
	// 		})
	// 		return
	// 	}
	// }

	g.connectedPlayers[p.UUID] = p

	g.log.Info().
		Stringer("uuid", p.UUID).
		Str("username", string(p.name)).
		Msg("player connected")

	dimensionCodec := codec

	g.sendJoinGameMessage(p, dimensionCodec)
	g.sendServerDifficulty(p)

	g.WritePacket(p, packet.ClientboundHeldItemChange{
		Slot: 0,
	})
	g.WritePacket(p, packet.ClientboundDeclareRecipes{
		Recipes: []packet.Recipe{}, // TODO: declare recipes
	})
	g.sendTags(p)
	g.WritePacket(p, packet.ClientboundEntityStatus{
		EntityID: 1,  // same EID as when joining
		Status:   23, // disable reduced debug screen info
	})
	g.WritePacket(p, packet.ClientboundPlayerPositionAndLook{
		X:          0,
		Y:          69,
		Z:          0,
		Yaw:        0,
		Pitch:      0,
		Flags:      0,
		TeleportID: 12,
	})
	g.WritePacket(p, packet.ClientboundPlayerInfo{
		Action: packet.PlayerInfoActionAddPlayer,
		Players: []packet.PlayerInfoPlayer{
			{
				UUID:           p.UUID,
				Name:           p.name,
				Gamemode:       packet.GamemodeSurvival,
				Ping:           100,
				HasDisplayName: false,
			},
		},
	})
	g.WritePacket(p, packet.ClientboundPlayerInfo{
		Action: packet.PlayerInfoUpdateLatency,
		Players: []packet.PlayerInfoPlayer{
			{
				Ping: 32,
			},
		},
	})

	go g.handleIncomingPlayerMessages(p)
}

func (g *Game) handleIncomingPlayerMessages(p *Player) {
	for {
		pkg, err := p.conn.ReadPacket()
		if err != nil {
			g.log.Error().
				Err(err).
				Stringer("player", p.UUID).
				Msg("read packet failed, disconnect")
			g.Disconnect(p)
			return
		}

	retryLoop:
		for {
			select {
			case g.incomingMessageQueue <- incomingMessage{
				source: p,
				pkg:    pkg.(packet.Serverbound),
			}:
				break retryLoop
			default:
				waitTime := 10 * time.Millisecond
				g.log.Warn().
					Int("queue-size", len(g.incomingMessageQueue)).
					Stringer("backoff", waitTime).
					Msg("queue congested, retrying after backoff")
				time.Sleep(waitTime)
			}
		}
	}
}

func (g *Game) sendServerDifficulty(p *Player) {
	g.WritePacket(p, packet.ClientboundServerDifficulty{
		Difficulty:       packet.DifficultyNormal,
		DifficultyLocked: true,
	})
}

func (g *Game) sendJoinGameMessage(p *Player, dimensionCodec *nbt.Compound) {
	g.WritePacket(p, packet.ClientboundJoinGame{
		EntityID:         1,
		Hardcore:         false,
		Gamemode:         packet.GamemodeSurvival,
		PreviousGamemode: packet.GamemodeUnknown,
		WorldNames: []id.ID{
			id.ParseID("world"),
		},
		DimensionCodec: dimensionCodec,
		Dimension: dimensionCodec.
			Value["minecraft:dimension_type"].(*nbt.Compound).
			Value["value"].(*nbt.List).
			Value[0].(*nbt.Compound).
			Value["element"],
		WorldName:           id.ParseID("world"),
		HashedSeed:          g.world.WorldGenSettings.Seed,
		MaxPlayers:          100,
		ViewDistance:        5,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		Debug:               true,
		Flat:                false,
	})
}

func (g *Game) loadPlayerEntity(p *Player) error {
	data, err := g.world.LoadNBTPlayerdata(p.UUID)
	if err != nil {
		return fmt.Errorf("load nbt data: %w", err)
	}

	if err := entity.PlayerFromNBTIntoPlayer(data, p.Player); err != nil {
		return fmt.Errorf("decode nbt: %w", err)
	}

	return nil
}

func (g *Game) start(ctx context.Context) {
	g.prepare()

	// TODO: maybe more workers?
	go g.workIncomingMessageQueue(ctx)

	g.log.Debug().
		Stringer("tick", TickDuration).
		Msg("starting tick loop")
	ticker := time.NewTicker(TickDuration)
	lastTime := time.Now()
tickLoop:
	for {
		select {
		case <-ctx.Done():
			break tickLoop
		case t := <-ticker.C:
			sinceLast := time.Since(lastTime)
			if sinceLast > 2*TickDuration {
				g.log.Info().
					Stringer("sinceLast", sinceLast).
					Int("skipped", int(sinceLast/TickDuration)-1).
					Msg("can't keep up, skipping ticks")
			}
			lastTime = t

			g.tick()
			g.currentTick++
		}
	}

	g.log.Debug().
		Msg("stopped tick loop")
}

func (g *Game) workIncomingMessageQueue(ctx context.Context) {
workLoop:
	for {
		var msg incomingMessage
		select {
		case <-ctx.Done():
			break workLoop
		case msg = <-g.incomingMessageQueue:
		}

		g.log.Trace().
			IPAddr("source", msg.source.conn.IP()).
			Str("player", msg.source.name).
			Str("type", msg.pkg.Name()).
			Msg("processing message")

		g.processPacket(msg.source, msg.pkg)
	}

	g.log.Debug().
		Msg("stopped message worker")
}

func (g *Game) loadWorld() error {
	worldLoadStart := time.Now()
	loaded, err := LoadWorld(g.log, g.fs)
	if err != nil {
		return fmt.Errorf("load world: %w", err)
	}
	g.world = loaded
	g.log.Debug().
		Stringer("took", time.Since(worldLoadStart)).
		Msg("loaded world")
	return nil
}

func (g *Game) loadChunkAtCoord(coord voxel.V2) (*chunk.Chunk, error) {
	if loaded, ok := g.loadedChunks[coord]; ok {
		return loaded, nil
	}
	start := time.Now()
	regionCoord := voxel.V2{coord.X >> 5, coord.Z >> 5}
	region, err := g.world.loadRegion(regionCoord)
	if err != nil {
		return nil, err
	}
	chunk, err := region.loadChunk(coord)
	if err != nil {
		return nil, fmt.Errorf("load chunk: %w", err)
	}

	// remember that this chunk is currently loaded
	g.loadChunk(chunk)
	g.log.Debug().
		Stringer("took", time.Since(start)).
		Stringer("chunk", coord).
		Msg("load chunk")
	return chunk, nil
}

func (g *Game) loadChunk(ch *chunk.Chunk) {
	g.loadedChunks[ch.Coord] = ch
}

// generateAndLoadChunk will generate and load the chunk with the given coordinates.
// This will fail if the chunk already exists and is loaded, but not if the chunk
// already exists and is not loaded. If the chunk already exists (and is not loaded),
// The chunk will be overwritten.
func (g *Game) generateAndLoadChunk(coord voxel.V2) (*chunk.Chunk, error) {
	if _, ok := g.loadedChunks[coord]; ok {
		return nil, fmt.Errorf("chunk %s already exists and is loaded", coord)
	}
	chunk := g.world.generator.GenerateChunk(coord)
	g.loadChunk(chunk)
	return chunk, nil
}
