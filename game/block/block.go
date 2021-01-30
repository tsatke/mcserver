package block

import "github.com/tsatke/mcserver/game/id"

type Block struct {
	Name       id.ID
	Properties map[string]interface{}
}
