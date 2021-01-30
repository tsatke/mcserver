package game

import (
	"sync"

	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/entity"
	"github.com/tsatke/mcserver/network"
)

type Player struct {
	sync.Mutex

	// tempUUID is not intended for use except for the exact point in time where a player connects.
	// This is used to pass the player UUID into the game. DON'T USE IT ANYWHERE!
	tempUUID uuid.UUID
	name     string
	conn     *network.Conn

	// client holds attributes regarding the player client, such as the brand, settings and others.
	client playerClient

	*entity.Player
}

type playerClient struct {
	// brand is the brand of the client that the player uses.
	// If this is empty, the client may not have sent a plugin message with channel minecraft:brand yet.
	// There is no guarantee that he will send it at any point. The Mojang minecraft client sends 'vanilla'
	// as its brand.
	brand    string
	settings playerClientSettings
}

type playerClientSettings struct {
	locale       string
	viewDistance int
}

func NewPlayer(uuid uuid.UUID, name string, conn *network.Conn) *Player {
	return &Player{
		tempUUID: uuid,
		name:     name,
		conn:     conn,
	}
}

func (p *Player) Disconnect() {
	_ = p.conn.Close()
}
