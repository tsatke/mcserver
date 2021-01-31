package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundDeclareRecipes{}))
}

type Recipe struct {
	Type id.ID
	ID   id.ID
	// TODO: Data not supported yet
}

type ClientboundDeclareRecipes struct {
	Recipes []Recipe
}

func (ClientboundDeclareRecipes) ID() ID       { return IDClientboundDeclareRecipes }
func (ClientboundDeclareRecipes) Name() string { return "Declare Recipes" }

func (c ClientboundDeclareRecipes) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteVarInt("num recipes", len(c.Recipes))
	for i, recipe := range c.Recipes {
		enc.WriteID("recipe["+strconv.Itoa(i)+"] type", recipe.Type)
		enc.WriteID("recipe["+strconv.Itoa(i)+"] id", recipe.ID)
	}

	return
}
