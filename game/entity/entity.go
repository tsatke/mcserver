package entity

import "github.com/tsatke/mcserver/game/id"

type Entity interface {
	EntityID() id.ID // correlates to the actual type of this entity
}
