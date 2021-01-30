package entity

import (
	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/id"
)

type Data struct {
	ID                id.ID
	Pos               [3]float64
	Motion            [3]float64
	Rotation          [2]float32
	FallDistance      float32
	Fire              int16
	Air               int16
	OnGround          bool
	NoGravity         bool
	Invulnerable      bool
	PortalCooldown    int
	UUID              uuid.UUID
	CustomName        string
	CustomNameVisible bool
	Silent            bool
	Passengers        []interface{} // probably Mob or Player
	Glowing           bool
	Tags              []interface{} // to be done
}

func (d *Data) EntityID() id.ID {
	return d.ID
}

type CanBreed struct {
	InLove    int
	Age       int
	ForcedAge int
	LoveCause uuid.UUID
}

type CanBeAngry struct {
	AngerTime int
	AngryAt   uuid.UUID
}

type CanBeTamed struct {
	Owner   uuid.UUID
	Sitting bool
}
